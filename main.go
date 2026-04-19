package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"go-csitems-parser/internal/i18n"
	"go-csitems-parser/internal/itemsgame"
	"go-csitems-parser/internal/parsers"
	"go-csitems-parser/internal/parsers/meta"

	"github.com/rs/zerolog"
)

var (
	cpuProfile = flag.String("cpuprofile", "", "write CPU profile to file (analyse with: go tool pprof -http=:8080 <file>)")
	memProfile = flag.String("memprofile", "", "write heap profile to file (analyse with: go tool pprof -http=:8080 <file>)")
)

func main() {
	flag.Parse()

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()

	if err := execute(&logger); err != nil {
		logger.Error().Err(err).Msg("schema-gen failed")
		os.Exit(1)
	}
}

// execute wires pprof profiling around the pipeline when the relevant flags
// are set. CPU profiling covers exactly the run() call; the heap profile is
// captured after run() finishes (post-GC) so it reflects the final resident
// set rather than mid-pipeline transients.
func execute(logger *zerolog.Logger) error {
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			return fmt.Errorf("create cpu profile: %w", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			return fmt.Errorf("start cpu profile: %w", err)
		}
		logger.Info().Str("file", *cpuProfile).Msg("CPU profiling started")
	}

	runErr := run(logger)

	if *cpuProfile != "" {
		pprof.StopCPUProfile()
		logger.Info().Str("file", *cpuProfile).Msg("CPU profile written")
	}

	if runErr != nil {
		return runErr
	}

	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			return fmt.Errorf("create mem profile: %w", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			return fmt.Errorf("write mem profile: %w", err)
		}
		logger.Info().Str("file", *memProfile).Msg("heap profile written")
	}

	return nil
}

func run(logger *zerolog.Logger) error {
	itemsGame, err := itemsgame.Load("./files/items_game.txt")
	if err != nil {
		return err
	}
	logger.Info().Msg("Successfully loaded items_game.txt")

	ctx := logger.WithContext(context.Background())

	factory, err := i18n.Load(ctx, "./files/translations", "english")
	if err != nil {
		return err
	}

	knifeSkinMap, err := itemsgame.LoadKnifeSkinsMap("./files/knife_skins.json")
	if err != nil {
		return err
	}

	in := &parsers.Inputs{
		IG:                itemsGame,
		T:                 factory.Get("English"),
		KnifeSkinMap:      knifeSkinMap,
		SkinRarityMap:     meta.SkinWeaponRarityMap(ctx, itemsGame),
		StickerItemSetMap: meta.StickerItemSetMap(ctx, itemsGame),
	}

	start := time.Now()
	results, err := parsers.Default.Run(ctx, in, runtime.NumCPU())
	if err != nil {
		return err
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

	return nil
}
