package parsers

import (
	"context"
	"strings"

	"go-csitems-parser/models"

	"github.com/rs/zerolog"
)

// GetStickerItemSetMap parses "client_loot_lists" and returns a map of
// sticker kit name → item_set_id (the list key prefix without the rarity suffix).
//
// For example, a sub-list key "crate_sticker_pack_kat2014_01_rare" contains
// items like "[kat2014_3dmax]sticker", so the function maps:
//
//	"kat2014_3dmax" → "crate_sticker_pack_kat2014_01"
func GetStickerItemSetMap(ctx context.Context, ig *models.ItemsGame) map[string]string {
	logger := zerolog.Ctx(ctx)

	client_loot_lists, err := ig.Get("client_loot_lists")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get client_loot_lists from items_game.txt")
		return nil
	}

	result := make(map[string]string)

	for _, subList := range client_loot_lists.GetChilds() {
		rarity := GetLootListRarity(subList.Key)
		if rarity == "default" {
			continue
		}

		// Derive the item_set_id by stripping the rarity suffix (and the underscore before it).
		itemSetId := strings.TrimSuffix(subList.Key, "_"+rarity)

		for _, item := range subList.GetChilds() {
			// Sticker items are formatted as "[kit_name]sticker"
			matches := lootListItemRegexp.FindStringSubmatch(item.Key)
			if len(matches) < 3 || matches[2] != "sticker" {
				continue
			}
			kitName := matches[1]
			if _, exists := result[kitName]; !exists {
				result[kitName] = itemSetId
			}
		}
	}

	logger.Info().Msgf("Built sticker item-set map with %d entries", len(result))
	return result
}
