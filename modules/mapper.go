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
	Image 	   string                                 `json:"image"`
	StickerAmount int                                    `json:"sticker_count"`
	Type          string                                 `json:"type"`
	Paints        map[int]models.SchemaWeaponPaintKitMap `json:"paints"`
}

type CollectibleMap struct {
	MarketHashName string `json:"market_hash_name"`
	Rarity         string `json:"rarity"`
	Image          string `json:"image"`
}

// buildPaintKitIndex builds a name → *PaintKit lookup map to avoid O(n) scans.
func buildPaintKitIndex(paint_kits *[]models.PaintKit) map[string]*models.PaintKit {
	idx := make(map[string]*models.PaintKit, len(*paint_kits))
	for i := range *paint_kits {
		pk := &(*paint_kits)[i]
		idx[pk.Name] = pk
	}
	return idx
}

// itemSetIndex groups item-set items and metadata by weapon class for fast lookup.
type itemSetIndex struct {
	// paintKits maps weapon_class → []paint_kit_name
	paintKits map[string][]string
	// combo maps "weapon_class|paint_kit_name" → combination data
	combo map[string]*models.PaintKitWeaponCombinationData
}

func buildItemSetIndex(item_sets *[]models.ItemSet) *itemSetIndex {
	idx := &itemSetIndex{
		paintKits: make(map[string][]string),
		combo:     make(map[string]*models.PaintKitWeaponCombinationData),
	}
	for i := range *item_sets {
		is := &(*item_sets)[i]
		data := &models.PaintKitWeaponCombinationData{
			ItemSetId:     is.Key,
			CanBeStatTrak: is.HasCrate,
			CanBeSouvenir: is.HasSouvenir,
		}
		for _, item := range is.Items {
			idx.paintKits[item.WeaponClass] = append(idx.paintKits[item.WeaponClass], item.PaintKitName)
			key := item.WeaponClass + "|" + item.PaintKitName
			idx.combo[key] = data
		}
	}
	return idx
}

func GetWeaponPaintKits(
	weapons *[]models.BaseWeapon,
	paint_kits *[]models.PaintKit,
	item_sets *[]models.ItemSet,
) map[int]SchemaWeaponSkinMap {
	pkIdx := buildPaintKitIndex(paint_kits)
	isIdx := buildItemSetIndex(item_sets)

	weapon_skin_map := make(map[int]SchemaWeaponSkinMap, len(*weapons))

	for _, weapon := range *weapons {
		current := SchemaWeaponSkinMap{
			Name:          weapon.Name,
			StickerAmount: weapon.NumStickers,
			Image: fmt.Sprintf("econ/default_generated/%s_light", weapon.ClassName),
			Type:          "weapon",
			Paints:        make(map[int]models.SchemaWeaponPaintKitMap),
		}

		for _, paint_kit_name := range isIdx.paintKits[weapon.ClassName] {
			pk, ok := pkIdx[paint_kit_name]
			if !ok {
				continue
			}

			if data := isIdx.combo[weapon.ClassName+"|"+pk.Name]; data != nil {
				pk.StatTrak = data.CanBeStatTrak
				pk.Souvenir = data.CanBeSouvenir
				pk.ItemSetId = data.ItemSetId
			}

			current.Paints[pk.DefinitionIndex] = models.SchemaWeaponPaintKitMap{
				DefinitionIndex: pk.DefinitionIndex,
				Float:           pk.Wear,
				Rarity:          pk.Rarity,
				Image:           fmt.Sprintf("econ/default_generated/%s_%s_light", weapon.ClassName, pk.Name),
				Name:            pk.MarketHashName,
				ItemSetId:       pk.ItemSetId,
				Souvenir:        pk.Souvenir,
				StatTrak:        pk.StatTrak,
			}
		}

		weapon_skin_map[weapon.DefinitionIndex] = current
	}

	return weapon_skin_map
}

func GetKnifePaintKits(
	knives *[]models.BaseWeapon,
	paint_kits *[]models.PaintKit,
	knife_map map[string][]string,
) map[int]SchemaWeaponSkinMap {
	pkIdx := buildPaintKitIndex(paint_kits)
	weapon_skin_map := make(map[int]SchemaWeaponSkinMap, len(*knives))

	for _, knife := range *knives {
		current := SchemaWeaponSkinMap{
			Name:          knife.Name,
			StickerAmount: 0,
			Type:          "knife",
			Paints:        make(map[int]models.SchemaWeaponPaintKitMap),
		}

		for _, pk_name := range knife_map[knife.ClassName] {
			pk, ok := pkIdx[pk_name]
			if !ok {
				continue
			}
			current.Paints[pk.DefinitionIndex] = models.SchemaWeaponPaintKitMap{
				DefinitionIndex: pk.DefinitionIndex,
				Float:           pk.Wear,
				Rarity:          pk.Rarity,
				Image:           fmt.Sprintf("%s_%s", knife.ClassName, pk.Name),
				Name:            pk.MarketHashName,
				ItemSetId:       pk.ItemSetId,
				Souvenir:        false,
				StatTrak:        true,
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

func GetGlovePaintKits(
	gloves *[]models.BaseWeapon,
	paint_kits *[]models.PaintKit,
	glove_map map[string][]string,
) map[int]SchemaWeaponSkinMap {
	pkIdx := buildPaintKitIndex(paint_kits)
	weapon_skin_map := make(map[int]SchemaWeaponSkinMap, len(*gloves))

	for _, glove := range *gloves {
		current := SchemaWeaponSkinMap{
			Name:          glove.Name,
			StickerAmount: 0,
			Type:          "glove",
			Paints:        make(map[int]models.SchemaWeaponPaintKitMap),
		}

		for _, pk_name := range glove_map[glove.ClassName] {
			pk, ok := pkIdx[pk_name]
			if !ok {
				continue
			}
			current.Paints[pk.DefinitionIndex] = models.SchemaWeaponPaintKitMap{
				DefinitionIndex: pk.DefinitionIndex,
				Float:           pk.Wear,
				Rarity:          pk.Rarity,
				Image:           fmt.Sprintf("econ/default_generated/%s_%s_light", glove.ClassName, pk.Name),
				Name:            pk.MarketHashName,
				ItemSetId:       pk.ItemSetId,
				Souvenir:        false,
				StatTrak:        false,
			}
		}

		weapon_skin_map[glove.DefinitionIndex] = current
	}

	return weapon_skin_map
}

func GetPaintKitByName(paint_kits *[]models.PaintKit, name string) *models.PaintKit {
	for i := range *paint_kits {
		if (*paint_kits)[i].Name == name {
			return &(*paint_kits)[i]
		}
	}
	return nil
}

func GetWeaponByClass(weapons *[]models.BaseWeapon, weapon_class string) *models.BaseWeapon {
	for i := range *weapons {
		if (*weapons)[i].ClassName == weapon_class {
			return &(*weapons)[i]
		}
	}
	return nil
}

func GetItemSetPaintKitsForWeapon(item_sets *[]models.ItemSet, weapon_name string) []string {
	paint_kits := make([]string, 0)
	for _, item_set := range *item_sets {
		for _, item := range item_set.Items {
			if item.WeaponClass == weapon_name {
				paint_kits = append(paint_kits, item.PaintKitName)
			}
		}
	}
	return paint_kits
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

