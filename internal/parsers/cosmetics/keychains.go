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

// Keychains extracts keychain_definitions (charms).
type Keychains struct{ base.Parser }

func NewKeychains() *Keychains { return &Keychains{Parser: base.New("keychains")} }

func (p *Keychains) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	defs, err := in.IG.Get("keychain_definitions")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get keychain_definitions, is items_game.txt valid?")
		return nil, nil
	}

	var out []models.Keychain
	defer p.LogCount(ctx, "keychains", func() int { return len(out) })()

	for _, mk := range defs.GetChilds() {
		definition_index, _ := strconv.Atoi(mk.Key)
		name, _ := mk.GetString("name")
		loc_name, _ := mk.GetString("loc_name")
		image_inventory, _ := mk.GetString("image_inventory")
		item_rarity, _ := mk.GetString("item_rarity")
		is_commodity, _ := mk.GetBool("is_commodity")
		pedestal_display_model, _ := mk.GetString("pedestal_display_model")

		tags, tagsErr := mk.Get("tags")

		current := models.Keychain{
			DefinitionIndex: definition_index,
			Name:            name,
			MarketHashName:  marketname.GenerateMarketHashName(in.T, loc_name, nil, "keychain"),
			Rarity:          item_rarity,
			ImageInventory:  image_inventory,
			IsCommodity:     is_commodity,
			Model:           pedestal_display_model,
		}

		if tagsErr == nil && tags != nil {
			current.IsSpecialCharm = true
			// image_inventory is already read from items_game above; it's
			// already the canonical path (e.g. "econ/keychains/aus2025/kc_aus2025"),
			// no override needed.
			current.MarketHashName = marketname.GenerateMarketHashName(in.T, loc_name, nil, "kc_sticker_display_case")
		}

		out = append(out, current)
	}

	return out, nil
}
