package cosmetics

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"go-csitems-parser/internal/itemsgame"
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

var yearCoinRegex = regexp.MustCompile(`\d+yearcoin`)
var seasonCoinRegex = regexp.MustCompile(`season\d+_coin`)

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
		name, _ := item.GetString("name")
		image_inventory, _ := item.GetString("image_inventory")
		rarity, _ := item.GetString("item_rarity")
		prefab, _ := item.GetString("prefab")

		item_description, _ := item.GetString("item_description")
		description := item_description
		if item_description != "" {
			if translated, err := in.T.GetValueByKey(item_description); err == nil {
				description = translated
			}
		}

		var display_model string
		if attributes, err := item.Get("attributes"); err == nil {
			display_model, _ = attributes.GetString("pedestal display model")
		}

		tournament_event_id, _ := itemsgame.GetTournamentEventId(item)
		if tournament_event_id < 0 {
			tournament_event_id = 0
		}

		if prefab == "map_token" {
			rarity = "ancient"
		}

		out = append(out, models.Collectible{
			DefinitionIndex:   definition_index,
			Name:              name,
			Type:              collectibleType(image_inventory, prefab, item_name),
			Model:             display_model,
			Prefab:            prefab,
			Description:       description,
			Rarity:            rarity,
			ImageInventory:    image_inventory,
			TournamentEventId: tournament_event_id,
			MarketHashName:    marketname.GenerateMarketHashName(in.T, item_name, nil, "collectible"),
			Genuine:           prefabQuality[prefab] == "genuine",
		})
	}

	return out, nil
}

// collectibleType classifies a collectible into a stable snake_case category
// derived from its image, prefab, and item_name.
func collectibleType(image_inventory, prefab, item_name string) string {
	if prefab == "" || image_inventory == "" {
		return "unknown"
	}

	switch {
	case prefab == "premier_season_coin":
		return "premier_season_coin"
	case strings.Contains(image_inventory, "service_medal"):
		return "service_medal"
	case yearCoinRegex.MatchString(image_inventory):
		return "years_of_service"
	case strings.Contains(item_name, "#CSGO_Collectible_Map"):
		return "map_contributor"
	case strings.HasPrefix(item_name, "#CSGO_TournamentJournal"),
		strings.HasPrefix(item_name, "#CSGO_CollectibleCoin"):
		return "pick_em"
	case strings.HasPrefix(item_name, "#CSGO_Collectible_Pin"),
		strings.HasPrefix(item_name, "#CSGO_Collectible_CommunitySeason"):
		return "map_pin"
	case seasonCoinRegex.MatchString(prefab):
		return "operation"
	case prefab == "majors_trophy":
		return "tournament_finalist"
	default:
		return "unknown"
	}
}

func isItemCollectible(name string) bool {
	if len(name) == 0 {
		return false
	}
	return strings.HasPrefix(name, "#CSGO_Collectible") ||
		strings.HasPrefix(name, "#CSGO_TournamentJournal")
}
