// Package pipeline holds the infrastructure types shared by every parser:
// the Parser interface, the Inputs struct that threads state between tiers,
// and the Tier / Pipeline types that drive parallel tier dispatch.
//
// It lives in its own package so concrete parser packages can import it
// without pulling in the registry (which imports the concrete packages).
package pipeline

import (
	"context"
	"runtime"
	"sync"

	"go-csitems-parser/internal/i18n"
	"go-csitems-parser/internal/models"

	"golang.org/x/sync/errgroup"
)

// Inputs carries everything any parser in the pipeline might read. Raw inputs
// (IG, T, KnifeSkinMap, the pre-computed index maps) are populated before the
// pipeline runs; per-tier outputs are written back via each parser's Commit
// method before the next tier starts.
type Inputs struct {
	IG           *models.ItemsGame
	T            i18n.Translator
	KnifeSkinMap map[string][]string

	SkinRarityMap     map[string]string
	StickerItemSetMap map[string]string

	SouvenirPackages []models.SouvenirPackage
	WeaponCases      []models.WeaponCase
	Weapons          []models.BaseWeapon
	Knives           []models.BaseWeapon
	Gloves           []models.BaseWeapon
	PaintKits        []models.PaintKit

	StickerKits []models.StickerKit
	ItemSets    []models.ItemSet
}

// Parser is the contract every concrete parser satisfies. The usual route is
// to embed base.Parser (for Name / Exported / Commit / timing) and define a
// Parse method; override Commit on the concrete type when the output feeds
// later tiers.
type Parser interface {
	Name() string
	Parse(ctx context.Context, in *Inputs) (any, error)
	Commit(in *Inputs, result any)
	Exported() bool
}

// Tier bundles parsers that run concurrently. AfterCommit, if set, runs once
// every parser in the tier has finished and committed.
type Tier struct {
	Parsers     []Parser
	AfterCommit func(*Inputs)
}

// Pipeline is an ordered list of tiers. Tiers run in sequence; parsers within
// a tier run in parallel.
type Pipeline []Tier

// Run executes the pipeline, bounding per-tier concurrency at limit (defaults
// to runtime.NumCPU). The returned map is keyed by parser Name.
func (pl Pipeline) Run(ctx context.Context, in *Inputs, limit int) (map[string]any, error) {
	if limit <= 0 {
		limit = runtime.NumCPU()
	}

	results := make(map[string]any)
	var mu sync.Mutex

	for _, tier := range pl {
		g, gctx := errgroup.WithContext(ctx)
		g.SetLimit(limit)
		for _, p := range tier.Parsers {
			g.Go(func() error {
				r, err := p.Parse(gctx, in)
				mu.Lock()
				results[p.Name()] = r
				mu.Unlock()
				return err
			})
		}
		if err := g.Wait(); err != nil {
			return results, err
		}
		for _, p := range tier.Parsers {
			p.Commit(in, results[p.Name()])
		}
		if tier.AfterCommit != nil {
			tier.AfterCommit(in)
		}
	}
	return results, nil
}

// Export bundles a parser result with its export name for the writer step.
type Export struct {
	Name  string
	Value any
}

// Exports returns the subset of results that should be written to disk, in
// the order parsers were registered.
func (pl Pipeline) Exports(results map[string]any) []Export {
	var out []Export
	for _, tier := range pl {
		for _, p := range tier.Parsers {
			if !p.Exported() {
				continue
			}
			out = append(out, Export{Name: p.Name(), Value: results[p.Name()]})
		}
	}
	return out
}
