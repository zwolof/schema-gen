// Package containers holds parsers for item containers: weapon cases,
// souvenir packages, sticker capsules, and misc self-opening capsules
// (patches, graffiti, etc.).
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

// WeaponCases extracts the weapon_case prefab items from items_game.txt.
type WeaponCases struct{ base.Parser }

func NewWeaponCases() *WeaponCases { return &WeaponCases{Parser: base.New("weapon_cases")} }

func (w *WeaponCases) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	items, err := in.IG.Get("items")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get items from items_game.txt")
		return nil, nil
	}

	var out []models.WeaponCase
	defer w.LogCount(ctx, "weapon cases", func() int { return len(out) })()

	for _, item := range items.GetChilds() {
		prefab, _ := item.GetString("prefab")
		if prefab != "weapon_case" {
			continue
		}

		definition_index, _ := strconv.Atoi(item.Key)
		item_name, _ := item.GetString("item_name")
		image_inventory, _ := item.GetString("image_inventory")
		item_description, _ := item.GetString("item_description")
		first_sale_date, _ := item.GetString("first_sale_date")

		associated_items, _ := item.Get("associated_items")
		case_key_def_idx := -1
		if associated_items != nil {
			v := associated_items.GetChilds()[0].Key
			if v != "" {
				case_key_def_idx, _ = strconv.Atoi(v)
			}
		}

		out = append(out, models.WeaponCase{
			DefinitionIndex: definition_index,
			ImageInventory:  image_inventory,
			Key:             lookupCaseKey(in.IG, case_key_def_idx),
			ItemSetId:       itemsgame.GetContainerItemSet(item, "ItemSet"),
			MarketHashName:  marketname.GenerateMarketHashName(in.T, item_name, nil, "weapon_case"),
			FirstSaleDate:   first_sale_date,
			Description:     item_description,
		})
	}

	return out, nil
}

func (w *WeaponCases) Commit(in *pipeline.Inputs, result any) {
	if r, ok := result.([]models.WeaponCase); ok {
		in.WeaponCases = r
	}
}

// lookupCaseKey finds the weapon_case_key item matching definitionIndex.
func lookupCaseKey(ig *models.ItemsGame, definitionIndex int) *models.WeaponCaseKey {
	items, err := ig.Get("items")
	if err != nil {
		return nil
	}

	for _, item := range items.GetChilds() {
		defIdx, _ := strconv.Atoi(item.Key)
		if defIdx != definitionIndex {
			continue
		}
		prefab, _ := item.GetString("prefab")
		if !strings.Contains(prefab, "weapon_case_key") {
			continue
		}
		name, _ := item.GetString("name")
		imageInventory, _ := item.GetString("image_inventory")
		return &models.WeaponCaseKey{
			DefinitionIndex: defIdx,
			Name:            name,
			ImageInventory:  imageInventory,
		}
	}
	return &models.WeaponCaseKey{}
}
