package models

// BaseWeapon is the shape of a weapon/knife/glove prefab extracted from
// items_game.txt, used as the base item for skin maps.
type BaseWeapon struct {
	DefinitionIndex int    `json:"definition_index"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	ClassName       string `json:"classname"`
	NumStickers     int    `json:"sticker_count"`
	ImageInventory  string `json:"image"`
}

// PaintKitWearRange carries the [min, max] float wear range for a paint kit.
type PaintKitWearRange struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
}

// PaintKit is a weapon paint finish.
type PaintKit struct {
	DefinitionIndex   int               `json:"definition_index"`
	Name              string            `json:"name"`
	MarketHashName    string            `json:"market_hash_name"`
	Wear              PaintKitWearRange `json:"float"`
	Rarity            string            `json:"rarity"`
	Souvenir          bool              `json:"souvenir"`
	StatTrak          bool              `json:"stattrak"`
	ItemSetId         string            `json:"item_set_id,omitempty"`
	UseLegacyModel    bool              `json:"use_legacy_model"`
	DescriptionString string            `json:"description_string"`
	DescriptionTag    string            `json:"description_tag"`
	Style             int               `json:"style"`
}

// PaintKitWeaponCombinationData carries the per-(weapon, paint-kit) flags
// that determine whether the combination can roll StatTrak/Souvenir.
type PaintKitWeaponCombinationData struct {
	ItemSetId     string `json:"item_set_id"`
	CanBeStatTrak bool   `json:"can_be_stattrak"`
	CanBeSouvenir bool   `json:"can_be_souvenir"`
}

// SchemaWeaponPaintKitMap is the per-paint-kit entry in the exported weapon
// skin map.
type SchemaWeaponPaintKitMap struct {
	DefinitionIndex int               `json:"definition_index"`
	Name            string            `json:"name"`
	ItemSetId       string            `json:"item_set_id"`
	Image           string            `json:"image"`
	Rarity          string            `json:"rarity"`
	Float           PaintKitWearRange `json:"float"`
	Souvenir        bool              `json:"souvenir"`
	StatTrak        bool              `json:"stattrak"`
}
