package meta

import (
	"context"
	"strings"

	"go-csitems-parser/internal/marketname"
	"go-csitems-parser/internal/models"
	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/rs/zerolog"
)

// Collections extracts item_sets again but for the collection schema
// (filtering out "_characters" agent sets).
type Collections struct{ base.Parser }

func NewCollections() *Collections { return &Collections{Parser: base.New("collections")} }

// svgCollections is the set of collection name tokens (without the "#CSGO_"
// prefix) whose set-icon image is an SVG rather than a PNG.
var svgCollections = map[string]bool{
	"set_timed_drops_achroma":   true,
	"set_timed_drops_exuberant": true,
	"set_community_37":          true,
}

func (c *Collections) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	itemSets, err := in.IG.Get("item_sets")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get item_sets from items_game.txt")
		return nil, nil
	}

	var out []models.Collection
	defer c.LogCount(ctx, "collections", func() int { return len(out) })()

	for _, s := range itemSets.GetChilds() {
		name, _ := s.GetString("name")
		if strings.Contains(name, "_characters") {
			continue
		}

		current := models.Collection{
			Key:         s.Key,
			Name:        marketname.GenerateMarketHashName(in.T, name, nil, "collection"),
			Image:       "econ/set_icons/" + s.Key,
			UseSvgImage: svgCollections[s.Key],
		}

		for _, wpncase := range in.WeaponCases {
			if wpncase.ItemSetId == nil || *wpncase.ItemSetId != current.Key {
				continue
			}
			current.HasCrate = true
			break
		}

		for _, sv := range in.SouvenirPackages {
			for _, id := range sv.ItemSetIds {
				if id == current.Key {
					current.HasSouvenir = true
					break
				}
			}
			if current.HasSouvenir {
				break
			}
		}

		out = append(out, current)
	}

	return out, nil
}
