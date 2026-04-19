package models

// StickerKit is a single sticker variant extracted from sticker_kits.
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
}

// CustomStickers is an aggregated per-event/team/player sticker entry.
type CustomStickers struct {
	GeneratedId string `json:"generated_id"`
	Group       int    `json:"group"`
	Name        string `json:"name"`
	Count       int    `json:"count"`
}

// SchemaCustomSticker is the public schema shape for a CustomStickers entry.
type SchemaCustomSticker struct {
	Group int    `json:"group"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// StickerCapsule is a sticker-pack container (crate_sticker_pack_*,
// crate_signature_pack_*, etc.).
type StickerCapsule struct {
	DefinitionIndex int     `json:"definition_index"`
	Name            string  `json:"name"`
	MarketHashName  string  `json:"market_hash_name"`
	ItemDescription string  `json:"item_description"`
	ImageInventory  string  `json:"image"`
	ItemSetId       *string `json:"item_set_id"`
}
