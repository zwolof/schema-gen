package meta

import (
	"context"
	"strconv"

	"go-csitems-parser/internal/marketname"
	"go-csitems-parser/internal/models"
	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/rs/zerolog"
)

// PaintKits extracts paint_kits, annotating each with its rarity from the
// paint_kits_rarity table. Publishes into Inputs.PaintKits for Tier-2 skin
// builders.
type PaintKits struct{ base.Parser }

func NewPaintKits() *PaintKits { return &PaintKits{Parser: base.New("paint_kits")} }

func (p *PaintKits) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	paintKits, err := in.IG.Get("paint_kits")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get paint_kits from items_game.txt")
		return nil, nil
	}

	rarityMap := paintKitRarityMap(ctx, in.IG)
	if rarityMap == nil {
		logger.Error().Msg("Failed to get paint_kits_rarity from items_game.txt")
		return nil, nil
	}

	var out []models.PaintKit
	defer p.LogCount(ctx, "paintkits", func() int { return len(out) })()

	for _, r := range paintKits.GetChilds() {
		name, _ := r.GetString("name")
		if name == "workshop_default" {
			continue
		}

		definitionIndex, _ := strconv.Atoi(r.Key)
		wearMin, _ := r.GetFloat32("wear_remap_min")
		wearMax, _ := r.GetFloat32("wear_remap_max")
		descriptionTag, _ := r.GetString("description_tag")
		descriptionString, _ := r.GetString("description_string")

		// Fallback: some kits ship with max=0, which makes the JSON wear range
		// useless. Treat that as "full 0-1 range".
		if wearMax == 0.0 {
			wearMax = 1.0
		}

		translatedDescription, _ := in.T.GetValueByKey(descriptionString)

		current := models.PaintKit{
			DefinitionIndex: definitionIndex,
			Name:            name,
			Description:     translatedDescription,
			Wear:            models.PaintKitWearRange{Min: wearMin, Max: wearMax},
			MarketHashName:  marketname.GenerateMarketHashName(in.T, descriptionTag, &name, "paint_kit"),
		}

		rarity, exists := rarityMap[current.Name]
		if !exists {
			logger.Warn().Msgf("No rarity found for paint kit '%s' (definition index: %d)", current.Name, current.DefinitionIndex)
		}
		current.Rarity = rarity

		out = append(out, current)
	}

	return out, nil
}

func (p *PaintKits) Commit(in *pipeline.Inputs, result any) {
	if r, ok := result.([]models.PaintKit); ok {
		in.PaintKits = r
	}
}

func paintKitRarityMap(ctx context.Context, ig *models.ItemsGame) map[string]string {
	logger := zerolog.Ctx(ctx)

	rarities, err := ig.Get("paint_kits_rarity")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get paint_kits_rarity from items_game.txt")
		return nil
	}

	m, err := rarities.ToStringMap()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to convert paint_kits_rarity to string map")
		return nil
	}
	return *m
}
