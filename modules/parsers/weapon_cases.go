package parsers

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ParseWeaponCases(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.WeaponCase {
	logger := zerolog.Ctx(ctx)

	start := time.Now()
	// logger.Info().Msg("Parsing weapon cases...")

	items, err := ig.Get("items")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get collectibles from items_game.txt")
		return nil
	}

	var weapon_cases []models.WeaponCase

	// Iterate through all items in the "items" section
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

		// Get child key called "attributes"
		associated_items, _ := item.Get("associated_items")

		case_key_def_idx := -1 // Default to -1 if not found
		if associated_items != nil {
			value := associated_items.GetChilds()[0].Key

			if value != "" {
				case_key_def_idx, _ = strconv.Atoi(value)
			}
		}

		// If case_key_def_idx is still -1, we cannot find the key for this case
		case_key := GetWeaponCaseKeyByDefIndex(ig, case_key_def_idx)
		item_set := modules.GetContainerItemSet(item, t, "ItemSet")

		// Create the weapon case model
		var current = models.WeaponCase{
			DefinitionIndex: definition_index,
			ImageInventory:  image_inventory,
			Key:            case_key,
			ItemSetId:      item_set,
			MarketHashName: modules.GenerateMarketHashName(t, item_name, nil, "weapon_case"),
			FirstSaleDate:  first_sale_date,
			Description:    item_description,
		}

		weapon_cases = append(weapon_cases, current)
	}

	// Save music kits to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' weapon cases in %s", len(weapon_cases), duration)

	return weapon_cases
}

func GetWeaponCaseKeyByDefIndex(ig *models.ItemsGame, definitionIndex int) *models.WeaponCaseKey {
	items, err := ig.Get("items")

	if err != nil {
		log.Error().Err(err).Msg("gg Failed to get items from items_game.txt")
		return nil
	}

	var current models.WeaponCaseKey
	for _, item := range items.GetChilds() {
		def_idx, _ := strconv.Atoi(item.Key)
		if def_idx != definitionIndex {
			continue
		}

		prefab, _ := item.GetString("prefab")

		if !strings.Contains(prefab, "weapon_case_key") {
			continue
		}

		name, _ := item.GetString("name")
		image_inventory, _ := item.GetString("image_inventory")

		// item_name, _ := item.GetString("item_name")
		// item_description, _ := item.GetString("item_description")
		// first_sale_date, _ := item.GetString("first_sale_date")
		current = models.WeaponCaseKey{
			DefinitionIndex: def_idx,
			Name:            name,
			ImageInventory:  image_inventory,
			// Prefab:          prefab,
			// ItemName:        item_name,
			// ItemDescription: item_description,
			// FirstSaleDate:   first_sale_date,
		}

		break // We found the item, no need to continue
	}

	// if current.Prefab == "" {
	// 	log.Error().Msgf("No weapon case key found for definition index %d", definitionIndex)
	// 	return nil
	// }

	return &current
}
