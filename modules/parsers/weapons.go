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

var invalid_weapon_prefabs = []string{
	"grenade",
	"equipment",
	"weapon_fire_grenade_prefab",
	"weapon_hegrenade_prefab",
}

var special_cases = []string{
	"weapon_knife_t",
	"weapon_knife",
}

func ParseWeapons(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.BaseWeapon {
	logger := zerolog.Ctx(ctx)

	start := time.Now()

	prefabs, err := ig.Get("prefabs")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get prefabs")
		return nil
	}

	game_info, err := ig.Get("game_info")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get game_info")
		return nil
	}

	items, err := ig.Get("items")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get items")
		return nil
	}

	max_num_stickers, _ := game_info.GetInt("max_num_stickers")

	buildWeapon := func(className string, kv interface {
		GetString(string) (string, error)
	}, defIdx int) models.BaseWeapon {
		item_name, _ := kv.GetString("item_name")
		item_description, _ := kv.GetString("item_description")
		image_inventory, _ := kv.GetString("image_inventory")

		translated_name, err := t.GetValueByKey(item_name)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to translate item name for weapon %s", item_name)
			translated_name = item_name
		}
		translated_description, err := t.GetValueByKey(item_description)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to translate item description for weapon %s", item_description)
			translated_description = item_description
		}

		return models.BaseWeapon{
			DefinitionIndex: defIdx,
			Name:            translated_name,
			Description:     translated_description,
			ClassName:       className,
			ImageInventory:  image_inventory,
			NumStickers:     max_num_stickers,
		}
	}

	var weapons []models.BaseWeapon

	for _, w := range prefabs.GetChilds() {
		if !strings.HasPrefix(w.Key, "weapon_") || !strings.HasSuffix(w.Key, "_prefab") {
			continue
		}

		item_class := strings.TrimSuffix(w.Key, "_prefab")
		def_idx := GetBaseWeaponDefinitionIndex(item_class, ig)

		if def_idx == -1 {
			logger.Error().Msgf("Failed to get definition index for weapon class '%s'", item_class)
			continue
		}
		_, err := w.Get("paint_data")
		if err != nil && item_class != "weapon_taser" {
			continue
		}

		weapons = append(weapons, buildWeapon(item_class, w, def_idx))
	}

	// Valve did not give these items a dedicated _prefab entry — they live directly
	// in the "items" section and are found by their "name" key.
	specialSet := make(map[string]struct{}, len(special_cases))
	for _, sc := range special_cases {
		specialSet[sc] = struct{}{}
	}

	for _, item := range items.GetChilds() {
		name, _ := item.GetString("name")
		if _, ok := specialSet[name]; !ok {
			continue
		}
		definition_index, _ := strconv.Atoi(item.Key)
		weapons = append(weapons, buildWeapon(name, item, definition_index))
	}

	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' weapons in %s", len(weapons), duration)

	return weapons
}

func GetBaseWeaponDefinitionIndex(class string, ig *models.ItemsGame) int {
	items, err := ig.Get("items")
	if err != nil {
		return -1
	}

	for _, w := range items.GetChilds() {
		name, _ := w.GetString("name")

		if name == class {
			definition_index, _ := strconv.Atoi(w.Key)
			return definition_index
		}
	}

	return -1 // Class not found
}
