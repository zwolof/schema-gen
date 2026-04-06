package modules

import (
	"fmt"
	"go-csitems-parser/models"
)

type WeaponToPaintKitMap struct {
	DefinitionIndex int    `json:"definition_index"`
	Name            string `json:"name"`
	StickerCount    int    `json:"sticker_count"`
	Type            string `json:"type"`
}

type WeaponSkinMap struct {
	BaseItem models.BaseWeapon `json:"base_item"`
	Paints   []models.PaintKit `json:"paints"`
}

type SchemaWeaponSkinMap struct {
	Name          string                                 `json:"name"`
	StickerAmount int                                    `json:"sticker_count"`
	Type          string                                 `json:"type"`
	Paints        map[int]models.SchemaWeaponPaintKitMap `json:"paints"`
}

type CollectibleMap struct {
	MarketHashName string `json:"market_hash_name"`
	Rarity         string `json:"rarity"`
	Image          string `json:"image"`
}

func GetWeaponPaintKits(
	weapons *[]models.BaseWeapon,
	paint_kits *[]models.PaintKit,
	item_sets *[]models.ItemSet,
) map[int]SchemaWeaponSkinMap {
	weapon_skin_map := make(map[int]SchemaWeaponSkinMap, 0)

	for _, weapon := range *weapons {
		// Create a new glove skin map entry
		current := SchemaWeaponSkinMap{
			Name:          weapon.Name,
			StickerAmount: weapon.NumStickers,
			Type:          "weapon",
			Paints:        make(map[int]models.SchemaWeaponPaintKitMap),
		}

		item_set_paint_kits := GetItemSetPaintKitsForWeapon(item_sets, weapon.ClassName)
		for _, paint_kit_name := range item_set_paint_kits {
			paint_kit := GetPaintKitByName(paint_kits, paint_kit_name)

			if paint_kit == nil {
				continue
			}

			data := GetPaintKitWeaponCombinationData(item_sets, weapon.ClassName, paint_kit.Name)

			if data != nil {
				paint_kit.StatTrak = data.CanBeStatTrak
				paint_kit.Souvenir = data.CanBeSouvenir
				paint_kit.ItemSetId = data.ItemSetId
			}

			current.Paints[paint_kit.DefinitionIndex] = models.SchemaWeaponPaintKitMap{
				DefinitionIndex: paint_kit.DefinitionIndex,
				Float:           paint_kit.Wear,
				Rarity:          paint_kit.Rarity,
				Image:           fmt.Sprintf("%s_%s", weapon.ClassName, paint_kit.Name),
				Name:            paint_kit.MarketHashName,
				ItemSetId:       paint_kit.ItemSetId,
				Souvenir:        paint_kit.Souvenir,
				StatTrak:        paint_kit.StatTrak,
			}
		}

		weapon_skin_map[weapon.DefinitionIndex] = current
	}

	return weapon_skin_map
}

func GetItemSetPaintKitsForWeapon(
	item_sets *[]models.ItemSet,
	weapon_name string,
) []string {
	paint_kits := make([]string, 0)

	for _, item_set := range *item_sets {
		// if item_set.Type != models.ItemSetTypePaintKits {
		// 	continue
		// }

		for _, item := range item_set.Items {
			if item.WeaponClass == weapon_name {
				paint_kits = append(paint_kits, item.PaintKitName)
			}
		}
	}

	return paint_kits
}

func GetKnifePaintKits(
	knives *[]models.BaseWeapon,
	paint_kits *[]models.PaintKit,
	knife_map map[string][]string,
) map[int]SchemaWeaponSkinMap {
	weapon_skin_map := make(map[int]SchemaWeaponSkinMap, 0)

	for _, knife := range *knives {
		// Create a new glove skin map entry
		current := SchemaWeaponSkinMap{
			Name:          knife.Name,
			StickerAmount: 0, // Knives don't have stickers
			Type:          "knife",
			Paints:        make(map[int]models.SchemaWeaponPaintKitMap),
		}

		knife_map_value, ok := knife_map[knife.ClassName]
		if !ok {
			continue
		}

		for _, pk_name := range knife_map_value {
			for _, paint_kit := range *paint_kits {
				if paint_kit.Name != pk_name {
					continue
				}

				// Add the paint kit to the current glove skin map
				current.Paints[paint_kit.DefinitionIndex] = models.SchemaWeaponPaintKitMap{
					DefinitionIndex: paint_kit.DefinitionIndex,
					Float:           paint_kit.Wear,
					Rarity:          paint_kit.Rarity,
					Image:           fmt.Sprintf("%s_%s", knife.ClassName, paint_kit.Name),
					Name:            paint_kit.MarketHashName,
					ItemSetId:       paint_kit.ItemSetId,
					Souvenir:        false, // Knives can NOT be Souvenir
					StatTrak:        true,  // Knives can always be StatTrak
				}
			}
		}

		weapon_skin_map[knife.DefinitionIndex] = current
	}

	return weapon_skin_map
}

type GloveSkinMap struct {
	BaseItem  models.BaseWeapon `json:"base_item"`
	PaintKits []models.PaintKit `json:"paint_kits"`
}

func GetGlovePaintKits(gloves *[]models.BaseWeapon, paint_kits *[]models.PaintKit, glove_map map[string][]string) map[int]SchemaWeaponSkinMap {
	weapon_skin_map := make(map[int]SchemaWeaponSkinMap, 0)

	for _, glove := range *gloves {
		// Create a new glove skin map entry
		current := SchemaWeaponSkinMap{
			Name:          glove.Name,
			StickerAmount: 0,
			Type:          "glove",
			Paints:        make(map[int]models.SchemaWeaponPaintKitMap),
		}

		glove_map_value, ok := glove_map[glove.ClassName]
		if !ok {
			continue
		}

		for _, pk_name := range glove_map_value {
			for _, paint_kit := range *paint_kits {
				if paint_kit.Name != pk_name {
					continue
				}

				// Add the paint kit to the current glove skin map
				current.Paints[paint_kit.DefinitionIndex] = models.SchemaWeaponPaintKitMap{
					DefinitionIndex: paint_kit.DefinitionIndex,
					Float:           paint_kit.Wear,
					Rarity:          paint_kit.Rarity,
					Image:           fmt.Sprintf("%s_%s", glove.ClassName, paint_kit.Name),
					Name:            paint_kit.MarketHashName,
					ItemSetId:       paint_kit.ItemSetId,
					Souvenir:        false, // Gloves can NOT be Souvenir
					StatTrak:        false, // Gloves can NOT be StatTrak
				}
			}
		}

		weapon_skin_map[glove.DefinitionIndex] = current
	}

	return weapon_skin_map
}

func GetPaintKitByName(paint_kits *[]models.PaintKit, name string) *models.PaintKit {
	for _, paint_kit := range *paint_kits {
		if paint_kit.Name == name {
			return &paint_kit
		}
	}
	return nil
}

func GetWeaponByClass(weapons *[]models.BaseWeapon, weapon_class string) *models.BaseWeapon {
	for _, wpn := range *weapons {
		if wpn.ClassName == weapon_class {
			return &wpn
		}
	}
	return nil
}

func GetPaintKitWeaponCombinationData(item_sets *[]models.ItemSet, cn string, pk string) *models.PaintKitWeaponCombinationData {
	for _, item_set := range *item_sets {
		for _, item := range item_set.Items {
			if item.WeaponClass == cn && item.PaintKitName == pk {
				return &models.PaintKitWeaponCombinationData{
					ItemSetId:     item_set.Key,
					CanBeStatTrak: item_set.HasCrate,
					CanBeSouvenir: item_set.HasSouvenir,
				}
			}
		}
	}
	return nil
}
