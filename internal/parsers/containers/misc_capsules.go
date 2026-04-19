package containers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-csitems-parser/internal/itemsgame"
	"go-csitems-parser/internal/marketname"
	"go-csitems-parser/internal/models"
	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/rs/zerolog"
)

var miscCapsulePrefabs = []string{
	"patch_capsule",
	"graffiti_box",
	"weapon_case_selfopening_collection",
}

var prefabTagType = map[string]string{
	"patch_capsule":                      "PatchCapsule",
	"stockh2021_patch_capsule_prefab":    "PatchCapsule",
	"graffiti_box":                       "SprayCapsule",
	"weapon_case_selfopening_collection": "ItemSet",
	"weapon_case_base":                   "ItemSet",
}

// MiscCapsules extracts patch/graffiti/self-opening containers that don't fit
// elsewhere (registered as "misc_capsules" in the pipeline).
type MiscCapsules struct{ base.Parser }

func NewMiscCapsules() *MiscCapsules {
	return &MiscCapsules{Parser: base.New("misc_capsules")}
}

func (m *MiscCapsules) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	items, err := in.IG.Get("items")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get items from items_game.txt")
		return nil, nil
	}

	var out []models.StickerCapsule
	defer m.LogCount(ctx, "self-opening capsules", func() int { return len(out) })()

	for _, item := range items.GetChilds() {
		prefab, _ := item.GetString("prefab")
		name, _ := item.GetString("name")
		if !isMiscSelfOpeningCapsule(prefab, name) {
			continue
		}

		itemName, _ := item.GetString("item_name")
		tagType := prefabTagType[prefab]
		if tagType == "" {
			logger.Warn().Msgf("Unknown tag type for item: %s", itemName)
			continue
		}

		itemSet := itemsgame.GetContainerItemSet(item, tagType)
		if itemSet == nil {
			itemSet = itemsgame.GetSupplyCrateSeries(item, in.IG)
			if itemSet == nil {
				fmt.Println("Item set is nil again, skipping item:", itemName)
				continue
			}
		}

		definition_index, _ := strconv.Atoi(item.Key)
		image_inventory, _ := item.GetString("image_inventory")
		item_description, _ := item.GetString("item_description")

		out = append(out, models.StickerCapsule{
			DefinitionIndex: definition_index,
			Name:            name,
			ImageInventory:  image_inventory,
			ItemDescription: item_description,
			ItemSetId:       itemSet,
			MarketHashName:  marketname.GenerateMarketHashName(in.T, itemName, nil, "self_opening_capsule"),
		})
	}

	return out, nil
}

func isMiscSelfOpeningCapsule(prefab, name string) bool {
	if strings.HasPrefix(name, "crate_xray_") || strings.HasPrefix(name, "crate_musickit_") {
		return true
	}
	for _, p := range miscCapsulePrefabs {
		if strings.Contains(prefab, p) {
			return true
		}
	}
	return false
}
