package parsers

import (
	"context"
	"fmt"
	"slices"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

type StickerTypeParams struct {
	TournamentPlayerId int `json:"tournament_player_id"`
	TournamentTeamId   int `json:"tournament_team_id"`
	TournamentEventId  int `json:"tournament_event_id"`
}

var sticker_types = []string{
	"normal",
	"player",
	"team",
	"event",
}

var sticker_effects = []string{
	"normal",
	"foil",
	"holo",
	"glitter",
	"gold",
}

var group_id_to_sub_id = map[int]string{
	1: "E",
	2: "P",
	3: "T",
}

func ParseCustomStickers(ctx context.Context, ig *models.ItemsGame, sticker_kits []models.StickerKit, t *modules.Translator) []models.CustomStickers {
	logger := zerolog.Ctx(ctx)

	start := time.Now()

	// Store all custom stickers
	var items []models.CustomStickers

	if len(sticker_kits) == 0 {
		logger.Warn().Msg("No sticker kits found, skipping custom stickers parsing")
		return items
	}

	var unique_tournament_ids []int
	var unique_player_ids []int
	var unique_team_ids []int

	// Handle all the event stickers
	for _, sticker_kit := range sticker_kits {

		// Tournament stickers only
		if sticker_kit.Tournament != nil && sticker_kit.Player == nil {
			tournament_id := sticker_kit.Tournament.Id

			if tournament_id <= 0 {
				logger.Debug().Msgf("Skipping sticker kit with invalid tournament ID: %s", sticker_kit.Name)
				continue // Skip if the tournament ID is invalid
			}

			// We need to get the sticker type, and based on that, we can process it
			current_type := GetStickerType(0, tournament_id, 0)
			current_effect := GetStickerEffect(sticker_kit.StickerMaterial)

			// Get the count of stickers for this event and type
			count := GetStickerCountByTournamentId(&sticker_kits, tournament_id, current_effect, false)

			if count == 0 {
				continue // Skip if no stickers found for the team and type
			}

			generated_id := GenerateCustomStickerId(tournament_id, group_id_to_sub_id[1], &current_effect, &current_type)

			if CustomStickerExists(items, generated_id) {
				continue // Skip if the custom sticker already exists
			}

			items = append(items, models.CustomStickers{
				GeneratedId: generated_id,
				Group:       2,
				Count:       count,
				Name:        modules.GenerateCustomStickerMarketHashName_Event(t, tournament_id, &current_effect),
			})

			if !slices.Contains(unique_tournament_ids, tournament_id) {
				unique_tournament_ids = append(unique_tournament_ids, tournament_id)
			}
		}

		// Team stickers only
		if sticker_kit.Team != nil && sticker_kit.Player == nil {
			team_id := sticker_kit.Team.Id

			if team_id <= 0 {
				logger.Debug().Msgf("Skipping sticker kit with invalid team ID: %s", sticker_kit.Name)
				continue // Skip if the team ID is invalid
			}

			// We need to get the sticker type, and based on that, we can process it
			current_type := GetStickerType(0, 0, team_id)
			current_effect := GetStickerEffect(sticker_kit.StickerMaterial)

			// Get the count of stickers for this event and type
			count := GetStickerCountByTeamId(&sticker_kits, team_id, current_effect, false)

			if count == 0 {
				continue // Skip if no stickers found for the team and type
			}

			generated_id := GenerateCustomStickerId(team_id, group_id_to_sub_id[3], &current_effect, &current_type)

			if CustomStickerExists(items, generated_id) {
				continue // Skip if the custom sticker already exists
			}

			items = append(items, models.CustomStickers{
				GeneratedId: generated_id,
				Group:       3,
				Count:       count,
				Name:        modules.GenerateCustomStickerMarketHashName_Team(t, team_id, &current_effect),
			})

			if !slices.Contains(unique_team_ids, team_id) {
				unique_team_ids = append(unique_team_ids, team_id)
			}
		}

		// Player stickers only
		if sticker_kit.Player != nil {
			player_id := sticker_kit.Player.Id

			if player_id <= 0 {
				logger.Debug().Msgf("Skipping sticker kit with invalid player ID: %s", sticker_kit.Name)
				continue // Skip if the player ID is invalid
			}

			// We need to get the sticker type, and based on that, we can process it
			current_type := GetStickerType(player_id, 0, 0)
			current_effect := GetStickerEffect(sticker_kit.StickerMaterial)

			// Get the count of stickers for this event and type
			count := GetStickerCountByPlayerId(&sticker_kits, player_id, current_effect, false)

			if count == 0 {
				continue // Skip if no stickers found for the player and type
			}

			generated_id := GenerateCustomStickerId(player_id, group_id_to_sub_id[2], &current_effect, &current_type)

			if CustomStickerExists(items, generated_id) {
				continue // Skip if the custom sticker already exists
			}

			items = append(items, models.CustomStickers{
				GeneratedId: generated_id,
				Group:       2,
				Count:       count,
				Name:        modules.GenerateCustomStickerMarketHashName_Player(t, sticker_kit.Player, &current_effect),
			})

			if !slices.Contains(unique_player_ids, player_id) {
				unique_player_ids = append(unique_player_ids, player_id)
			}
		}
	}

	// now we need to get the total per player/team/event
	for _, tournament_id := range unique_tournament_ids {
		curr := GetTotalStickerForSubId(ig, t, sticker_kits, 1, tournament_id)
		items = append(items, curr)
	}

	for _, player_id := range unique_player_ids {
		curr := GetTotalStickerForSubId(ig, t, sticker_kits, 2, player_id)

		items = append(items, curr)
	}

	for _, team_id := range unique_team_ids {
		curr := GetTotalStickerForSubId(ig, t, sticker_kits, 3, team_id)
		items = append(items, curr)
	}

	// Save music kits to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' custom stickers in %s", len(items), duration)

	return items
}

func GetTotalStickerForSubId(ig *models.ItemsGame, t *modules.Translator, sticker_kits []models.StickerKit, group_id int, sub_id int) models.CustomStickers {
	var count int = 0

	// Generate a unique ID for the custom sticker
	generated_id := GenerateCustomStickerId(sub_id, group_id_to_sub_id[group_id], nil, nil)

	for _, kit := range sticker_kits {
		switch group_id {

		// Events
		case 1:
			if kit.Tournament == nil {
				continue // Skip if the sticker kit does not have a tournament
			}

			if kit.Tournament.Id == sub_id {
				count++
			}

		// Players
		case 2:
			if kit.Player == nil {
				continue // Skip if the sticker kit does not have a player
			}

			if kit.Player.Id == sub_id {
				count++
			}

		// Teams
		case 3:
			// Teams are special, we need to ignore ones with player-id
			if kit.Team == nil || kit.Player != nil {
				continue // Skip if the sticker kit does not have a team
			}

			if kit.Team.Id == sub_id {
				count++
			}
		}
	}

	var name string
	switch group_id {
	case 1:
		event := modules.GetTournamentData(t, sub_id)
		name = modules.GenerateCustomStickerMarketHashName_Event(t, event.Id, nil)
	case 2:
		player := modules.GetPlayerByAccountId(ig, sub_id)
		name = modules.GenerateCustomStickerMarketHashName_Player(t, player, nil)
	case 3:
		team := modules.GetTournamentTeamData(t, sub_id)
		name = modules.GenerateCustomStickerMarketHashName_Team(t, team.Id, nil)
	}

	return models.CustomStickers{
		GeneratedId: generated_id,
		Group:       group_id,
		Count:       count,
		Name:        name,
	}
}

func CustomStickerExists(items []models.CustomStickers, generated_id string) bool {
	for _, item := range items {
		if item.GeneratedId == generated_id {
			return true // Found a duplicate
		}
	}
	return false // No duplicates found
}

var sticker_effect_id_suffix = map[string]string{
	"glitter":    "G",
	"holo":       "H",
	"foil":       "F",
	"gold":       "G",
	"lenticular": "L",
	"normal":     "N",
}

var sticker_type_id_suffix = map[string]string{
	"normal": "N",
	"player": "P",
	"team":   "T",
	"event":  "E",
}

func GenerateCustomStickerId(team_id int, _type string, effect *string, sticker_type *string) string {
	if effect == nil || sticker_type == nil {
		return fmt.Sprintf("C%dA%s", team_id, _type)
	}
	return fmt.Sprintf("C%d%s%s%s", team_id, sticker_effect_id_suffix[*effect], sticker_type_id_suffix[*sticker_type], _type)
}

func GetStickerCountByPlayerId(sticker_kits *[]models.StickerKit, player_id int, sticker_effect string, ignore_effect bool) int {
	var count int
	for _, cs := range *sticker_kits {
		if cs.Player == nil {
			continue
		}
		if cs.Player.Id != player_id || (!ignore_effect && cs.Effect != sticker_effect) {
			continue
		}
		count++
	}
	return count
}

func GetStickerCountByTeamId(sticker_kits *[]models.StickerKit, team_id int, sticker_effect string, ignore_effect bool) int {
	var count int
	for _, cs := range *sticker_kits {
		if cs.Team == nil {
			continue
		}
		if cs.Team.Id != team_id || (!ignore_effect && cs.Effect != sticker_effect) {
			continue
		}
		count++
	}
	return count
}

func GetStickerCountByTournamentId(sticker_kits *[]models.StickerKit, tournament_id int, sticker_effect string, ignore_effect bool) int {
	var count int
	for _, cs := range *sticker_kits {
		if cs.Tournament == nil {
			continue
		}
		if cs.Tournament.Id != tournament_id || (!ignore_effect && cs.Effect != sticker_effect) {
			continue
		}
		count++
	}
	return count
}
