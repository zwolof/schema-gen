package parsers

import (
	"context"
	"strconv"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

func ParseKeychains(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.Keychain {
	logger := zerolog.Ctx(ctx)

	start := time.Now()

	keychain_definitions, err := ig.Get("keychain_definitions")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get keychain_definitions, is items_game.txt valid?")
		return nil
	}

	var agents []models.Keychain
	for _, mk := range keychain_definitions.GetChilds() {
		definition_index, _ := strconv.Atoi(mk.Key)
		name, _ := mk.GetString("name")

		if name == "kc_aus2025" {
			continue // Skip the AUS 2025 keychain, it's not a valid keychain
		}

		loc_name, _ := mk.GetString("loc_name")
		image_inventory, _ := mk.GetString("image_inventory")
		item_rarity, _ := mk.GetString("item_rarity")

		current := models.Keychain{
			DefinitionIndex: definition_index,
			MarketHashName:  modules.GenerateMarketHashName(t, loc_name, nil, "keychain"),
			Rarity:          item_rarity,
			ImageInventory:  image_inventory,
		}

		agents = append(agents, current)
	}

	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' keychains in %s", len(agents), duration)

	return agents
}
