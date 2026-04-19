package cosmetics

import (
	"context"
	"strconv"

	"go-csitems-parser/internal/marketname"
	"go-csitems-parser/internal/models"
	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/rs/zerolog"
)

// Agents extracts player agent (customplayertradable) items.
type Agents struct{ base.Parser }

func NewAgents() *Agents { return &Agents{Parser: base.New("agents")} }

func (a *Agents) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	items, err := in.IG.Get("items")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get items, is items_game.txt valid?")
		return nil, nil
	}

	var out []models.PlayerAgent
	defer a.LogCount(ctx, "agents", func() int { return len(out) })()

	for _, agent := range items.GetChilds() {
		prefab, _ := agent.GetString("prefab")
		if prefab != "customplayertradable" {
			continue
		}

		definition_index, _ := strconv.Atoi(agent.Key)
		item_name, _ := agent.GetString("item_name")
		item_rarity, _ := agent.GetString("item_rarity")
		image_inventory, _ := agent.GetString("image_inventory")

		out = append(out, models.PlayerAgent{
			DefinitionIndex: definition_index,
			MarketHashName:  marketname.GenerateMarketHashName(in.T, item_name, nil, "agent"),
			ImageInventory:  image_inventory,
			Rarity:          item_rarity,
		})
	}

	return out, nil
}
