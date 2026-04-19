package parsers

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

func ParseCollectibles(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.Collectible {
	logger := zerolog.Ctx(ctx)

	start := time.Now()

	items, err := ig.Get("items")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get collectibles from items_game.txt")
		return nil
	}

	var collectibles []models.Collectible

	// First pass: collect every item_name that has an "attendance_pin" prefab.
	// These are pins that were physically distributed at live events and therefore
	// exist with the Genuine quality in-game.  The commodity_pin counterpart is
	// the regular (non-Genuine) purchasable version of the same pin.
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

		if item_name == "" || !IsItemCollectible(item_name) {
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

		collectibles = append(collectibles, models.Collectible{
			DefinitionIndex: definition_index,
			MarketHashName:  modules.GenerateMarketHashName(t, item_name, nil, "collectible"),
			ImageInventory:  image_inventory,
			Rarity:          rarity,
			Genuine:         hasGenuine,
		})
	}

	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' collectibles in %s", len(collectibles), duration)

	return collectibles
}

func GetCollectibleType(
	image_inventory string,
	prefab string,
	item_name string,
	tournament_event_id int,
) string {
	if prefab == "" {
		return "unknown"
	}

	if image_inventory == "" {
		return "unknown"
	}

	if prefab == "premier_season_coin" {
		return "premier_season_coin"
	}

	if strings.Contains(image_inventory, "service_medal") {
		return "service_medal"
	}

	// image_inventory looks like "10yearcoin", "5yearcoin", etc.
	reg1 := regexp.MustCompile(`\d+yearcoin`)
	if reg1.MatchString(image_inventory) {
		return "years_of_service"
	}

	if strings.Contains(item_name, "#CSGO_Collectible_Map") {
		return "map_contributor"
	}

	if strings.HasPrefix(item_name, "#CSGO_TournamentJournal") || strings.HasPrefix(item_name, "#CSGO_CollectibleCoin") {
		return "pickem"
	}

	if strings.HasPrefix(item_name, "#CSGO_Collectible_Pin") {
		return "collectible_pin"
	}

	// This is a bit odd, idk what Valve was thinking
	if strings.HasPrefix(item_name, "#CSGO_Collectible_CommunitySeason") {
		return "collectible_pin"
	}

	// Create a regex for season1_coin, // season2_coin, etc.
	reg2 := regexp.MustCompile(`season\d+_coin`)
	if reg2.MatchString(prefab) {
		return "operation_coin"
	}

	if prefab == "majors_trophy" {
		return "tournament_trophy"
	}

	// katowice_2014_finalist
	return "unknown"
}

func IsItemCollectible(item_name string) bool {
	if len(item_name) == 0 {
		return false
	}

	//if it starts with "Collectible" or "Collectible_"
	if strings.HasPrefix(item_name, "#CSGO_Collectible") || strings.HasPrefix(item_name, "#CSGO_TournamentJournal") {
		return true
	}

	return false
}

func CanCollectibleBeGenuine(prefab string) bool {
	return false
}