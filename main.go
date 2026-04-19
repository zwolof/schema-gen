package main

import (
	"context"
	"os"
	"runtime"
	"sync"
	"time"

	"go-csitems-parser/internal/i18n"
	"go-csitems-parser/internal/itemsgame"
	"go-csitems-parser/internal/parsers"
	"go-csitems-parser/internal/parsers/meta"

	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()

	itemsGame := itemsgame.Load("./files/items_game.txt")
	if itemsGame == nil {
		logger.Error().Msg("Failed to load items_game.txt, please check the file path and format.")
		panic("items_game.txt is nil, exiting...")
	}
	logger.Info().Msg("Successfully loaded items_game.txt")

	ctx := logger.WithContext(context.Background())
	factory, err := i18n.Load(ctx, "./files/translations", "english")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to load translations")
		return
	}

	in := &parsers.Inputs{
		IG:                itemsGame,
		T:                 factory.Get("English"),
		KnifeSkinMap:      itemsgame.LoadKnifeSkinsMap("./files/knife_skins.json"),
		SkinRarityMap:     meta.SkinWeaponRarityMap(ctx, itemsGame),
		StickerItemSetMap: meta.StickerItemSetMap(ctx, itemsGame),
	}

	start := time.Now()
	results, err := parsers.Default.Run(ctx, in, runtime.NumCPU())
	if err != nil {
		logger.Error().Err(err).Msg("pipeline failed")
	}
	logger.Debug().Msgf("[go-items] Parsed all items in %s", time.Since(start))

	exportStart := time.Now()
	var wg sync.WaitGroup
	for _, e := range parsers.Default.Exports(results) {
		wg.Go(func() {
			ExportToJsonFile(e.Value, e.Name)
		})
	}
	wg.Wait()
	logger.Debug().Msgf("[go-items] Exported all files in %s", time.Since(exportStart))
}
