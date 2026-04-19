package meta

import (
	"context"
	"regexp"

	"go-csitems-parser/internal/models"
	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/baldurstod/vdf"
	"github.com/rs/zerolog"
)

var itemSetKeyRegexp = regexp.MustCompile(`^\[(.+?)\](.+)$`)

// ItemSets extracts item_sets and enriches each with HasCrate/HasSouvenir
// derived from the pre-computed Inputs.WeaponCases / SouvenirPackages.
// Publishes into Inputs.ItemSets for Tier-2 parsers.
type ItemSets struct{ base.Parser }

func NewItemSets() *ItemSets { return &ItemSets{Parser: base.Internal("item_sets")} }

func (is *ItemSets) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	raw, err := in.IG.Get("item_sets")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get item_sets from items_game.txt")
		return nil, nil
	}

	var out []models.ItemSet
	defer is.LogCount(ctx, "item sets", func() int { return len(out) })()

	for _, s := range raw.GetChilds() {
		name, _ := s.GetString("name")

		current := models.ItemSet{
			Key:  s.Key,
			Name: name,
			Type: models.ItemSetTypePaintKits,
		}

		itemsetItems, _ := s.Get("items")
		items := itemSetPaintKits(itemsetItems)

		if len(items) == 0 {
			agents := itemSetAgents(itemsetItems)
			if len(agents) > 0 {
				current.Agents = agents
				current.Type = models.ItemSetTypeAgents
			} else {
				continue
			}
		} else {
			current.Items = items
		}

		for _, wpncase := range in.WeaponCases {
			if wpncase.ItemSetId == nil || *wpncase.ItemSetId != current.Key {
				continue
			}
			current.HasCrate = true
			break
		}

		for _, sv := range in.SouvenirPackages {
			if sv.ItemSetId == nil || *sv.ItemSetId != current.Key {
				continue
			}
			current.HasSouvenir = true
			break
		}

		out = append(out, current)
	}

	return out, nil
}

func (is *ItemSets) Commit(in *pipeline.Inputs, result any) {
	if r, ok := result.([]models.ItemSet); ok {
		in.ItemSets = r
	}
}

func itemSetAgents(kv *vdf.KeyValue) []string {
	agents := make([]string, 0)
	for _, item := range kv.GetChilds() {
		agents = append(agents, item.Key)
	}
	return agents
}

func itemSetPaintKits(kv *vdf.KeyValue) []models.ItemSetItem {
	skins := make([]models.ItemSetItem, 0)
	for _, skin := range kv.GetChilds() {
		res := itemSetKeyRegexp.FindStringSubmatch(skin.Key)
		if len(res) < 3 {
			continue
		}
		skins = append(skins, models.ItemSetItem{
			PaintKitName: res[1],
			WeaponClass:  res[2],
		})
	}
	return skins
}
