package parsers

import (
	"context"
	"fmt"
	"time"

	"go-csitems-parser/models"

	"github.com/rs/zerolog"
)

const stickerCDNBase = "https://cs2cdn.com/econ/stickers"

// ParseStickerSlabs returns a map of sticker definition_index → CDN image URL,
// derived from each kit's sticker_material path.
//
// URL format: https://cs2cdn.com/econ/stickers/{sticker_material}.png
// e.g. sticker_material "cologne2014/teamdignitas_holo"
//      → "https://cs2cdn.com/econ/stickers/cologne2014/teamdignitas_holo.png"
func ParseStickerSlabs(ctx context.Context, kits []models.StickerKit) map[int]string {
	logger := zerolog.Ctx(ctx)
	start := time.Now()

	slabs := make(map[int]string, len(kits))

	for _, kit := range kits {
		if kit.StickerMaterial == "" {
			continue
		}
		slabs[kit.DefinitionIndex] = fmt.Sprintf("econ/stickers/%s_1355_37.png", kit.StickerMaterial)
	}

	logger.Info().Msgf("Built '%d' sticker slab URLs in %s", len(slabs), time.Since(start))

	return slabs
}
