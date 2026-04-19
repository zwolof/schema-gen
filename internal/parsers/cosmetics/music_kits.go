package cosmetics

import (
	"context"
	"strconv"
	"strings"

	"go-csitems-parser/internal/marketname"
	"go-csitems-parser/internal/models"
	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/rs/zerolog"
)

// MusicKits extracts music_definitions entries.
type MusicKits struct{ base.Parser }

func NewMusicKits() *MusicKits { return &MusicKits{Parser: base.New("music_kits")} }

func (p *MusicKits) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	defs, err := in.IG.Get("music_definitions")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get music_definitions")
		return nil, nil
	}

	var out []models.MusicKit
	defer p.LogCount(ctx, "music-kits", func() int { return len(out) })()

	for _, mk := range defs.GetChilds() {
		definition_index, _ := strconv.Atoi(mk.Key)
		name, _ := mk.GetString("name")
		loc_name, _ := mk.GetString("loc_name")
		image_inventory, _ := mk.GetString("image_inventory")

		if strings.Contains(name, "valve_") {
			continue
		}

		out = append(out, models.MusicKit{
			DefinitionIndex: definition_index,
			Name:            name,
			ImageInventory:  image_inventory,
			MarketHashName:  marketname.GenerateMarketHashName(in.T, loc_name, nil, "music_kit"),
		})
	}

	return out, nil
}
