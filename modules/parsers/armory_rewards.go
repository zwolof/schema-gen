package parsers

import (
	"context"
	"strings"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

func ParseArmoryRewards(
	ctx context.Context, 
	ig *models.ItemsGame, 
	item_sets *[]models.ItemSet,
	t *modules.Translator,
) []models.OperationPointsRedeemableItem {
	logger := zerolog.Ctx(ctx)

	start := time.Now()

	seasonal_operations, err := ig.Get("seasonaloperations")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get seasonal operations from items_game.txt")
		return nil
	}

	var armory_items []models.OperationPointsRedeemableItem
	var current_item models.OperationPointsRedeemableItem

	for _, item := range seasonal_operations.GetChilds() {
		redeemable_goods, _ := item.GetString("redeemable_goods")

		if redeemable_goods != "xpshop" {
			continue
		}
		
		goods := item.GetChilds()

		if len(goods) == 0 {
			continue
		}

		for _, g := range goods {
			if g.Key != "operational_point_redeemable" {
				logger.Warn().Msgf("Unexpected redeemable good type '%s' for item '%s'", g.Key, item.Key)
				continue
			}

			item_name, _ := g.GetString("item_name")
			ui_image_thumbnail, _ := g.GetString("ui_image_thumbnail")

			item_set_id := strings.Replace(item_name, "lootlist:", "", 1)
			callout, _ := g.GetString("callout")
			points, _ := g.GetInt("points")
			ui_order, _ := g.GetInt("ui_order")

			translated_name, err := t.GetValueByKey(callout)
			if err != nil {
				logger.Error().Err(err).Msgf("Failed to translate callout '%s'", callout)
				continue
			}

			current_item = models.OperationPointsRedeemableItem{
				Points:         points,
				Name:            translated_name,
				ItemSetId:      item_set_id,
				UIImageThumbnail:  ui_image_thumbnail,
				UIOrder:        ui_order,
			}

			armory_items = append(armory_items, current_item)
		}
	}

	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' armory items in %s", len(armory_items), duration)

	return armory_items
}
