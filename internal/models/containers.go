package models

// WeaponCase is a standard weapon crate.
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

// WeaponCaseKey is the matching key item for a weapon crate.
type WeaponCaseKey struct {
	DefinitionIndex int    `json:"definition_index"`
	Name            string `json:"name"`
	ItemName        string `json:"item_name"`
	ItemDescription string `json:"item_description"`
	FirstSaleDate   string `json:"first_sale_date"`
	Prefab          string `json:"prefab"`
	ImageInventory  string `json:"image"`
}

// WeaponCaseItemSet is a lightweight (id, name) reference to the item set a
// weapon case unlocks.
type WeaponCaseItemSet struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// SouvenirPackage is a tournament souvenir crate.
type SouvenirPackage struct {
	DefinitionIndex int             `json:"definition_index"`
	MarketHashName  string          `json:"market_hash_name"`
	ImageInventory  string          `json:"image"`
	KeychainSetId   *string         `json:"keychain_set_id"`
	ItemSetId       *string         `json:"item_set_id"`
	Tournament      *TournamentData `json:"tournament"`
}
