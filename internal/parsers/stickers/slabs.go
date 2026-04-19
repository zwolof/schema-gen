package stickers

import (
	"context"
	"fmt"
	"time"

	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/rs/zerolog"
)

// Slabs returns a map of sticker definition_index → CDN image URL.
// URL format: econ/stickers/{sticker_material}_1355_37
type Slabs struct{ base.Parser }

func NewSlabs() *Slabs { return &Slabs{Parser: base.New("sticker_slabs")} }

func (s *Slabs) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)
	start := time.Now()

	out := make(map[int]string, len(in.StickerKits))
	for _, kit := range in.StickerKits {
		if kit.StickerMaterial == "" {
			continue
		}
		out[kit.DefinitionIndex] = fmt.Sprintf("econ/stickers/%s_1355_37", kit.StickerMaterial)
	}

	logger.Info().Msgf("Built '%d' sticker slab URLs in %s", len(out), time.Since(start))
	return out, nil
}
