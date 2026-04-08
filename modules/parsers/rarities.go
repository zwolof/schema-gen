package parsers

import (
	"context"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

func ParseRarities(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) map[string]models.Rarity {
	logger := zerolog.Ctx(ctx)

	start := time.Now()

	rarities, err := ig.Get("rarities")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get rarities from items_game.txt")
		return nil
	}

	colors, _ := ig.Get("colors")
	color_map := make(map[string]models.GenericColor)

	for _, clr := range colors.GetChilds() {
		color_name, _ := clr.GetString("color_name")
		hex_color, _ := clr.GetString("hex_color")

		color_map[clr.Key] = models.GenericColor{
			Key:       clr.Key,
			ColorName: color_name,
			HexColor:  hex_color,
		}
	}

	items := make(map[string]models.Rarity)
	for _, r := range rarities.GetChilds() {
		loc_key, _ := r.GetString("loc_key")
		loc_key_weapon, _ := r.GetString("loc_key_weapon")
		loc_key_character, _ := r.GetString("loc_key_character")

		if loc_key == "" || loc_key_weapon == "" || loc_key_character == "" {
			logger.Warn().Msgf("Rarity '%s' is missing one of the localization keys, skipping", r.Key)
			continue
		}

		translated_rarity, _ := t.GetValueByKey(loc_key)
		translated_weapon, _ := t.GetValueByKey(loc_key_weapon)
		translated_character, _ := t.GetValueByKey(loc_key_character)

		current := models.Rarity{
			LocRarity:    translated_rarity,
			LocWeapon:    translated_weapon,
			LocCharacter: translated_character,
		}

		// Get color Data
		color_str, _ := r.GetString("color")

		// loop through the color map to find the matching color
		if color_str != "" {
			if colorData, ok := color_map[color_str]; ok {

				current.Hex = colorData.HexColor
				// current.ColorName = colorData.ColorName
			}
		}

		items[r.Key] = current
	}

	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' rarities in %s", len(items), duration)

	return items
}
