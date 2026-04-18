package parsers

import (
	"context"
	"strconv"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

func ParseHighlightReels(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.HighlightReel {
	logger := zerolog.Ctx(ctx)

	start := time.Now()

	highlight_reels, err := ig.Get("highlight_reels")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get highlight_reels from items_game.txt")
		return nil
	}

	var items []models.HighlightReel
	for _, r := range highlight_reels.GetChilds() {
		id, _ := r.GetString("id")
		tournament_event, _ := r.GetString("tournament event id")
		tournament_event_stage_, _ := r.GetString("tournament event stage id")
		map_name, _ := r.GetString("map")

		// Get map name, this is weird but okay
		// Some tournament events might not have an ID, so we need to handle that
		tournament_event_id_int, _ := strconv.Atoi(tournament_event)
		tournament_event_stage_id_int, _ := strconv.Atoi(tournament_event_stage_)

		// Handle teams
		tournament_event_team0, _ := r.GetString("tournament event team0 id")
		tournament_event_team1, _ := r.GetString("tournament event team1 id")

		// convert to int
		team0, _ := t.GetValueByKey("CSGO_TeamID_" + tournament_event_team0)
		team1, _ := t.GetValueByKey("CSGO_TeamID_" + tournament_event_team1)

		var teams = models.HighlightReelTeams{
			TeamZero: team0, // This is the team with the player who made the highlight
			TeamOne:  team1, // This is the opposing team
		}

		reel_description, _ := t.GetValueByKey("HighlightDesc_" + id)

		tournament := modules.GetTournamentData(t, tournament_event_id_int)
		stage := modules.GetTournamentStageData(t, tournament_event_stage_id_int)
		current := models.HighlightReel{
			Id:             id,
			Tournament:     tournament,
			Stage:          stage,
			Map:            map_name,
			MarketHashName: modules.GenerateHighlightReelMarketHashName(t, id, tournament_event_id_int),
			ReelDescription: reel_description,
			Teams:          teams,
		}

		items = append(items, current)
	}

	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' highlight reels in %s", len(items), duration)

	return items
}
