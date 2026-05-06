package meta

import (
	"context"

	"go-csitems-parser/internal/models"
	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/rs/zerolog"
)

// Rarities extracts the items_game.txt rarities section, merged with
// translation and colour lookups.
type Rarities struct{ base.Parser }

func NewRarities() *Rarities { return &Rarities{Parser: base.New("rarities")} }

func (r *Rarities) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	rarities, err := in.IG.Get("rarities")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get rarities from items_game.txt")
		return nil, nil
	}

	colors, _ := in.IG.Get("colors")
	colorMap := make(map[string]models.GenericColor)
	if colors != nil {
		for _, clr := range colors.GetChilds() {
			colorName, _ := clr.GetString("color_name")
			hexColor, _ := clr.GetString("hex_color")
			colorMap[clr.Key] = models.GenericColor{
				Key:       clr.Key,
				ColorName: colorName,
				HexColor:  hexColor,
			}
		}
	}

	out := make(map[string]models.Rarity)
	defer r.LogCount(ctx, "rarities", func() int { return len(out) })()

	for _, entry := range rarities.GetChilds() {
		locKey, _ := entry.GetString("loc_key")
		locKeyWeapon, _ := entry.GetString("loc_key_weapon")
		locKeyCharacter, _ := entry.GetString("loc_key_character")

		if locKey == "" || locKeyWeapon == "" || locKeyCharacter == "" {
			logger.Warn().Msgf("Rarity '%s' is missing one of the localization keys, skipping", entry.Key)
			continue
		}

		translatedRarity, _ := in.T.GetValueByKey(locKey)
		translatedWeapon, _ := in.T.GetValueByKey(locKeyWeapon)
		translatedChar, _ := in.T.GetValueByKey(locKeyCharacter)

		current := models.Rarity{
			LocRarity:    translatedRarity,
			LocWeapon:    translatedWeapon,
			LocCharacter: translatedChar,
		}

		colorStr, _ := entry.GetString("color")
		if colorStr != "" {
			if cd, ok := colorMap[colorStr]; ok {
				current.Hex = cd.HexColor
			}
		}

		out[entry.Key] = current
	}

	return out, nil
}

// RarityKeys parses the rarities block from items_game.txt and returns the set
// of valid rarity token strings (e.g. "common", "uncommon", "rare" …).
// Used to drive loot-list rarity detection without hard-coding rarity names.
func RarityKeys(ctx context.Context, ig *models.ItemsGame) map[string]bool {
	rarities, err := ig.Get("rarities")
	if err != nil {
		return nil
	}
	keys := make(map[string]bool, len(rarities.GetChilds()))
	for _, entry := range rarities.GetChilds() {
		keys[entry.Key] = true
	}
	return keys
}
