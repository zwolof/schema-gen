package parsers

import (
	"context"

	"go-csitems-parser/models"

	"github.com/rs/zerolog"
)

// GetSkinWeaponRarityMap parses the "client_loot_lists" section of items_game and returns a
// map of "[paintkit]weapon" → rarity string (e.g. "mythical", "rare", "uncommon").
//
// It does this by examining every sub-list key whose name ends in a known rarity suffix
// (e.g. "crate_community_11_mythical" → "mythical") and collecting all items directly
// nested under that key.  Keys that resolve to "default" (no recognisable suffix) are
// skipped because they are aggregator lists that simply reference other sub-lists.
func GetSkinWeaponRarityMap(ctx context.Context, ig *models.ItemsGame) map[string]string {
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
			// Aggregator list or unrecognised suffix — skip.
			continue
		}

		for _, item := range subList.GetChilds() {
			// Items are formatted as "[paintkit]weapon_class"
			matches := lootListItemRegexp.FindStringSubmatch(item.Key)
			if len(matches) < 3 {
				// Not a direct item entry (e.g. it's a nested list reference).
				continue
			}
			// Reconstruct the canonical key exactly as it appears in the file.
			key := "[" + matches[1] + "]" + matches[2]
			// First writer wins — if the same skin appears in multiple lists, keep
			// the rarity from the first (usually lower-rarity) occurrence.
			if _, exists := result[key]; !exists {
				result[key] = rarity
			}
		}
	}

	logger.Info().Msgf("Built skin-weapon rarity map with %d entries", len(result))
	return result
}
