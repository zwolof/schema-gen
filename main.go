package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"
	"go-csitems-parser/modules/parsers"

	"github.com/rs/zerolog"
)

type ItemSchema struct {
	Collections    []models.Collection                   `json:"collections"`
	Rarities       []models.SchemaRarity                 `json:"rarities"`
	Stickers       map[int]modules.SchemaItemWithImage   `json:"stickers"`
	Keychains      map[int]modules.SchemaItemWithImage   `json:"keychains"`
	Collectibles   map[int]models.SchemaGenericeMap      `json:"collectibles"`
	Containers     map[int]string                        `json:"containers"`
	Agents         map[int]models.SchemaGenericeMap      `json:"agents"`
	CustomStickers map[string]models.SchemaCustomSticker `json:"custom_stickers"`
	MusicKits      map[int]models.SchemaGenericeMap      `json:"music_kits"`
	Weapons        map[int]modules.SchemaWeaponSkinMap   `json:"weapons"`

	HighlightReels []models.HighlightReel `json:"highlight_reels"`
}

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()

	itemsGame := modules.LoadItemsGame("./files/items_game.txt")
	if itemsGame == nil {
		logger.Error().Msg("Failed to load items_game.txt, please check the file path and format.")
		panic("items_game.txt is nil, exiting...")
	}
	logger.Info().Msg("Successfully loaded items_game.txt")

	ctx := logger.WithContext(context.Background())

	factory := modules.LoadAllTranslations(ctx, "./files/translations")
	if factory == nil {
		logger.Error().Msg("Failed to load translations")
		return
	}

	t := factory.GetTranslator("English")
	start := time.Now()

	// --- Phase 1: independent parsers run concurrently ---
	var (
		player_agents     []models.PlayerAgent
		souvenir_packages []models.SouvenirPackage
		musicKits         []models.MusicKit
		collectibles      []models.Collectible
		weapon_cases      []models.WeaponCase
		rarities          []models.Rarity
		keychains         []models.Keychain
		weapons           []models.BaseWeapon
		gloves            []models.BaseWeapon
		knives            []models.BaseWeapon
		highlight_reels   []models.HighlightReel
		sticker_capsules  []models.StickerCapsule
		misc_capsules     []models.StickerCapsule
		sticker_kits      []models.StickerKit
		paint_kits        []models.PaintKit
	)

	var wg sync.WaitGroup
	wg.Add(14)
	go func() { defer wg.Done(); player_agents = parsers.ParseAgents(ctx, itemsGame, t) }()
	go func() { defer wg.Done(); souvenir_packages = parsers.ParseSouvenirPackages(ctx, itemsGame, t) }()
	go func() { defer wg.Done(); musicKits = parsers.ParseMusicKits(ctx, itemsGame, t) }()
	go func() { defer wg.Done(); collectibles = parsers.ParseCollectibles(ctx, itemsGame, t) }()
	go func() { defer wg.Done(); weapon_cases = parsers.ParseWeaponCases(ctx, itemsGame, t) }()
	go func() { defer wg.Done(); rarities = parsers.ParseRarities(ctx, itemsGame, t) }()
	go func() { defer wg.Done(); keychains = parsers.ParseKeychains(ctx, itemsGame, t) }()
	go func() { defer wg.Done(); weapons = parsers.ParseWeapons(ctx, itemsGame, t) }()
	go func() { defer wg.Done(); gloves = parsers.ParseGloves(ctx, itemsGame, t) }()
	go func() { defer wg.Done(); knives = parsers.ParseKnives(ctx, itemsGame, t) }()
	go func() { defer wg.Done(); highlight_reels = parsers.ParseHighlightReels(ctx, itemsGame, t) }()
	go func() { defer wg.Done(); sticker_capsules = parsers.ParseStickerCapsules(ctx, itemsGame, t) }()
	go func() { defer wg.Done(); misc_capsules = parsers.ParseSelfOpeningCrates(ctx, itemsGame, t) }()
	go func() { defer wg.Done(); paint_kits = parsers.ParsePaintKits(ctx, itemsGame, t) }()
	wg.Wait()

	// sticker_kits must finish before custom_stickers
	sticker_kits = parsers.ParseStickerKits(ctx, itemsGame, t)
	custom_stickers := parsers.ParseCustomStickers(ctx, itemsGame, sticker_kits, t)

	// item_sets depends on souvenir_packages and weapon_cases
	item_sets := parsers.ParseItemSets(ctx, itemsGame, souvenir_packages, weapon_cases, t)

	logger.Debug().Msgf("[go-items] Parsed all items in %s", time.Since(start))

	// --- Phase 2: mapping (depends on parsed results) ---
	knife_skin_map := modules.LoadKnifeSkinsMap("./files/knife_skins.json")
	knife_skins := modules.GetKnifePaintKits(&knives, &paint_kits, knife_skin_map)
	weapon_skins := modules.GetWeaponPaintKits(&weapons, &paint_kits, &item_sets)
	glove_skins := modules.GetGlovePaintKits(&gloves, &paint_kits, knife_skin_map)

	// Paint kits need to be exported after item_sets are resolved
	collections := parsers.ParseCollections(ctx, itemsGame, souvenir_packages, weapon_cases, t)

	// --- Phase 3: export concurrently ---
	exports := []struct {
		v    any
		name string
	}{
		{paint_kits, "paint_kits"},
		{weapons, "weapons"},
		{sticker_capsules, "sticker_capsules"},
		{custom_stickers, "custom_stickers"},
		{sticker_kits, "sticker_kits"},
		{misc_capsules, "misc_capsules"},
		{knife_skins, "knife_skins"},
		{weapon_skins, "weapon_skins"},
		{glove_skins, "glove_skins"},
		{collections, "collections"},
		{rarities, "rarities"},
		{keychains, "keychains"},
		{player_agents, "agents"},
		{musicKits, "music_kits"},
		{collectibles, "collectibles"},
		{souvenir_packages, "souvenir_packages"},
		{weapon_cases, "weapon_cases"},
		{highlight_reels, "highlight_reels"},
	}
	var exportWg sync.WaitGroup
	exportWg.Add(len(exports))
	for _, e := range exports {
		e := e
		go func() { defer exportWg.Done(); ExportToJsonFile(e.v, e.name) }()
	}
	exportWg.Wait()

	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}
