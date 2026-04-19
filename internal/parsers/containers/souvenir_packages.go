package containers

import (
	"context"
	"slices"
	"strconv"

	"go-csitems-parser/internal/i18n"
	"go-csitems-parser/internal/itemsgame"
	"go-csitems-parser/internal/marketname"
	"go-csitems-parser/internal/models"
	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/baldurstod/vdf"
	"github.com/rs/zerolog"
)

var souvenirPrefabs = []string{
	"weapon_case_souvenirpkg",
	"aus2025_souvenir_crate_promo_prefab",
}

// SouvenirPackages extracts the souvenir crate items from items_game.txt.
type SouvenirPackages struct{ base.Parser }

func NewSouvenirPackages() *SouvenirPackages {
	return &SouvenirPackages{Parser: base.New("souvenir_packages")}
}

func (s *SouvenirPackages) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	items, err := in.IG.Get("items")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get items")
		return nil, nil
	}
	clientLootLists, _ := in.IG.Get("client_loot_lists")

	var out []models.SouvenirPackage
	defer s.LogCount(ctx, "souvenir packages", func() int { return len(out) })()

	for _, c := range items.GetChilds() {
		prefab, _ := c.GetString("prefab")
		if !slices.Contains(souvenirPrefabs, prefab) {
			continue
		}

		definition_index, _ := strconv.Atoi(c.Key)
		item_name, _ := c.GetString("item_name")
		image_inventory, _ := c.GetString("image_inventory")
		tournament_event_id, _ := itemsgame.GetTournamentEventId(c)
		name, _ := c.GetString("name")

		out = append(out, models.SouvenirPackage{
			DefinitionIndex: definition_index,
			ImageInventory:  image_inventory,
			ItemSetId:       itemsgame.GetContainerItemSet(c, "ItemSet"),
			MarketHashName:  marketname.GenerateMarketHashName(in.T, item_name, nil, "souvenir_package"),
			KeychainSetId:   lookupKeychainSetId(clientLootLists, name),
			Tournament:      i18n.GetTournamentData(in.T, tournament_event_id),
		})
	}

	return out, nil
}

func (s *SouvenirPackages) Commit(in *pipeline.Inputs, result any) {
	if r, ok := result.([]models.SouvenirPackage); ok {
		in.SouvenirPackages = r
	}
}

func lookupKeychainSetId(ig *vdf.KeyValue, name string) *string {
	if ig == nil {
		return nil
	}
	for _, item := range ig.GetChilds() {
		if item.Key != name {
			continue
		}
		kc, err := item.GetString("match_highlight_reel_keychain")
		if err != nil {
			continue
		}
		if kc == "" {
			return nil
		}
		return &kc
	}
	return nil
}
