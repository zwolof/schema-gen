package parsers

import (
	"context"
	"slices"
	"strconv"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/baldurstod/vdf"
	"github.com/rs/zerolog"
)

var valid_prefabs = []string{
	"weapon_case_souvenirpkg",
	"aus2025_souvenir_crate_promo_prefab",
}

func ParseSouvenirPackages(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.SouvenirPackage {
	logger := zerolog.Ctx(ctx)

	start := time.Now()

	items, err := ig.Get("items")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get items")
		return nil
	}
	client_loot_lists, _ := ig.Get("client_loot_lists")

	var souvenir_packages []models.SouvenirPackage
	for _, c := range items.GetChilds() {
		prefab, _ := c.GetString("prefab")

		if !slices.Contains(valid_prefabs, prefab) {
			continue
		}

		definition_index, _ := strconv.Atoi(c.Key)
		item_name, _ := c.GetString("item_name")
		image_inventory, _ := c.GetString("image_inventory")

		item_set := modules.GetContainerItemSet(c, t, "ItemSet")
		tournament_event_id, _ := modules.GetTournamentEventId(c)

		name, _ := c.GetString("name")

		current := models.SouvenirPackage{
			DefinitionIndex: definition_index,
			ImageInventory:  image_inventory,
			ItemSetId:       item_set,
			MarketHashName:  modules.GenerateMarketHashName(t, item_name, nil, "souvenir_package"),
			KeychainSetId:   GetKeychainSetId(client_loot_lists, name),
			Tournament:      modules.GetTournamentData(t, tournament_event_id),
		}

		souvenir_packages = append(souvenir_packages, current)
	}

	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' souvenir packages in %s", len(souvenir_packages), duration)

	return souvenir_packages
}

func GetKeychainSetId(ig *vdf.KeyValue, name string) *string {
	var kc_capsule_id string

	for _, item := range ig.GetChilds() {
		if item.Key != name {
			continue
		}

		match_highlight_reel_keychain, err := item.GetString("match_highlight_reel_keychain")

		// If this key does not exist, skip it
		if err != nil {
			continue
		}

		kc_capsule_id = match_highlight_reel_keychain
		break
	}

	if kc_capsule_id == "" {
		return nil
	}

	return &kc_capsule_id
}
