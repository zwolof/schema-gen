package i18n

import (
	"strconv"

	"go-csitems-parser/internal/models"
)

// GetTournamentData resolves a tournament event id to its short-name data.
// Returns nil when the id is 0 or the translation is missing.
func GetTournamentData(t Translator, id int) *models.TournamentData {
	return lookupTournamentData(t, "CSGO_Tournament_Event_NameShort_"+strconv.Itoa(id), id)
}

// GetTournamentStageData resolves a tournament event stage id.
func GetTournamentStageData(t Translator, id int) *models.TournamentData {
	return lookupTournamentData(t, "CSGO_Tournament_Event_Stage_"+strconv.Itoa(id), id)
}

// GetTournamentTeamData resolves a tournament team id to its localised name.
func GetTournamentTeamData(t Translator, id int) *models.TournamentData {
	return lookupTournamentData(t, "CSGO_TeamID_"+strconv.Itoa(id), id)
}

func lookupTournamentData(t Translator, key string, id int) *models.TournamentData {
	name, _ := t.GetValueByKey(key)
	if id == 0 || name == "" {
		return nil
	}
	return &models.TournamentData{Id: id, Name: name}
}
