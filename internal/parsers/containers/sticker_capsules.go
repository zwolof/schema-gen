package containers

import (
	"context"
	"strconv"
	"strings"

	"go-csitems-parser/internal/itemsgame"
	"go-csitems-parser/internal/marketname"
	"go-csitems-parser/internal/models"
	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/rs/zerolog"
)

var stickerCapsulePrefabs = []string{
	"crate_sticker_pack_",
	"crate_signature_pack_",
}

// StickerCapsules extracts sticker-pack containers from items_game.txt.
type StickerCapsules struct{ base.Parser }

func NewStickerCapsules() *StickerCapsules {
	return &StickerCapsules{Parser: base.New("sticker_capsules")}
}

func (s *StickerCapsules) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	items, err := in.IG.Get("items")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get items from items_game.txt")
		return nil, nil
	}

	var out []models.StickerCapsule
	defer s.LogCount(ctx, "sticker capsules", func() int { return len(out) })()

	for _, item := range items.GetChilds() {
		name, _ := item.GetString("name")
		if !isStickerCapsule(name) {
			continue
		}

		itemSet := itemsgame.GetSupplyCrateSeries(item, in.IG)
		if itemSet == nil {
			continue
		}

		definition_index, _ := strconv.Atoi(item.Key)
		item_name, _ := item.GetString("item_name")
		image_inventory, _ := item.GetString("image_inventory")
		item_description, _ := item.GetString("item_description")

		out = append(out, models.StickerCapsule{
			DefinitionIndex: definition_index,
			Name:            name,
			ImageInventory:  image_inventory,
			ItemDescription: item_description,
			ItemSetId:       itemSet,
			MarketHashName:  marketname.GenerateMarketHashName(in.T, item_name, nil, "sticker_capsule"),
		})
	}

	return out, nil
}

func isStickerCapsule(name string) bool {
	for _, p := range stickerCapsulePrefabs {
		if strings.Contains(name, p) {
			return true
		}
	}
	return false
}
