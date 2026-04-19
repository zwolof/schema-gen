// Package skinmap builds per-weapon / per-knife / per-glove skin maps for
// the exported schema. Builder is the interface the pipeline depends on;
// PaintKitBuilder is the default implementation.
package skinmap

import (
	"fmt"

	"go-csitems-parser/internal/models"
)

// Builder produces the schema-shaped skin maps. It also exposes
// EnrichPaintKits as the pre-Tier-2 mutation step that populates
// Souvenir/StatTrak/ItemSetId on each paint kit.
type Builder interface {
	Weapons(weapons []models.BaseWeapon, paintKits []models.PaintKit, itemSets []models.ItemSet, rarityMap map[string]string) map[int]SchemaWeaponSkinMap
	Knives(knives []models.BaseWeapon, paintKits []models.PaintKit, knifeSkinMap map[string][]string, rarityMap map[string]string) map[int]SchemaWeaponSkinMap
	Gloves(gloves []models.BaseWeapon, paintKits []models.PaintKit, gloveSkinMap map[string][]string, rarityMap map[string]string) map[int]SchemaWeaponSkinMap
	EnrichPaintKits(paintKits []models.PaintKit, itemSets []models.ItemSet, weapons []models.BaseWeapon)
}

// SchemaWeaponSkinMap is the per-base-item shape emitted by the Builder.
type SchemaWeaponSkinMap struct {
	Name          string                                 `json:"name"`
	Image         string                                 `json:"image"`
	StickerAmount int                                    `json:"sticker_count"`
	Type          string                                 `json:"type"`
	Paints        map[int]models.SchemaWeaponPaintKitMap `json:"paints"`
}

// PaintKitBuilder is the default file-driven Builder.
type PaintKitBuilder struct{}

// Default is the package-level Builder singleton used by the pipeline.
var Default Builder = PaintKitBuilder{}

// ---- Builder interface methods ----

func (PaintKitBuilder) Weapons(weapons []models.BaseWeapon, paintKits []models.PaintKit, itemSets []models.ItemSet, rarityMap map[string]string) map[int]SchemaWeaponSkinMap {
	pkIdx := buildPaintKitIndex(paintKits)
	isIdx := buildItemSetIndex(itemSets)

	out := make(map[int]SchemaWeaponSkinMap, len(weapons))
	for _, weapon := range weapons {
		current := SchemaWeaponSkinMap{
			Name:          weapon.Name,
			StickerAmount: weapon.NumStickers,
			Image:         fmt.Sprintf("econ/default_generated/%s_light", weapon.ClassName),
			Type:          "weapon",
			Paints:        make(map[int]models.SchemaWeaponPaintKitMap),
		}

		for _, pkName := range isIdx.paintKits[weapon.ClassName] {
			pk, ok := pkIdx[pkName]
			if !ok {
				continue
			}

			// Use per-(weapon, paint_kit) combo rather than pk fields — pk may
			// carry a different weapon's last-write combo when a paint kit
			// appears in multiple weapon item sets.
			var statTrak, souvenir bool
			var itemSetID string
			if data := isIdx.combo[weapon.ClassName+"|"+pk.Name]; data != nil {
				statTrak = data.CanBeStatTrak
				souvenir = data.CanBeSouvenir
				itemSetID = data.ItemSetId
			}

			rarity := pk.Rarity
			if r, ok := rarityMap["["+pk.Name+"]"+weapon.ClassName]; ok {
				rarity = r
			}

			current.Paints[pk.DefinitionIndex] = models.SchemaWeaponPaintKitMap{
				DefinitionIndex: pk.DefinitionIndex,
				Float:           pk.Wear,
				Rarity:          rarity,
				Image:           fmt.Sprintf("econ/default_generated/%s_%s_light", weapon.ClassName, pk.Name),
				Name:            pk.MarketHashName,
				ItemSetId:       itemSetID,
				Souvenir:        souvenir,
				StatTrak:        statTrak,
			}
		}
		out[weapon.DefinitionIndex] = current
	}
	return out
}

func (PaintKitBuilder) Knives(knives []models.BaseWeapon, paintKits []models.PaintKit, knifeSkinMap map[string][]string, _ map[string]string) map[int]SchemaWeaponSkinMap {
	pkIdx := buildPaintKitIndex(paintKits)
	out := make(map[int]SchemaWeaponSkinMap, len(knives))

	for _, knife := range knives {
		current := SchemaWeaponSkinMap{
			Name:   knife.Name,
			Type:   "knife",
			Paints: make(map[int]models.SchemaWeaponPaintKitMap),
		}

		for _, pkName := range knifeSkinMap[knife.ClassName] {
			pk, ok := pkIdx[pkName]
			if !ok {
				continue
			}

			image := fmt.Sprintf("econ/default_generated/%s_%s_light", knife.ClassName, pk.Name)
			if pk.Name == "default" {
				image = fmt.Sprintf("econ/weapons/base_weapons/%s", knife.ClassName)
			}

			current.Paints[pk.DefinitionIndex] = models.SchemaWeaponPaintKitMap{
				DefinitionIndex: pk.DefinitionIndex,
				Float:           pk.Wear,
				Rarity:          "ancient",
				Image:           image,
				Name:            pk.MarketHashName,
				// Legacy: knives ran before the enrichment pass, so ItemSetId
				// was always "" — preserved for downstream byte-exactness.
				ItemSetId: "",
				Souvenir:  false,
				StatTrak:  true,
			}
		}

		out[knife.DefinitionIndex] = current
	}

	return out
}

func (PaintKitBuilder) Gloves(gloves []models.BaseWeapon, paintKits []models.PaintKit, gloveSkinMap map[string][]string, _ map[string]string) map[int]SchemaWeaponSkinMap {
	pkIdx := buildPaintKitIndex(paintKits)
	out := make(map[int]SchemaWeaponSkinMap, len(gloves))

	for _, glove := range gloves {
		current := SchemaWeaponSkinMap{
			Name:   glove.Name,
			Type:   "glove",
			Paints: make(map[int]models.SchemaWeaponPaintKitMap),
		}

		for _, pkName := range gloveSkinMap[glove.ClassName] {
			pk, ok := pkIdx[pkName]
			if !ok {
				continue
			}

			current.Paints[pk.DefinitionIndex] = models.SchemaWeaponPaintKitMap{
				DefinitionIndex: pk.DefinitionIndex,
				Float:           pk.Wear,
				Rarity:          "ancient",
				Image:           fmt.Sprintf("econ/default_generated/%s_%s_light", glove.ClassName, pk.Name),
				Name:            pk.MarketHashName,
				// Gloves never carry an item-set id — see Knives above.
				ItemSetId: "",
				Souvenir:  false,
				StatTrak:  false,
			}
		}

		out[glove.DefinitionIndex] = current
	}

	return out
}

// EnrichPaintKits mutates paintKits in place so Tier-2 reads can proceed
// concurrently. Scoped to weapon classes only, matching the legacy behaviour
// where GetWeaponPaintKits was the sole site that mutated paint kits.
func (PaintKitBuilder) EnrichPaintKits(paintKits []models.PaintKit, itemSets []models.ItemSet, weapons []models.BaseWeapon) {
	pkIdx := buildPaintKitIndex(paintKits)
	isIdx := buildItemSetIndex(itemSets)

	for _, w := range weapons {
		for _, pkName := range isIdx.paintKits[w.ClassName] {
			pk, ok := pkIdx[pkName]
			if !ok {
				continue
			}

			if data := isIdx.combo[w.ClassName+"|"+pkName]; data != nil {
				pk.StatTrak = data.CanBeStatTrak
				pk.Souvenir = data.CanBeSouvenir
				pk.ItemSetId = data.ItemSetId
			}
		}
	}
}

// --- private indexes ---

// buildPaintKitIndex maps paint-kit name → *PaintKit (pointer into the
// caller's slice) so the builders can O(1) their joins. The pointer is used
// so EnrichPaintKits can mutate through it.
func buildPaintKitIndex(paintKits []models.PaintKit) map[string]*models.PaintKit {
	idx := make(map[string]*models.PaintKit, len(paintKits))
	for i := range paintKits {
		pk := &paintKits[i]
		idx[pk.Name] = pk
	}
	return idx
}

type itemSetIndex struct {
	paintKits map[string][]string                              // weapon_class → []paint_kit_name
	combo     map[string]*models.PaintKitWeaponCombinationData // "weapon_class|paint_kit_name" → combo
}

func buildItemSetIndex(itemSets []models.ItemSet) *itemSetIndex {
	idx := &itemSetIndex{
		paintKits: make(map[string][]string),
		combo:     make(map[string]*models.PaintKitWeaponCombinationData),
	}
	for i := range itemSets {
		is := &itemSets[i]
		data := &models.PaintKitWeaponCombinationData{
			ItemSetId:     is.Key,
			CanBeStatTrak: is.HasCrate,
			CanBeSouvenir: is.HasSouvenir,
		}
		for _, item := range is.Items {
			idx.paintKits[item.WeaponClass] = append(idx.paintKits[item.WeaponClass], item.PaintKitName)
			idx.combo[item.WeaponClass+"|"+item.PaintKitName] = data
		}
	}
	return idx
}
