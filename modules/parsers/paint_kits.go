package parsers

import (
	"context"
	"strconv"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

func ParsePaintKits(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.PaintKit {
	logger := zerolog.Ctx(ctx)

	start := time.Now()

	paint_kits, err := ig.Get("paint_kits")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get paint_kits from items_game.txt")
		return nil
	}

	raritymap := GetPaintKitRarityStringMap(ctx, ig)

	if raritymap == nil {
		logger.Error().Msg("Failed to get paint_kits_rarity from items_game.txt")
		return nil
	}

	var items []models.PaintKit
	for _, r := range paint_kits.GetChilds() {
		name, _ := r.GetString("name")

		// skip if name equal to "workshop_default"
		if name == "workshop_default" {
			continue
		}

		definition_index, _ := strconv.Atoi(r.Key)
		wear_remap_min, _ := r.GetFloat32("wear_remap_min")
		wear_remap_max, _ := r.GetFloat32("wear_remap_max")
		description_tag, _ := r.GetString("description_tag")
		description_string, _ := r.GetString("description_string")

		float_ranges := models.PaintKitWearRange{
			Min: wear_remap_min,
			Max: wear_remap_max,
		}

		if float_ranges.Max == 0.0 {
			float_ranges.Max = 1.0
		}

		translated_description, _ := t.GetValueByKey(description_string)

		current := models.PaintKit{
			DefinitionIndex: definition_index,
			Name:            name,
			Wear:            float_ranges,
			Description:     translated_description,
			MarketHashName:  modules.GenerateMarketHashName(t, description_tag, &name, "paint_kit"),
		}

		val, exists := raritymap[current.Name]
		if !exists {
			logger.Warn().Msgf("No rarity found for paint kit '%s' (definition index: %d)", current.Name, current.DefinitionIndex)
		}

		current.Rarity = val

		items = append(items, current)
	}

	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' paintkits in %s", len(items), duration)

	return items
}

func GetPaintKitRarityStringMap(ctx context.Context, ig *models.ItemsGame) map[string]string {
	paint_kits_rarity, err := ig.Get("paint_kits_rarity")
	logger := zerolog.Ctx(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get paint_kits_rarity from items_game.txt")
		return nil
	}

	// Create a map to hold the rarity strings
	rmap, err := paint_kits_rarity.ToStringMap()

	if err != nil {
		logger.Error().Err(err).Msg("Failed to convert paint_kits_rarity to string map")
		return nil
	}

	return *rmap
}
