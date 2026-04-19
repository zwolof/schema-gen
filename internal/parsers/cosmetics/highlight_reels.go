package cosmetics

import (
	"context"
	"strconv"

	"go-csitems-parser/internal/i18n"
	"go-csitems-parser/internal/marketname"
	"go-csitems-parser/internal/models"
	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/rs/zerolog"
)

// HighlightReels extracts the per-match highlight reels that Valve ships
// inside items_game.txt.
type HighlightReels struct{ base.Parser }

func NewHighlightReels() *HighlightReels {
	return &HighlightReels{Parser: base.New("highlight_reels")}
}

func (p *HighlightReels) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	reels, err := in.IG.Get("highlight_reels")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get highlight_reels from items_game.txt")
		return nil, nil
	}

	var out []models.HighlightReel
	defer p.LogCount(ctx, "highlight reels", func() int { return len(out) })()

	for _, r := range reels.GetChilds() {
		definition_index, _ := strconv.Atoi(r.Key)

		id, _ := r.GetString("id")
		eventStr, _ := r.GetString("tournament event id")
		stageStr, _ := r.GetString("tournament event stage id")
		mapName, _ := r.GetString("map")

		eventID, _ := strconv.Atoi(eventStr)
		stageID, _ := strconv.Atoi(stageStr)

		team0Str, _ := r.GetString("tournament event team0 id")
		team1Str, _ := r.GetString("tournament event team1 id")
		team0ID, _ := strconv.Atoi(team0Str)
		team1ID, _ := strconv.Atoi(team1Str)

		team0, _ := in.T.GetValueByKey("CSGO_TeamID_" + team0Str)
		team1, _ := in.T.GetValueByKey("CSGO_TeamID_" + team1Str)

		title, _ := in.T.GetValueByKey("HighlightReel_" + id)
		description, _ := in.T.GetValueByKey("HighlightDesc_" + id)

		out = append(out, models.HighlightReel{
			DefinitionIndex: definition_index,
			Id:              id,
			Tournament:      i18n.GetTournamentData(in.T, eventID),
			Stage:           i18n.GetTournamentStageData(in.T, stageID),
			Map:             mapName,
			MarketHashName:  marketname.GenerateHighlightReelMarketHashName(in.T, id, eventID),
			ReelDescription: description,
			ReelTitle:       title,
			Teams:           models.HighlightReelTeams{TeamZero: team0, TeamOne: team1},
			VideoUrl:        marketname.GenerateHighlightReelVideoURL(id, mapName, eventID, stageID, team0ID, team1ID, "ww"),
		})
	}

	return out, nil
}
