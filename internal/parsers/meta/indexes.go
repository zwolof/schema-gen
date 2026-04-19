package meta

import (
	"context"
	"strings"

	"go-csitems-parser/internal/models"

	"github.com/rs/zerolog"
)

// SkinWeaponRarityMap parses client_loot_lists and returns a map of
// "[paintkit]weaponclass" → rarity string. First writer wins.
func SkinWeaponRarityMap(ctx context.Context, ig *models.ItemsGame) map[string]string {
	logger := zerolog.Ctx(ctx)

	clientLootLists, err := ig.Get("client_loot_lists")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get client_loot_lists from items_game.txt")
		return nil
	}

	result := make(map[string]string)
	for _, subList := range clientLootLists.GetChilds() {
		rarity := lootListRarity(subList.Key)
		if rarity == "default" {
			continue
		}
		for _, item := range subList.GetChilds() {
			matches := lootListItemRegexp.FindStringSubmatch(item.Key)
			if len(matches) < 3 {
				continue
			}
			key := "[" + matches[1] + "]" + matches[2]
			if _, exists := result[key]; !exists {
				result[key] = rarity
			}
		}
	}

	logger.Info().Msgf("Built skin-weapon rarity map with %d entries", len(result))
	return result
}

// StickerItemSetMap parses client_loot_lists and returns a map of
// sticker-kit name → item_set_id (sub-list key with rarity suffix stripped).
func StickerItemSetMap(ctx context.Context, ig *models.ItemsGame) map[string]string {
	logger := zerolog.Ctx(ctx)

	clientLootLists, err := ig.Get("client_loot_lists")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get client_loot_lists from items_game.txt")
		return nil
	}

	result := make(map[string]string)
	for _, subList := range clientLootLists.GetChilds() {
		rarity := lootListRarity(subList.Key)
		if rarity == "default" {
			continue
		}
		itemSetID := strings.TrimSuffix(subList.Key, "_"+rarity)

		for _, item := range subList.GetChilds() {
			matches := lootListItemRegexp.FindStringSubmatch(item.Key)
			if len(matches) < 3 || matches[2] != "sticker" {
				continue
			}
			kit := matches[1]
			if _, exists := result[kit]; !exists {
				result[kit] = itemSetID
			}
		}
	}

	logger.Info().Msgf("Built sticker item-set map with %d entries", len(result))
	return result
}
