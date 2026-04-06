package parsers

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

var valid_capsule_prefabs = []string{
	"crate_sticker_pack_",
	"crate_signature_pack_",
}

func ParseStickerCapsules(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.StickerCapsule {
	logger := zerolog.Ctx(ctx)

	start := time.Now()

	items, err := ig.Get("items")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get collectibles from items_game.txt")
		return nil
	}

	var sticker_capsules []models.StickerCapsule

	for _, item := range items.GetChilds() {
		definition_index, _ := strconv.Atoi(item.Key)
		item_name, _ := item.GetString("item_name")
		name, _ := item.GetString("name")
		image_inventory, _ := item.GetString("image_inventory")
		item_description, _ := item.GetString("item_description")

		if !IsValidStickerCapsule(name) {
			continue
		}

		item_set := modules.GetSupplyCrateSeries(item, ig)

		if item_set == nil {
			continue
		}

		var current = models.StickerCapsule{
			DefinitionIndex: definition_index,
			Name:            name,
			ImageInventory:  image_inventory,
			ItemDescription: item_description,
			ItemSetId:       item_set,
			MarketHashName:  modules.GenerateMarketHashName(t, item_name, nil, "sticker_capsule"),
		}

		sticker_capsules = append(sticker_capsules, current)
	}

	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' sticker capsules in %s", len(sticker_capsules), duration)

	return sticker_capsules
}

func IsValidStickerCapsule(prefab string) bool {
	for _, valid_prefab := range valid_capsule_prefabs {
		if strings.Contains(prefab, valid_prefab) {
			return true
		}
	}
	return false
}
