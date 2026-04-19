package i18n

import (
	"fmt"

	"go-csitems-parser/internal/models"
)

// GetTournamentData resolves a tournament event id to its short-name data.
// Returns nil when the id is 0 or the translation is missing.
func GetTournamentData(t Translator, id int) *models.TournamentData {
	return lookupTournamentData(t, fmt.Sprintf("CSGO_Tournament_Event_NameShort_%d", id), id)
}

// GetTournamentStageData resolves a tournament event stage id.
func GetTournamentStageData(t Translator, id int) *models.TournamentData {
	return lookupTournamentData(t, fmt.Sprintf("CSGO_Tournament_Event_Stage_%d", id), id)
}

// GetTournamentTeamData resolves a tournament team id to its localised name.
func GetTournamentTeamData(t Translator, id int) *models.TournamentData {
	return lookupTournamentData(t, fmt.Sprintf("CSGO_TeamID_%d", id), id)
}

func lookupTournamentData(t Translator, key string, id int) *models.TournamentData {
	name, _ := t.GetValueByKey(key)
	if id == 0 || name == "" {
		return nil
	}

	return &models.TournamentData{Id: id, Name: name}
}
