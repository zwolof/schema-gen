package models

type StickerCapsule struct {
	DefinitionIndex int     `json:"definition_index"`
	Name            string  `json:"name"`
	MarketHashName  string  `json:"market_hash_name"`
	ItemDescription string  `json:"item_description"`
	ImageInventory  string  `json:"image"`
	ItemSetId       *string `json:"item_set_id"`
}
