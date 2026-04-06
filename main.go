package main

import (
	"context"
	"fmt"
	"os"
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

	// Set the global logger to use the console writer
	itemsGame := modules.LoadItemsGame("./files/items_game.txt")

	if itemsGame == nil {
		logger.Error().Msg("Failed to load items_game.txt, please check the file path and format.")
		panic("items_game.txt is nil, exiting...")
	} else {
		logger.Info().Msgf("Successfully loaded items_game.txt")
	}

	// Attach the Logger to the context.Context
	ctx := context.Background()
	ctx = logger.WithContext(ctx)

	factory := modules.LoadAllTranslations(ctx, "./files/translations")

	if factory == nil {
		logger.Error().Msg("Failed to load translations")
		return
	}

	translator := factory.GetTranslator("English")
	start := time.Now()

	player_agents := parsers.ParseAgents(ctx, itemsGame, translator)
	souvenir_packages := parsers.ParseSouvenirPackages(ctx, itemsGame, translator)
	musicKits := parsers.ParseMusicKits(ctx, itemsGame, translator)
	collectibles := parsers.ParseCollectibles(ctx, itemsGame, translator)
	weapon_cases := parsers.ParseWeaponCases(ctx, itemsGame, translator)
	rarities := parsers.ParseRarities(ctx, itemsGame, translator)
	keychains := parsers.ParseKeychains(ctx, itemsGame, translator)
	weapons := parsers.ParseWeapons(ctx, itemsGame, translator)
	gloves := parsers.ParseGloves(ctx, itemsGame, translator)
	knives := parsers.ParseKnives(ctx, itemsGame, translator)
	highlight_reels := parsers.ParseHighlightReels(ctx, itemsGame, translator)

	sticker_capsules := parsers.ParseStickerCapsules(ctx, itemsGame, translator)

	misc_capsules := parsers.ParseSelfOpeningCrates(ctx, itemsGame, translator)
	sticker_kits := parsers.ParseStickerKits(ctx, itemsGame, translator)
	custom_stickers := parsers.ParseCustomStickers(ctx, itemsGame, sticker_kits, translator)

	// Paint kits are tricky
	item_sets := parsers.ParseItemSets(ctx, itemsGame, souvenir_packages, weapon_cases, translator)
	paint_kits := parsers.ParsePaintKits(ctx, itemsGame, translator)

	ExportToJsonFile(paint_kits, "paint_kits")
	ExportToJsonFile(weapons, "weapons")
	ExportToJsonFile(sticker_capsules, "sticker_capsules")
	ExportToJsonFile(custom_stickers, "custom_stickers")
	ExportToJsonFile(sticker_kits, "sticker_kits")
	ExportToJsonFile(misc_capsules, "misc_capsules")

	// Special parsing for collections
	collections := parsers.ParseCollections(ctx, itemsGame, souvenir_packages, weapon_cases, translator)
	// Now we need to map whether or not an item has a souvenir variant or not, same for stattrak

	duration := time.Since(start)
	logger.Debug().Msgf("[go-items] Parsed all items in %s", duration)

	// Some knife stuff
	knife_skin_map := modules.LoadKnifeSkinsMap("./files/knife_skins.json")
	knife_skins := modules.GetKnifePaintKits(&knives, &paint_kits, knife_skin_map)
	weapon_skins := modules.GetWeaponPaintKits(&weapons, &paint_kits, &item_sets)
	glove_skins := modules.GetGlovePaintKits(&gloves, &paint_kits, knife_skin_map)

	// Create the final item schema
	ExportToJsonFile(knife_skins, "knife_skins")
	ExportToJsonFile(weapon_skins, "weapon_skins")
	ExportToJsonFile(glove_skins, "glove_skins")
	ExportToJsonFile(collections, "collections")
	ExportToJsonFile(rarities, "rarities")
	ExportToJsonFile(keychains, "keychains")
	ExportToJsonFile(player_agents, "agents")
	ExportToJsonFile(musicKits, "music_kits")
	ExportToJsonFile(collectibles, "collectibles")
	ExportToJsonFile(souvenir_packages, "souvenir_packages")
	ExportToJsonFile(weapon_cases, "weapon_cases")
	ExportToJsonFile(highlight_reels, "highlight_reels")
	
	// keep alive
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}
