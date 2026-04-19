package meta

import (
	"context"
	"strings"

	"go-csitems-parser/internal/models"
	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/rs/zerolog"
)

// ArmoryRewards extracts the Operation Armory redeemable-goods list from
// seasonaloperations.
type ArmoryRewards struct{ base.Parser }

func NewArmoryRewards() *ArmoryRewards {
	return &ArmoryRewards{Parser: base.New("armory_rewards")}
}

func (a *ArmoryRewards) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	ops, err := in.IG.Get("seasonaloperations")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get seasonal operations from items_game.txt")
		return nil, nil
	}

	var out []models.OperationPointsRedeemableItem
	defer a.LogCount(ctx, "armory items", func() int { return len(out) })()

	for _, item := range ops.GetChilds() {
		redeemableGoods, _ := item.GetString("redeemable_goods")
		if redeemableGoods != "xpshop" {
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

			itemName, _ := g.GetString("item_name")
			uiImageThumbnail, _ := g.GetString("ui_image_thumbnail")
			itemSetID := strings.Replace(itemName, "lootlist:", "", 1)
			callout, _ := g.GetString("callout")
			points, _ := g.GetInt("points")
			uiOrder, _ := g.GetInt("ui_order")

			translatedName, err := in.T.GetValueByKey(callout)
			if err != nil {
				logger.Error().Err(err).Msgf("Failed to translate callout '%s'", callout)
				continue
			}

			out = append(out, models.OperationPointsRedeemableItem{
				Points:           points,
				Name:             translatedName,
				ItemSetId:        itemSetID,
				UIImageThumbnail: uiImageThumbnail,
				UIOrder:          uiOrder,
			})
		}
	}

	return out, nil
}
