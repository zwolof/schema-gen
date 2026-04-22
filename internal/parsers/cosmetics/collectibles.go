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

// Collectibles extracts CS2 collectible items (pins, medals, map tokens).
// Items whose prefab is "attendance_pin" were physically distributed at live
// events and exist with the Genuine quality in-game; "commodity_pin" items
// are the regular purchasable counterparts and are never Genuine.
type Collectibles struct{ base.Parser }

func NewCollectibles() *Collectibles { return &Collectibles{Parser: base.New("collectibles")} }

func (p *Collectibles) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	items, err := in.IG.Get("items")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get collectibles from items_game.txt")
		return nil, nil
	}

	prefabs, err := in.IG.Get("prefabs")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get prefabs from items_game.txt")
		return nil, nil
	}

	// Build prefab name → item_quality map so we can resolve any prefab,
	// not just hard-coded "attendance_pin".
	prefabQuality := make(map[string]string)
	for _, pf := range prefabs.GetChilds() {
		if q, _ := pf.GetString("item_quality"); q != "" {
			prefabQuality[pf.Key] = q
		}
	}

	var out []models.Collectible
	defer p.LogCount(ctx, "collectibles", func() int { return len(out) })()

	for _, item := range items.GetChilds() {
		item_name, _ := item.GetString("item_name")
		if item_name == "" || !isItemCollectible(item_name) {
			continue
		}

		definition_index, _ := strconv.Atoi(item.Key)
		image_inventory, _ := item.GetString("image_inventory")
		rarity, _ := item.GetString("item_rarity")
		prefab, _ := item.GetString("prefab")

		if prefab == "map_token" {
			rarity = "ancient"
		}

		out = append(out, models.Collectible{
			DefinitionIndex: definition_index,
			MarketHashName:  marketname.GenerateMarketHashName(in.T, item_name, nil, "collectible"),
			ImageInventory:  image_inventory,
			Rarity:          rarity,
			Genuine:         prefabQuality[prefab] == "genuine",
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
