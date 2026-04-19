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
// Also flags items as Genuine when a matching "attendance_pin" exists — the
// non-commodity variant physically handed out at live events.
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

	// First pass: collect every item_name that has an "attendance_pin" prefab.
	// These pins were physically distributed at live events and therefore
	// exist with the Genuine quality in-game. The commodity_pin counterpart
	// is the regular (non-Genuine) purchasable version of the same pin.
	attendancePins := make(map[string]struct{})
	for _, item := range items.GetChilds() {
		prefab, _ := item.GetString("prefab")
		if prefab != "attendance_pin" {
			continue
		}
		if name, _ := item.GetString("item_name"); name != "" {
			attendancePins[name] = struct{}{}
		}
	}

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

		_, hasGenuine := attendancePins[item_name]

		out = append(out, models.Collectible{
			DefinitionIndex: definition_index,
			MarketHashName:  marketname.GenerateMarketHashName(in.T, item_name, nil, "collectible"),
			ImageInventory:  image_inventory,
			Rarity:          rarity,
			Genuine:         hasGenuine,
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
