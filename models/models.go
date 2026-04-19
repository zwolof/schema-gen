package models

import "github.com/baldurstod/vdf"

type ItemsGame struct {
	*vdf.KeyValue
}

type Localization struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type BaseWeapon struct {
	DefinitionIndex int    `json:"definition_index"`
	Name            string `json:"name"`
	Description	 string `json:"description"`
	ClassName       string `json:"classname"`
	NumStickers     int    `json:"sticker_count"`
	ImageInventory  string `json:"image"`

	// ItemName        string `json:"item_name"`
	// ItemDescription string `json:"item_description"`
	// Slot            string `json:"slot"`
}

type GenericColor struct {
	Key       string `json:"key"`
	ColorName string `json:"color_name"`
	HexColor  string `json:"hex_color"`
}

type CustomStickers struct {
	GeneratedId string `json:"generated_id"`
	Group       int    `json:"group"`
	Name        string `json:"name"`
	Count       int    `json:"count"`
	// Count []*vdf.KeyValue `json:"count"`
}

type StickerKit struct {
	DefinitionIndex int             `json:"definition_index"`
	MarketHashName  string          `json:"market_hash_name"`
	Name            string          `json:"name"`
	StickerMaterial string          `json:"sticker_material"`
	Image           string          `json:"image"`
	Rarity          string          `json:"rarity"`
	Effect          string          `json:"effect"`
	Type            string          `json:"type"`
	ItemSetId       string          `json:"item_set_id,omitempty"`
	Tournament      *TournamentData `json:"tournament"`
	Team            *TournamentData `json:"team"`
	Player          *TournamentData `json:"player"`

	// ItemName          string        `json:"item_name"`
	// DescriptionString string        `json:"description_string"`
	// TournamentEventId int           `json:"tournament_event_id"`
	// TournamentTeamId  int           `json:"tournament_team_id"`
}

type PaintKitWearRange struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
}

type PaintKit struct {
	DefinitionIndex   int               `json:"definition_index"`
	Name              string            `json:"name"`
	Description	   string            `json:"description"`
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

type ItemSetItem struct {
	PaintKitName string `json:"paintkit"`
	WeaponClass  string `json:"weapon"`
}

type LootListItem struct {
	Name  string `json:"item_name"`
	Class string `json:"item_class"`
}

type ClientLootList struct {
	LootListId   string                  `json:"loot_list_id"`
	Series       int                     `json:"series"`
	SubLootLists []ClientLootListSubList `json:"sub_loot_lists"`
}

type ClientLootListSubList struct {
	Rarity       string         `json:"rarity"`
	LootListName string         `json:"loot_list_name"`
	Items        []LootListItem `json:"items"`
}

type ItemSet struct {
	Key         string        `json:"key"`
	Name        string        `json:"name"`
	Type        ItemSetType   `json:"type"`
	Items       []ItemSetItem `json:"items"`
	Agents      []string      `json:"agents"`
	HasCrate    bool          `json:"has_crate"`
	HasSouvenir bool          `json:"has_souvenir"`
}

type OperationPointsRedeemableItem struct {
	Points         int    `json:"points"`
	Name            string `json:"name"`
	ItemSetId      string `json:"item_set_id"`
	UIImageThumbnail  string `json:"ui_image_thumbnail"`
	UIOrder        int `json:"order"`
}

type Rarity struct {
	LocRarity    string `json:"loc_rarity"`
	LocWeapon    string `json:"loc_weapon"`
	LocCharacter string `json:"loc_character"`
	Hex          string `json:"hex"`
}

type MusicKit struct {
	DefinitionIndex int    `json:"definition_index"`
	Name            string `json:"name"`
	ImageInventory  string `json:"image"`
	MarketHashName  string `json:"market_hash_name"`

	ItemName string `json:"item_name"`
	Model    string `json:"display_model"`
}

type Keychain struct {
	DefinitionIndex int    `json:"definition_index"`
	Name            string `json:"name"`
	MarketHashName  string `json:"market_hash_name"`
	LocDescription  string `json:"loc_description"`

	Rarity          string `json:"rarity"`
	Quality         string `json:"quality"`
	ImageInventory  string `json:"image"`
	Model           string `json:"display_model"`
	IsCommodity       bool   `json:"is_commodity"`
	IsSpecialCharm	bool   `json:"is_special_charm"`

	LootListId      string `json:"loot_list_id"`
}

type HighlightReelTeams struct {
	TeamZero string `json:"team_0"`
	TeamOne  string `json:"team_1"`
}

type HighlightReel struct {
	Id             string             `json:"id"`
	DefinitionIndex int                `json:"definition_index"`
	MarketHashName string             `json:"market_hash_name"`
	ReelTitle	   string             `json:"reel_title"`
	ReelDescription string             `json:"reel_description"`
	Map            string             `json:"map"`
	Teams          HighlightReelTeams `json:"teams"`
	Tournament     *TournamentData    `json:"tournament"`
	Stage          *TournamentData    `json:"stage"`
	VideoUrl	   string             `json:"video_url"`
}

type PlayerAgent struct {
	DefinitionIndex int    `json:"definition_index"`
	MarketHashName  string `json:"market_hash_name"`
	ImageInventory  string `json:"image"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Rarity          string `json:"rarity"`
}

type WeaponCaseItemSet struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type WeaponCaseKey struct {
	DefinitionIndex int    `json:"definition_index"`
	Name            string `json:"name"`
	ItemName        string `json:"item_name"`
	ItemDescription string `json:"item_description"`
	FirstSaleDate   string `json:"first_sale_date"`
	Prefab          string `json:"prefab"`
	ImageInventory  string `json:"image"`
}

type WeaponCase struct {
	DefinitionIndex int            `json:"definition_index"`
	Name            string         `json:"name"`
	MarketHashName  string         `json:"market_hash_name"`
	ImageInventory  string         `json:"image"`
	ItemSetId       *string        `json:"item_set_id"`
	Key             *WeaponCaseKey `json:"key"`

	Description   string `json:"description"`
	Prefab        string `json:"prefab"`
	Model         string `json:"model_player"`
	FirstSaleDate string `json:"first_sale_date"`
}

type SouvenirPackage struct {
	DefinitionIndex int             `json:"definition_index"`
	MarketHashName  string          `json:"market_hash_name"`
	ImageInventory  string          `json:"image"`
	KeychainSetId   *string         `json:"keychain_set_id"`
	ItemSetId       *string         `json:"item_set_id"`
	Tournament      *TournamentData `json:"tournament"`

	// Name              string             `json:"name"`
	// ItemName          string             `json:"item_name"`
	// ItemDescription   string             `json:"item_description"`
}

type TournamentData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Collectible struct {
	DefinitionIndex   int    `json:"definition_index"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	Model             string `json:"display_model"`
	Prefab            string `json:"prefab"`
	Description       string `json:"description"`
	Rarity            string `json:"rarity"`
	ImageInventory    string `json:"image"`
	TournamentEventId int    `json:"tournament_event_id"`
	MarketHashName    string `json:"market_hash_name"`
	// Genuine is true when this pin also exists as an "attendance_pin" —
	// meaning it was physically distributed at an event and can appear
	// with the Genuine quality in-game.
	Genuine           bool   `json:"genuine,omitempty"`
}

type Collection struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	HasCrate    bool   `json:"has_crate"`
	HasSouvenir bool   `json:"has_souvenir"`
	Image       string `json:"image"`
}

type ItemSetCollectionMap struct {
	Key         string        `json:"key"`
	Name        string        `json:"name"`
	Items       []ItemSetItem `json:"items"`
	HasCrate    bool          `json:"has_crate"`
	HasSouvenir bool          `json:"has_souvenir"`
}

type PaintKitWeaponCombinationData struct {
	ItemSetId     string `json:"item_set_id"`
	CanBeStatTrak bool   `json:"can_be_stattrak"`
	CanBeSouvenir bool   `json:"can_be_souvenir"`
}

type SchemaGenericeMap struct {
	MarketHashName string `json:"market_hash_name"`
	Rarity         string `json:"rarity"`
	Image          string `json:"image"`
}

type SchemaCustomSticker struct {
	Group int    `json:"group"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

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