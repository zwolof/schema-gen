package meta

import (
	"context"
	"regexp"

	"go-csitems-parser/internal/models"
	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/baldurstod/vdf"
	"github.com/rs/zerolog"
)

var itemSetKeyRegexp = regexp.MustCompile(`^\[(.+?)\](.+)$`)

// ItemSets extracts item_sets and enriches each with HasCrate/HasSouvenir
// derived from the pre-computed Inputs.WeaponCases / SouvenirPackages.
// Publishes into Inputs.ItemSets for Tier-2 parsers.
type ItemSets struct{ base.Parser }

func NewItemSets() *ItemSets { return &ItemSets{Parser: base.Internal("item_sets")} }

func (is *ItemSets) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	raw, err := in.IG.Get("item_sets")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get item_sets from items_game.txt")
		return nil, nil
	}

	var out []models.ItemSet
	defer is.LogCount(ctx, "item sets", func() int { return len(out) })()

	for _, s := range raw.GetChilds() {
		name, _ := s.GetString("name")

		current := models.ItemSet{
			Key:  s.Key,
			Name: name,
			Type: models.ItemSetTypePaintKits,
		}

		itemsetItems, _ := s.Get("items")
		items := itemSetPaintKits(itemsetItems)

		if len(items) == 0 {
			agents := itemSetAgents(itemsetItems)
			if len(agents) > 0 {
				current.Agents = agents
				current.Type = models.ItemSetTypeAgents
			} else {
				continue
			}
		} else {
			current.Items = items
		}

		for _, wpncase := range in.WeaponCases {
			if wpncase.ItemSetId == nil || *wpncase.ItemSetId != current.Key {
				continue
			}
			current.HasCrate = true
			break
		}

		for _, sv := range in.SouvenirPackages {
			if sv.ItemSetId == nil || *sv.ItemSetId != current.Key {
				continue
			}
			current.HasSouvenir = true
			break
		}

		out = append(out, current)
	}

	// Synthesize a pseudo item-set for skins removed from circulation.
	// Valve stores these in client_loot_lists["removed_items"] rather than in
	// item_sets, so they are invisible to the normal item-set pass above.
	// Example: the M4A4 | Howl (cu_m4a1_howling / weapon_m4a1).
	if removedSet := parseRemovedItemsSet(ctx, in.IG); removedSet != nil {
		out = append(out, *removedSet)
		logger.Info().Msgf("Synthesized removed_items item set with %d skin(s)", len(removedSet.Items))
	}

	return out, nil
}

func (is *ItemSets) Commit(in *pipeline.Inputs, result any) {
	if r, ok := result.([]models.ItemSet); ok {
		in.ItemSets = r
	}
}

// parseRemovedItemsSet reads client_loot_lists["removed_items"] and builds a
// synthetic ItemSet so that removed skins (e.g. M4A4 | Howl) flow through the
// normal Tier-2 weapon-skin builder unchanged.
func parseRemovedItemsSet(ctx context.Context, ig *models.ItemsGame) *models.ItemSet {
	logger := zerolog.Ctx(ctx)

	clientLootLists, err := ig.Get("client_loot_lists")
	if err != nil {
		logger.Warn().Err(err).Msg("parseRemovedItemsSet: client_loot_lists not found")
		return nil
	}

	var removedKV *vdf.KeyValue
	for _, sub := range clientLootLists.GetChilds() {
		if sub.Key == "removed_items" {
			removedKV = sub
			break
		}
	}
	if removedKV == nil {
		return nil
	}

	// Reuse itemSetPaintKits — it iterates children and matches [paintkit]weaponclass.
	// Non-matching keys (e.g. "public_list_contents") are silently skipped.
	items := itemSetPaintKits(removedKV)
	if len(items) == 0 {
		return nil
	}

	return &models.ItemSet{
		Key:         "removed_items",
		Name:        "Removed Items",
		Type:        models.ItemSetTypePaintKits,
		Items:       items,
		HasCrate:    false,
		HasSouvenir: false,
	}
}

func itemSetAgents(kv *vdf.KeyValue) []string {
	agents := make([]string, 0)
	for _, item := range kv.GetChilds() {
		agents = append(agents, item.Key)
	}
	return agents
}

func itemSetPaintKits(kv *vdf.KeyValue) []models.ItemSetItem {
	skins := make([]models.ItemSetItem, 0)
	for _, skin := range kv.GetChilds() {
		res := itemSetKeyRegexp.FindStringSubmatch(skin.Key)
		if len(res) < 3 {
			continue
		}
		skins = append(skins, models.ItemSetItem{
			PaintKitName: res[1],
			WeaponClass:  res[2],
		})
	}
	return skins
}
