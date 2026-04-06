package parsers

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/baldurstod/vdf"
	"github.com/rs/zerolog"
)

var lootListItemRegexp = regexp.MustCompile(`^\[(.+?)\](.+)$`)

// NOTE: UNFINISHED
func ParseClientLootLists(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.ClientLootList {
	logger := zerolog.Ctx(ctx)

	// Measure the time it takes to parse the loot lists, just
	// for performance monitoring and debugging purposes.
	start := time.Now()

	// We need to get all available "main" loot lists to then find those in the "client_loot_lists" section.
	// Instead of doing some funky string-checking, there is a direct connection between the
	// "revolving_loot_lists" and "client_loot_lists" sections in the items_game.txt file.
	// The "revolving_loot_lists" section contains the main loot lists, while the "client_loot_lists"
	// section contains the sub-lists that contain the actual items.
	revolving_loot_lists, err := ig.Get("revolving_loot_lists")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get revolving_loot_lists from items_game.txt")
		return nil
	}

	// At this point, we have all the main loot lists, and can continue to parse the "client_loot_lists" section.
	client_loot_lists, err := ig.Get("client_loot_lists")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get client_loot_lists from items_game.txt")
		return nil
	}

	revolving_loot_list_stringmap, _ := revolving_loot_lists.ToStringMap()
	if revolving_loot_list_stringmap == nil {
		logger.Warn().Msg("No revolving loot lists found in items_game.txt")
		return nil
	}

	buffer := make([]models.ClientLootList, 0)

	for series, list := range *revolving_loot_list_stringmap {
		if !IsValidLootListName(list) {
			// logger.Debug().Msgf("Skipping loot list '%s' for series '%s' as it is not a valid loot list name", list, series)
			continue
		}

		series_int, _ := strconv.Atoi(series)

		current := models.ClientLootList{
			Series:       series_int,
			LootListId:   list,
			SubLootLists: make([]models.ClientLootListSubList, 0),
		}

		all, err := client_loot_lists.Get(list)
		// logger.Debug().Msgf("Processing loot list '%s' for series '%d'", list, series_int)

		if err != nil {
			logger.Error().Err(err).Msgf("Failed to get client loot list '%s' from client_loot_lists", list)
			continue
		}

		for _, sub := range all.GetChilds() {
			current_sub_list := models.ClientLootListSubList{
				LootListName: sub.Key,
			}

			current_sub_list.Rarity = GetLootListRarity(sub.Key)
			current_sub_list.Items = GetLootListItems(client_loot_lists, sub.Key)

			current.SubLootLists = append(current.SubLootLists, current_sub_list)
		}

		buffer = append(buffer, current)
	}

	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' loot lists in %s", len(buffer), duration)

	return buffer
}

func IsValidLootListName(name string) bool {
	skip := []string{
		"_musickit",
		"_signature",
		"_signatures",
		"_xray_p250",
		"_dhw13_promo",
		"_promo_",
		"crate_ems14_",
		"storageunit1_",
		"crate_pins",
		"storageunit0_",
		"giftpkg_",
	}

	for _, s := range skip {
		if strings.Contains(name, s) {
			// Skip music kits and other irrelevant loot lists
			// log.Println("Skipping loot list:", name, "as it contains", s)
			return false
		}
	}
	return true
}

func GetLootListRarity(name string) string {
	rarity_endings := []string{
		"default",
		"common",
		"uncommon",
		"rare",
		"mythical",
		"legendary",
		"ancient",
		"immortal",
		"unusual",
	}

	for _, ending := range rarity_endings {
		if strings.HasSuffix(name, ending) {
			return ending
		}
	}

	// If no rarity ending is found, return "default"
	return "default"
}

func GetLootListItems(kv *vdf.KeyValue, loot_list string) []models.LootListItem {
	items := make([]models.LootListItem, 0)

	list, err := kv.Get(loot_list)
	if err != nil {
		return items
	}

	for _, v := range list.GetChilds() {
		res := lootListItemRegexp.FindStringSubmatch(v.Key)

		if len(res) < 3 {
			continue
		}

		name := res[1]
		class := res[2]

		items = append(items, models.LootListItem{
			Name:  name,
			Class: class,
		})
	}

	return items
}
