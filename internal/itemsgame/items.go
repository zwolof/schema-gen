// Package items loads and queries Valve's items_game.txt tree. Load parses
// the raw VDF file, merging duplicate root sections into a single KeyValue
// tree wrapped as *models.ItemsGame. The remaining helpers are pure readers
// over that tree — they memoise lookups that parsers hit many times
// (pro_players, revolving_loot_lists) via sync.Map.
package itemsgame

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"

	"go-csitems-parser/internal/models"

	"github.com/baldurstod/vdf"
)

// Load reads items_game.txt at path, merging root-level duplicate sections.
func Load(path string) *models.ItemsGame {
	fileData, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	fileData = bytes.Trim(fileData, "\xef\xbb\xbf")

	parser := vdf.VDF{}
	parsed := parser.Parse(fileData)

	kv, _ := parsed.Get("items_game")
	if kv == nil {
		panic("items_game.txt does not contain 'items_game' section")
	}

	mergeKeysAtRootLevel(kv)

	return &models.ItemsGame{KeyValue: kv}
}

// LoadKnifeSkinsMap reads the knife/glove paint-kit overrides from
// knife_skins.json at path.
func LoadKnifeSkinsMap(path string) map[string][]string {
	fileData, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Error reading file %s: %v", path, err))
	}
	if len(fileData) == 0 {
		panic(fmt.Sprintf("File %s is empty", path))
	}

	result := make(map[string][]string)
	json.Unmarshal(fileData, &result)
	return result
}

// mergeKeysAtRootLevel flattens duplicate top-level sections (items_game can
// contain multiple "prefabs" blocks etc.) into a single merged tree.
func mergeKeysAtRootLevel(root *vdf.KeyValue) {
	if root == nil {
		panic("root KeyValue is nil")
	}
	if len(root.GetChilds()) == 0 {
		panic("root KeyValue has no child keys")
	}

	cached := make(map[string][]*vdf.KeyValue)
	for _, item := range root.GetChilds() {
		section := item.Key
		cached[section] = append(cached[section], item.GetChilds()...)
	}

	newRoot := &vdf.KeyValue{
		Key:   "items_game",
		Value: make([]*vdf.KeyValue, 0),
	}
	for section, items := range cached {
		if len(items) == 0 {
			continue
		}
		newRoot.Value = append(newRoot.Value.([]*vdf.KeyValue), &vdf.KeyValue{
			Key:   section,
			Value: items,
		})
	}
	root.Value = newRoot.Value
}

// --- VDF helpers that read from an items_game tree ---

// GetTournamentEventId reads attributes.tournament event id.value from an item.
func GetTournamentEventId(item *vdf.KeyValue) (int, error) {
	attributes, err := item.Get("attributes")
	if err != nil {
		return -1, err
	}
	tournament, err := attributes.Get("tournament event id")
	if err != nil {
		return -1, err
	}
	return tournament.GetInt("value")
}

// GetContainerItemSet reads tags.<key>.tag_value from an item. key defaults
// to "ItemSet" when empty.
func GetContainerItemSet(item *vdf.KeyValue, key string) *string {
	tags, err := item.Get("tags")
	if err != nil {
		return nil
	}

	containerKey := "ItemSet"
	if key != "" {
		containerKey = key
	}

	itemSet, err := tags.Get(containerKey)
	if err != nil {
		return nil
	}

	tag, _ := itemSet.GetString("tag_value")
	return &tag
}

// GetSupplyCrateSeries reads "set supply crate series" from an item and
// resolves it against the cached revolving_loot_lists index.
func GetSupplyCrateSeries(item *vdf.KeyValue, ig *models.ItemsGame) *string {
	attributes, err := item.Get("attributes")
	if err != nil {
		return nil
	}
	set, err := attributes.Get("set supply crate series")
	if err != nil {
		return nil
	}
	seriesID, err := set.GetString("value")
	if err != nil {
		return nil
	}

	if value, ok := revolvingLootListMap(ig)[seriesID]; ok {
		return &value
	}
	return nil
}

// GetPlayerByAccountId looks up a pro player by Steam account id. Thread-safe.
func GetPlayerByAccountId(ig *models.ItemsGame, accountID int) *models.TournamentData {
	return proPlayerIndex(ig)[accountID]
}

// --- memoised indexes ---

// proPlayerIndexCache: *models.ItemsGame → map[int]*models.TournamentData.
var proPlayerIndexCache sync.Map

func proPlayerIndex(ig *models.ItemsGame) map[int]*models.TournamentData {
	if v, ok := proPlayerIndexCache.Load(ig); ok {
		return v.(map[int]*models.TournamentData)
	}
	out := make(map[int]*models.TournamentData)
	if players, err := ig.Get("pro_players"); err == nil && players != nil {
		for _, p := range players.GetChilds() {
			aid, _ := strconv.Atoi(p.Key)
			name, _ := p.GetString("name")
			out[aid] = &models.TournamentData{Id: aid, Name: name}
		}
	}
	actual, _ := proPlayerIndexCache.LoadOrStore(ig, out)
	return actual.(map[int]*models.TournamentData)
}

// revolvingLootListCache: *models.ItemsGame → map[string]string (series_id → value).
var revolvingLootListCache sync.Map

func revolvingLootListMap(ig *models.ItemsGame) map[string]string {
	if v, ok := revolvingLootListCache.Load(ig); ok {
		return v.(map[string]string)
	}
	out := make(map[string]string)
	if revolving, err := ig.Get("revolving_loot_lists"); err == nil && revolving != nil {
		for _, list := range revolving.GetChilds() {
			val, _ := list.ToString()
			out[list.Key] = val
		}
	}
	actual, _ := revolvingLootListCache.LoadOrStore(ig, out)
	return actual.(map[string]string)
}
