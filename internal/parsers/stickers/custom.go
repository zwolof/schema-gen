package stickers

import (
	"context"
	"strconv"
	"strings"

	"go-csitems-parser/internal/i18n"
	"go-csitems-parser/internal/itemsgame"
	"go-csitems-parser/internal/marketname"
	"go-csitems-parser/internal/models"
	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/rs/zerolog"
)

// groupSubID keys into the 1/2/3 "E/P/T" naming used by GenerateCustomStickerId.
var groupSubID = map[int]string{1: "E", 2: "P", 3: "T"}

var stickerEffectSuffix = map[string]string{
	"glitter":    "G",
	"holo":       "H",
	"foil":       "F",
	"gold":       "G",
	"lenticular": "L",
	"normal":     "N",
}

var stickerTypeSuffix = map[string]string{
	"normal": "N",
	"player": "P",
	"team":   "T",
	"event":  "E",
}

// kitIndex is a one-time lookup of sticker-kit counts keyed by
// tournament/team/player id + effect, plus total-count maps.
type kitIndex struct {
	countByPlayer     map[int]map[string]int
	countByTeam       map[int]map[string]int
	countByTournament map[int]map[string]int
	totalByPlayer     map[int]int
	totalByTeam       map[int]int
	totalByTournament map[int]int
}

func buildKitIndex(kits []models.StickerKit) *kitIndex {
	idx := &kitIndex{
		countByPlayer:     make(map[int]map[string]int),
		countByTeam:       make(map[int]map[string]int),
		countByTournament: make(map[int]map[string]int),
		totalByPlayer:     make(map[int]int),
		totalByTeam:       make(map[int]int),
		totalByTournament: make(map[int]int),
	}
	for i := range kits {
		k := &kits[i]
		if k.Player != nil && k.Player.Id > 0 {
			if idx.countByPlayer[k.Player.Id] == nil {
				idx.countByPlayer[k.Player.Id] = make(map[string]int)
			}
			idx.countByPlayer[k.Player.Id][k.Effect]++
			idx.totalByPlayer[k.Player.Id]++
		}
		if k.Team != nil && k.Team.Id > 0 {
			if idx.countByTeam[k.Team.Id] == nil {
				idx.countByTeam[k.Team.Id] = make(map[string]int)
			}
			idx.countByTeam[k.Team.Id][k.Effect]++
			// totalByTeam mirrors the legacy "team-only" count — kits with a
			// player set are excluded.
			if k.Player == nil {
				idx.totalByTeam[k.Team.Id]++
			}
		}
		if k.Tournament != nil && k.Tournament.Id > 0 {
			if idx.countByTournament[k.Tournament.Id] == nil {
				idx.countByTournament[k.Tournament.Id] = make(map[string]int)
			}
			idx.countByTournament[k.Tournament.Id][k.Effect]++
			idx.totalByTournament[k.Tournament.Id]++
		}
	}
	return idx
}

// Custom derives "custom sticker" aggregates (per-event, per-team, per-player
// with effect variants + totals) from the parsed sticker kits.
type Custom struct{ base.Parser }

func NewCustom() *Custom { return &Custom{Parser: base.New("custom_stickers")} }

func (c *Custom) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	if len(in.StickerKits) == 0 {
		logger.Warn().Msg("No sticker kits found, skipping custom stickers parsing")
		return []models.CustomStickers(nil), nil
	}

	var out []models.CustomStickers
	defer c.LogCount(ctx, "custom stickers", func() int { return len(out) })()

	idx := buildKitIndex(in.StickerKits)
	seen := make(map[string]struct{})
	seenTournament := make(map[int]struct{})
	seenPlayer := make(map[int]struct{})
	seenTeam := make(map[int]struct{})
	var uniqueTournamentIds, uniquePlayerIds, uniqueTeamIds []int

	for _, kit := range in.StickerKits {
		// Tournament stickers only.
		if kit.Tournament != nil && kit.Player == nil {
			id := kit.Tournament.Id
			if id <= 0 {
				logger.Debug().Msgf("Skipping sticker kit with invalid tournament ID: %s", kit.Name)
				continue
			}
			effect := StickerEffect(kit.StickerMaterial)
			typ := StickerType(0, id, 0)
			count := idx.countByTournament[id][effect]
			if count == 0 {
				continue
			}
			genID := GenerateCustomStickerId(id, groupSubID[1], &effect, &typ)
			if _, exists := seen[genID]; exists {
				continue
			}
			seen[genID] = struct{}{}
			out = append(out, models.CustomStickers{
				GeneratedId: genID,
				Group:       2,
				Count:       count,
				Name:        marketname.GenerateCustomStickerMarketHashName_Event(in.T, id, &effect),
			})
			if _, ok := seenTournament[id]; !ok {
				seenTournament[id] = struct{}{}
				uniqueTournamentIds = append(uniqueTournamentIds, id)
			}
		}

		// Team stickers only.
		if kit.Team != nil && kit.Player == nil {
			id := kit.Team.Id
			if id <= 0 {
				logger.Debug().Msgf("Skipping sticker kit with invalid team ID: %s", kit.Name)
				continue
			}
			effect := StickerEffect(kit.StickerMaterial)
			typ := StickerType(0, 0, id)
			count := idx.countByTeam[id][effect]
			if count == 0 {
				continue
			}
			genID := GenerateCustomStickerId(id, groupSubID[3], &effect, &typ)
			if _, exists := seen[genID]; exists {
				continue
			}
			seen[genID] = struct{}{}
			out = append(out, models.CustomStickers{
				GeneratedId: genID,
				Group:       3,
				Count:       count,
				Name:        marketname.GenerateCustomStickerMarketHashName_Team(in.T, id, &effect),
			})
			if _, ok := seenTeam[id]; !ok {
				seenTeam[id] = struct{}{}
				uniqueTeamIds = append(uniqueTeamIds, id)
			}
		}

		// Player stickers only.
		if kit.Player != nil {
			id := kit.Player.Id
			if id <= 0 {
				logger.Debug().Msgf("Skipping sticker kit with invalid player ID: %s", kit.Name)
				continue
			}
			effect := StickerEffect(kit.StickerMaterial)
			typ := StickerType(id, 0, 0)
			count := idx.countByPlayer[id][effect]
			if count == 0 {
				continue
			}
			genID := GenerateCustomStickerId(id, groupSubID[2], &effect, &typ)
			if _, exists := seen[genID]; exists {
				continue
			}
			seen[genID] = struct{}{}
			out = append(out, models.CustomStickers{
				GeneratedId: genID,
				Group:       2,
				Count:       count,
				Name:        marketname.GenerateCustomStickerMarketHashName_Player(in.T, kit.Player, &effect),
			})
			if _, ok := seenPlayer[id]; !ok {
				seenPlayer[id] = struct{}{}
				uniquePlayerIds = append(uniquePlayerIds, id)
			}
		}
	}

	for _, id := range uniqueTournamentIds {
		out = append(out, buildTotal(in.IG, in.T, idx, 1, id))
	}
	for _, id := range uniquePlayerIds {
		out = append(out, buildTotal(in.IG, in.T, idx, 2, id))
	}
	for _, id := range uniqueTeamIds {
		out = append(out, buildTotal(in.IG, in.T, idx, 3, id))
	}

	return out, nil
}

func buildTotal(ig *models.ItemsGame, t i18n.Translator, idx *kitIndex, groupID, subID int) models.CustomStickers {
	generatedID := GenerateCustomStickerId(subID, groupSubID[groupID], nil, nil)

	var count int
	var name string
	switch groupID {
	case 1:
		count = idx.totalByTournament[subID]
		event := i18n.GetTournamentData(t, subID)
		name = marketname.GenerateCustomStickerMarketHashName_Event(t, event.Id, nil)
	case 2:
		count = idx.totalByPlayer[subID]
		player := itemsgame.GetPlayerByAccountId(ig, subID)
		name = marketname.GenerateCustomStickerMarketHashName_Player(t, player, nil)
	case 3:
		count = idx.totalByTeam[subID]
		team := i18n.GetTournamentTeamData(t, subID)
		name = marketname.GenerateCustomStickerMarketHashName_Team(t, team.Id, nil)
	}

	return models.CustomStickers{
		GeneratedId: generatedID,
		Group:       groupID,
		Count:       count,
		Name:        name,
	}
}

// GenerateCustomStickerId is exported so external callers can reconstruct the
// canonical id format. Nil effect + stickerType produce the "A" aggregate form.
// Called once per custom sticker emitted — strings.Builder avoids the fmt.Sprintf
// reflection hit on the hot path.
func GenerateCustomStickerId(teamID int, kind string, effect, stickerType *string) string {
	var b strings.Builder
	// Worst case: 'C' + 10-digit int + 2 suffix chars + kind.
	b.Grow(13 + len(kind))
	b.WriteByte('C')
	b.WriteString(strconv.Itoa(teamID))

	if effect == nil || stickerType == nil {
		b.WriteByte('A')
	} else {
		b.WriteString(stickerEffectSuffix[*effect])
		b.WriteString(stickerTypeSuffix[*stickerType])
	}
	b.WriteString(kind)
	return b.String()
}
