// Package base provides the embeddable Parser foundation. Concrete parsers
// embed base.Parser to inherit Name, Exported, a no-op Commit, and the
// LogCount timing helper — then add their own Parse method (and optionally
// override Commit when the output feeds later tiers).
package base

import (
	"context"
	"time"

	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/rs/zerolog"
)

// Parser is the shared foundation every concrete parser embeds. Use New for
// parsers whose output is written to exported/ and Internal for those
// consumed only by later tiers.
type Parser struct {
	name     string
	exported bool
}

// New builds a Parser whose output should be written to exported/.
func New(name string) Parser { return Parser{name: name, exported: true} }

// Internal builds a Parser whose output is consumed only by later tiers.
func Internal(name string) Parser { return Parser{name: name, exported: false} }

func (p Parser) Name() string   { return p.name }
func (p Parser) Exported() bool { return p.exported }

// Commit is the default no-op. Concrete parsers override this to publish
// state back into *Inputs for later tiers to read.
func (p Parser) Commit(in *pipeline.Inputs, result any) {}

// LogCount returns a deferred logger. Canonical usage:
//
//	var out []models.Foo
//	defer p.LogCount(ctx, "foos", func() int { return len(out) })()
func (p Parser) LogCount(ctx context.Context, label string, count func() int) func() {
	start := time.Now()
	return func() {
		zerolog.Ctx(ctx).Info().Msgf("Parsed '%d' %s in %s", count(), label, time.Since(start))
	}
}
