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

		item_name, _ := w.GetString("item_name")
		item_description, _ := w.GetString("item_description")
		image_inventory, _ := w.GetString("image_inventory")
		max_num_stickers, _ := game_info.GetInt("max_num_stickers")

		translated_name, err := t.GetValueByKey(item_name)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to translate item name for weapon %s", item_name)
			translated_name = item_name // Fallback to original if translation fails
		}

		translated_description, err := t.GetValueByKey(item_description)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to translate item description for weapon %s", item_description)
			translated_description = item_description // Fallback to original if translation fails
		}

		current := models.BaseWeapon{
			DefinitionIndex: def_idx,
			Name:            translated_name,
			Description:     translated_description,
			ClassName:       item_class,
			ImageInventory:  image_inventory,
			NumStickers:     max_num_stickers,
		}

		weapons = append(weapons, current)
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
