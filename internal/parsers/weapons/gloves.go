package weapons

import (
	"context"
	"strconv"

	"go-csitems-parser/internal/marketname"
	"go-csitems-parser/internal/models"
	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/rs/zerolog"
)

// Gloves extracts hands_paintable prefab items and publishes them as
// Inputs.Gloves for the Tier-2 glove skin builder.
type Gloves struct{ base.Parser }

func NewGloves() *Gloves { return &Gloves{Parser: base.Internal("gloves")} }

func (g *Gloves) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	items, err := in.IG.Get("items")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get items")
		return nil, nil
	}

	var out []models.BaseWeapon
	defer g.LogCount(ctx, "gloves", func() int { return len(out) })()

	for _, w := range items.GetChilds() {
		prefab, _ := w.GetString("prefab")
		if prefab != prefabHandsPaintable {
			continue
		}

		definition_index, _ := strconv.Atoi(w.Key)
		classname, _ := w.GetString("name")
		item_name, _ := w.GetString("item_name")

		out = append(out, models.BaseWeapon{
			DefinitionIndex: definition_index,
			ClassName:       classname,
			Name:            marketname.GenerateMarketHashName(in.T, item_name, nil, "glove"),
		})
	}

	return out, nil
}

func (g *Gloves) Commit(in *pipeline.Inputs, result any) {
	if r, ok := result.([]models.BaseWeapon); ok {
		in.Gloves = r
	}
}
