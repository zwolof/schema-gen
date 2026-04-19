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

// Collectibles extracts commodity-pin items (service medals, pins, etc.).
type Collectibles struct{ base.Parser }

func NewCollectibles() *Collectibles { return &Collectibles{Parser: base.New("collectibles")} }

func (p *Collectibles) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	items, err := in.IG.Get("items")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get collectibles from items_game.txt")
		return nil, nil
	}

	var out []models.Collectible
	defer p.LogCount(ctx, "collectibles", func() int { return len(out) })()

	for _, item := range items.GetChilds() {
		item_name, _ := item.GetString("item_name")
		if item_name == "" || !isItemCollectible(item_name) {
			continue
		}

		prefab, _ := item.GetString("prefab")
		if prefab != "commodity_pin" {
			continue
		}

		definition_index, _ := strconv.Atoi(item.Key)
		image_inventory, _ := item.GetString("image_inventory")
		rarity, _ := item.GetString("item_rarity")

		out = append(out, models.Collectible{
			DefinitionIndex: definition_index,
			MarketHashName:  marketname.GenerateMarketHashName(in.T, item_name, nil, "collectible"),
			ImageInventory:  image_inventory,
			Rarity:          rarity,
		})
	}

	return out, nil
}

func isItemCollectible(name string) bool {
	if len(name) == 0 {
		return false
	}
	return strings.HasPrefix(name, "#CSGO_Collectible") ||
		strings.HasPrefix(name, "#CSGO_TournamentJournal")
}
