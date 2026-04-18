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

	var keychains []models.Keychain
	for _, mk := range keychain_definitions.GetChilds() {
		definition_index, _ := strconv.Atoi(mk.Key)
		name, _ := mk.GetString("name")

		loc_name, _ := mk.GetString("loc_name")
		image_inventory, _ := mk.GetString("image_inventory")
		item_rarity, _ := mk.GetString("item_rarity")
		is_commodity, _ := mk.GetBool("is_commodity")
		pedestal_display_model, _ := mk.GetString("pedestal_display_model")

		has_tags, tags_err := mk.Get("tags")

		current := models.Keychain{
			DefinitionIndex: definition_index,
			Name:            name,
			MarketHashName:  modules.GenerateMarketHashName(t, loc_name, nil, "keychain"),
			Rarity:          item_rarity,
			ImageInventory:  image_inventory,
			IsCommodity:     is_commodity,
			Model:           pedestal_display_model,
		}

		if tags_err == nil && has_tags != nil {
			current.IsSpecialCharm = true
		}

		keychains = append(keychains, current)
	}

	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' keychains in %s", len(keychains), duration)

	return keychains
}
