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

// Knives extracts melee_unusual prefab items and publishes them as
// Inputs.Knives for the Tier-2 knife skin builder.
type Knives struct{ base.Parser }

func NewKnives() *Knives { return &Knives{Parser: base.Internal("knives")} }

func (k *Knives) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	items, err := in.IG.Get("items")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get items")
		return nil, nil
	}

	var out []models.BaseWeapon
	defer k.LogCount(ctx, "knives", func() int { return len(out) })()

	for _, w := range items.GetChilds() {
		prefab, _ := w.GetString("prefab")
		if prefab != prefabMeleeUnusual {
			continue
		}

		definition_index, _ := strconv.Atoi(w.Key)
		item_name, _ := w.GetString("item_name")
		classname, _ := w.GetString("name")
		image_inventory, _ := w.GetString("image_inventory")

		out = append(out, models.BaseWeapon{
			DefinitionIndex: definition_index,
			ClassName:       classname,
			Name:            marketname.GenerateMarketHashName(in.T, item_name, nil, "knife"),
			ImageInventory:  "econ/default_generated/" + image_inventory + "_light",
		})
	}

	return out, nil
}

func (k *Knives) Commit(in *pipeline.Inputs, result any) {
	if r, ok := result.([]models.BaseWeapon); ok {
		in.Knives = r
	}
}
