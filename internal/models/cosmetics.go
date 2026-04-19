package models

// PlayerAgent is a customplayertradable item (agent skin).
type PlayerAgent struct {
	DefinitionIndex int    `json:"definition_index"`
	MarketHashName  string `json:"market_hash_name"`
	ImageInventory  string `json:"image"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Rarity          string `json:"rarity"`
}

// MusicKit is a music_definitions entry.
type MusicKit struct {
	DefinitionIndex int    `json:"definition_index"`
	Name            string `json:"name"`
	ImageInventory  string `json:"image"`
	MarketHashName  string `json:"market_hash_name"`

	ItemName string `json:"item_name"`
	Model    string `json:"display_model"`
}

// Collectible is a commodity_pin item (service medals, pins, coins).
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
}

// Keychain is a keychain_definitions (charm) entry.
type Keychain struct {
	DefinitionIndex int    `json:"definition_index"`
	Name            string `json:"name"`
	MarketHashName  string `json:"market_hash_name"`
	LocDescription  string `json:"loc_description"`

	Rarity         string `json:"rarity"`
	Quality        string `json:"quality"`
	ImageInventory string `json:"image"`
	Model          string `json:"display_model"`
	IsCommodity    bool   `json:"is_commodity"`
	IsSpecialCharm bool   `json:"is_special_charm"`

	LootListId string `json:"loot_list_id"`
}

// HighlightReel is a per-match highlight reel keychain source.
type HighlightReel struct {
	Id              string             `json:"id"`
	DefinitionIndex int                `json:"definition_index"`
	MarketHashName  string             `json:"market_hash_name"`
	ReelTitle       string             `json:"reel_title"`
	ReelDescription string             `json:"reel_description"`
	Map             string             `json:"map"`
	Teams           HighlightReelTeams `json:"teams"`
	Tournament      *TournamentData    `json:"tournament"`
	Stage           *TournamentData    `json:"stage"`
	VideoUrl        string             `json:"video_url"`
}

// HighlightReelTeams names the two teams in a highlight-reel match.
type HighlightReelTeams struct {
	TeamZero string `json:"team_0"`
	TeamOne  string `json:"team_1"`
}
