package models

import "github.com/baldurstod/vdf"

// ItemsGame wraps the parsed items_game.txt VDF tree.
type ItemsGame struct {
	*vdf.KeyValue
}

// Localization is a name + description pair loaded from csgo_<lang>.txt.
type Localization struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
