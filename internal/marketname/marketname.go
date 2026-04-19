// Package marketname builds the market_hash_name strings that appear in the
// exported schema. All generators take an i18n.Translator so they can swap
// in an alternative implementation for tests.
package marketname

import (
	"fmt"
	"strings"

	"go-csitems-parser/internal/i18n"
	"go-csitems-parser/internal/models"
)

var hashNamePrefixes = map[string]string{
	"sticker_kit": "Sticker | ",
	"music_kit":   "Music Kit | ",
	"keychain":    "Charm | ",
}

var stickerEffectNames = map[string]string{
	"glitter":    " (Glitter)",
	"holo":       " (Holo)",
	"foil":       " (Foil)",
	"gold":       " (Gold)",
	"lenticular": " (Lenticular)",
	"normal":     " (Normal)",
}

// dopplerPhaseMap resolves doppler paintkit material names to phase labels
// that get appended to the weapon market hash name.
var dopplerPhaseMap = map[string]string{
	// Doppler phases
	"am_ruby_marbleized":         "Ruby",
	"am_ruby_marbleized_b":       "Ruby",
	"am_sapphire_marbleized":     "Sapphire",
	"am_sapphire_marbleized_b":   "Sapphire",
	"am_blackpearl_marbleized":   "Black Pearl",
	"am_blackpearl_marbleized_b": "Black Pearl",

	"am_doppler_phase1":   "Phase 1",
	"am_doppler_phase2":   "Phase 2",
	"am_doppler_phase3":   "Phase 3",
	"am_doppler_phase4":   "Phase 4",
	"am_doppler_phase1_b": "Phase 1",
	"am_doppler_phase2_b": "Phase 2",
	"am_doppler_phase3_b": "Phase 3",
	"am_doppler_phase4_b": "Phase 4",

	// Gamma Doppler phases
	"am_emerald_marbleized":   "Emerald",
	"am_emerald_marbleized_b": "Emerald",
	"am_gamma_doppler_phase1": "Phase 1",
	"am_gamma_doppler_phase2": "Phase 2",
	"am_gamma_doppler_phase3": "Phase 3",
	"am_gamma_doppler_phase4": "Phase 4",

	// Gamma Doppler Glock phases
	"am_emerald_marbleized_glock":   "Emerald",
	"am_gamma_doppler_phase1_glock": "Phase 1",
	"am_gamma_doppler_phase2_glock": "Phase 2",
	"am_gamma_doppler_phase3_glock": "Phase 3",
	"am_gamma_doppler_phase4_glock": "Phase 4",
}

// GenerateMarketHashName translates name, optionally appends a doppler-phase
// suffix, and prefixes with the item-type tag used downstream.
func GenerateMarketHashName(t i18n.Translator, name string, extra *string, itemType string) string {
	value, err := t.GetValueByKey(name)
	if err != nil {
		fmt.Printf("Error translating name '%s': %v\n", name, err)
		value = name
	}

	if name == "#PaintKit_Default_Tag" {
		value = "Vanilla"
	}
	if name == "kc_sticker_display_case" {
		return value
	}

	if extra != nil && *extra != "" {
		if phase, ok := dopplerPhaseMap[*extra]; ok {
			value = fmt.Sprintf("%s (%s)", value, phase)
		}
	}

	if itemType == "knife" || itemType == "glove" {
		value = fmt.Sprintf("★ %s", value)
	}
	if prefix, ok := hashNamePrefixes[itemType]; ok {
		value = prefix + value
	}

	return value
}

// GenerateCustomStickerMarketHashName_Team builds the per-team sticker name.
func GenerateCustomStickerMarketHashName_Team(t i18n.Translator, teamID int, effect *string) string {
	teamName, _ := t.GetValueByKey(fmt.Sprintf("CSGO_TeamID_%d", teamID))
	if effect == nil {
		return fmt.Sprintf("Sticker | %s", teamName)
	}
	return fmt.Sprintf("Sticker | %s%s", teamName, stickerEffectNames[*effect])
}

// GenerateCustomStickerMarketHashName_Player builds the per-player sticker name.
func GenerateCustomStickerMarketHashName_Player(_ i18n.Translator, player *models.TournamentData, effect *string) string {
	if effect == nil {
		return fmt.Sprintf("Sticker | %s", player.Name)
	}
	return fmt.Sprintf("Sticker | %s%s", player.Name, stickerEffectNames[*effect])
}

// GenerateCustomStickerMarketHashName_Event builds the per-event sticker name.
func GenerateCustomStickerMarketHashName_Event(t i18n.Translator, eventID int, effect *string) string {
	translated, _ := t.GetValueByKey(fmt.Sprintf("CSGO_Tournament_Event_Location_%d", eventID))
	if effect == nil {
		return fmt.Sprintf("Sticker | %s", translated)
	}
	return fmt.Sprintf("Sticker | %s%s", translated, stickerEffectNames[*effect])
}

// GenerateHighlightReelMarketHashName builds the hash name for a highlight-
// reel keychain.
func GenerateHighlightReelMarketHashName(t i18n.Translator, name string, _ int) string {
	value, err := t.GetValueByKey("HighlightReel_" + name)
	if err != nil {
		fmt.Printf("Error translating name '%s': %v\n", name, err)
		value = name
	}

	tournamentID := ""
	if len(name) > 0 {
		parts := strings.Split(name, "_")
		if len(parts) > 0 {
			tournamentID = parts[0]
		}
	}

	capsule, err := t.GetValueByKey("keychain_kc_" + tournamentID)
	if err != nil {
		fmt.Printf("Error translating capsule name '%s': %v\n", tournamentID, err)
		capsule = "Unknown Capsule"
	}

	return fmt.Sprintf("Souvenir Charm | %s | %s", capsule, value)
}

const highlightReelCDNBase = "https://cdn.steamstatic.com/apps/csgo/videos/highlightreels"

// GenerateHighlightReelVideoURL builds the Steam CDN video URL for a
// highlight reel. All IDs are zero-padded to 3 digits. Region is "ww" or "cn".
func GenerateHighlightReelVideoURL(highlightID, mapName string, eventID, stageID, team0ID, team1ID int, region string) string {
	event := fmt.Sprintf("%03d", eventID)
	stage := fmt.Sprintf("%03d", stageID)
	t0 := fmt.Sprintf("%03d", team0ID)
	t1 := fmt.Sprintf("%03d", team1ID)

	matchDir := fmt.Sprintf("%sv%s_%s", t0, t1, stage)
	filename := fmt.Sprintf("%s_%sv%s_%s_%s_%s_%s_1080p.webm", event, t0, t1, stage, mapName, highlightID, region)

	return fmt.Sprintf("%s/%s/%s/%s", highlightReelCDNBase, event, matchDir, filename)
}

// GetSpecialCharmImage returns the CDN path for the hand-curated special
// charm images, or "" if the name doesn't match one.
func GetSpecialCharmImage(name string) string {
	switch name {
	case "kc_aus2025":
		return "econ/keychains/aus2025/kc_aus2025"
	case "kc_bud2025":
		return "econ/keychains/bud2025/kc_bud2025"
	case "kc_sticker_display_case":
		return "econ/keychains/sticker_display_case/kc_sticker_display_case"
	default:
		return ""
	}
}

// GetInventoryImageUrl builds the canonical cs2cdn URL for an inventory image.
func GetInventoryImageUrl(basePath, imageInventory string) string {
	if imageInventory == "" {
		return ""
	}
	return fmt.Sprintf("https://cs2cdn.com/econ/%s/%s.png", basePath, imageInventory)
}

// ItemWear describes a CS2 float wear bracket.
type ItemWear struct {
	Name    string  `json:"name"`
	MinWear float64 `json:"min_wear"`
	MaxWear float64 `json:"max_wear"`
}

// ItemWears is the canonical wear-bracket table.
var ItemWears = map[string]ItemWear{
	"Factory New":    {Name: "Factory New", MinWear: 0.00, MaxWear: 0.07},
	"Minimal Wear":   {Name: "Minimal Wear", MinWear: 0.07, MaxWear: 0.15},
	"Field-Tested":   {Name: "Field-Tested", MinWear: 0.15, MaxWear: 0.38},
	"Well-Worn":      {Name: "Well-Worn", MinWear: 0.38, MaxWear: 0.45},
	"Battle-Scarred": {Name: "Battle-Scarred", MinWear: 0.45, MaxWear: 1.00},
}
