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

	player_agents := parsers.ParseAgents(ctx, itemsGame, t)
	souvenir_packages := parsers.ParseSouvenirPackages(ctx, itemsGame, t)
	musicKits := parsers.ParseMusicKits(ctx, itemsGame, t)
	collectibles := parsers.ParseCollectibles(ctx, itemsGame, t)
	weapon_cases := parsers.ParseWeaponCases(ctx, itemsGame, t)
	rarities := parsers.ParseRarities(ctx, itemsGame, t)
	keychains := parsers.ParseKeychains(ctx, itemsGame, t)
	weapons := parsers.ParseWeapons(ctx, itemsGame, t)
	gloves := parsers.ParseGloves(ctx, itemsGame, t)
	knives := parsers.ParseKnives(ctx, itemsGame, t)
	highlight_reels := parsers.ParseHighlightReels(ctx, itemsGame, t)
	sticker_capsules := parsers.ParseStickerCapsules(ctx, itemsGame, t)
	misc_capsules := parsers.ParseSelfOpeningCrates(ctx, itemsGame, t)
	paint_kits := parsers.ParsePaintKits(ctx, itemsGame, t)
	sticker_kits := parsers.ParseStickerKits(ctx, itemsGame, t)
	custom_stickers := parsers.ParseCustomStickers(ctx, itemsGame, sticker_kits, t)
	item_sets := parsers.ParseItemSets(ctx, itemsGame, souvenir_packages, weapon_cases, t)
	armory_rewards := parsers.ParseArmoryRewards(ctx, itemsGame, &item_sets, t)
	collections := parsers.ParseCollections(ctx, itemsGame, souvenir_packages, weapon_cases, t)

	logger.Debug().Msgf("[go-items] Parsed all items in %s", time.Since(start))

	knife_skin_map := modules.LoadKnifeSkinsMap("./files/knife_skins.json")
	knife_skins := modules.GetKnifePaintKits(&knives, &paint_kits, knife_skin_map)
	weapon_skins := modules.GetWeaponPaintKits(&weapons, &paint_kits, &item_sets)
	glove_skins := modules.GetGlovePaintKits(&gloves, &paint_kits, knife_skin_map)

	ExportToJsonFile(player_agents, "agents")
	ExportToJsonFile(souvenir_packages, "souvenir_packages")
	ExportToJsonFile(musicKits, "music_kits")
	ExportToJsonFile(collectibles, "collectibles")
	ExportToJsonFile(weapon_cases, "weapon_cases")
	ExportToJsonFile(rarities, "rarities")
	ExportToJsonFile(keychains, "keychains")
	ExportToJsonFile(weapons, "weapons")
	ExportToJsonFile(highlight_reels, "highlight_reels")
	ExportToJsonFile(sticker_capsules, "sticker_capsules")
	ExportToJsonFile(misc_capsules, "misc_capsules")
	ExportToJsonFile(paint_kits, "paint_kits")
	ExportToJsonFile(sticker_kits, "sticker_kits")
	ExportToJsonFile(custom_stickers, "custom_stickers")
	ExportToJsonFile(armory_rewards, "armory_rewards")
	ExportToJsonFile(collections, "collections")
	ExportToJsonFile(knife_skins, "knife_skins")
	ExportToJsonFile(weapon_skins, "weapon_skins")
	ExportToJsonFile(glove_skins, "glove_skins")

	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}
