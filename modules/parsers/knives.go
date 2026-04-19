package parsers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

func ParseKnives(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.BaseWeapon {
	logger := zerolog.Ctx(ctx)

	start := time.Now()
	// logger.Info().Msg("Parsing music kits...")

	items, err := ig.Get("items")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get items")
		return nil
	}

	var knives []models.BaseWeapon
	for _, w := range items.GetChilds() {
		prefab, _ := w.GetString("prefab")

		if prefab != "melee_unusual" {
			// Skip non-knife items
			continue
		}

		definition_index, _ := strconv.Atoi(w.Key)
		item_name, _ := w.GetString("item_name")
		classname, _ := w.GetString("name")
		image_inventory, _ := w.GetString("image_inventory")

		current := models.BaseWeapon{
			DefinitionIndex: definition_index,
			ClassName:       classname,
			Name:            modules.GenerateMarketHashName(t, item_name, nil, "knife"),
			ImageInventory:  fmt.Sprintf("econ/default_generated/%s_light", image_inventory),
		}

		knives = append(knives, current)
	}

	// Save knives to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' knives in %s", len(knives), duration.String())

	return knives
}
