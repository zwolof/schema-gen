package models

// Rarity is a rarity bucket from the items_game.txt rarities table.
type Rarity struct {
	LocRarity    string `json:"loc_rarity"`
	LocWeapon    string `json:"loc_weapon"`
	LocCharacter string `json:"loc_character"`
	Hex          string `json:"hex"`
}

// GenericColor is an entry from the items_game.txt colors table.
type GenericColor struct {
	Key       string `json:"key"`
	ColorName string `json:"color_name"`
	HexColor  string `json:"hex_color"`
}

// Collection is the exported schema shape of an item_sets entry.
type Collection struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	HasCrate    bool   `json:"has_crate"`
	HasSouvenir bool   `json:"has_souvenir"`
	Image       string `json:"image"`
}

// ItemSet is the internal shape of an item_sets entry with its members.
type ItemSet struct {
	Key         string        `json:"key"`
	Name        string        `json:"name"`
	Type        ItemSetType   `json:"type"`
	Items       []ItemSetItem `json:"items"`
	Agents      []string      `json:"agents"`
	HasCrate    bool          `json:"has_crate"`
	HasSouvenir bool          `json:"has_souvenir"`
}

// ItemSetItem is a (paint-kit, weapon-class) pair inside an ItemSet.
type ItemSetItem struct {
	PaintKitName string `json:"paintkit"`
	WeaponClass  string `json:"weapon"`
}

// ItemSetCollectionMap is a combined ItemSet + Collection export shape.
type ItemSetCollectionMap struct {
	Key         string        `json:"key"`
	Name        string        `json:"name"`
	Items       []ItemSetItem `json:"items"`
	HasCrate    bool          `json:"has_crate"`
	HasSouvenir bool          `json:"has_souvenir"`
}

// TournamentData is an id+name pair used for tournament events, stages, teams
// and pro players.
type TournamentData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// OperationPointsRedeemableItem is a single Operation Armory redeemable reward.
type OperationPointsRedeemableItem struct {
	Points           int    `json:"points"`
	Name             string `json:"name"`
	ItemSetId        string `json:"item_set_id"`
	UIImageThumbnail string `json:"ui_image_thumbnail"`
	UIOrder          int    `json:"order"`
}

// SchemaGenericeMap is a generic (market_hash_name, rarity, image) triple
// used for agents/music_kits/collectibles in the public schema.
type SchemaGenericeMap struct {
	MarketHashName string `json:"market_hash_name"`
	Rarity         string `json:"rarity"`
	Image          string `json:"image"`
}

// LootListItem is one row of a client_loot_lists sub-list.
type LootListItem struct {
	Name  string `json:"item_name"`
	Class string `json:"item_class"`
}

// ClientLootList is the parsed shape of a top-level client_loot_lists entry.
type ClientLootList struct {
	LootListId   string                  `json:"loot_list_id"`
	Series       int                     `json:"series"`
	SubLootLists []ClientLootListSubList `json:"sub_loot_lists"`
}

// ClientLootListSubList is one rarity-suffixed sub-list within a loot list.
type ClientLootListSubList struct {
	Rarity       string         `json:"rarity"`
	LootListName string         `json:"loot_list_name"`
	Items        []LootListItem `json:"items"`
}
