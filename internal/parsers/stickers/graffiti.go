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

// Graffiti extracts spray_* entries from sticker_kits and exports them as
// graffiti.json. These are graffiti/spray items, distinct from stickers.
type Graffiti struct{ base.Parser }

func NewGraffiti() *Graffiti { return &Graffiti{Parser: base.New("graffiti")} }

func (g *Graffiti) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	kits, err := in.IG.Get("sticker_kits")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get sticker_kits for graffiti parser")
		return nil, nil
	}

	tints := loadGraffitiTints(in)

	var out []models.Graffiti
	defer g.LogCount(ctx, "graffiti", func() int { return len(out) })()

	for _, item := range kits.GetChilds() {
		definition_index, _ := strconv.Atoi(item.Key)
		if definition_index <= 0 {
			continue
		}

		name, _ := item.GetString("name")
		if !strings.HasPrefix(name, "spray_") {
			continue
		}

		item_name, _ := item.GetString("item_name")
		sticker_material, _ := item.GetString("sticker_material")
		item_rarity, _ := item.GetString("item_rarity")
		tournament_event_id, _ := item.GetInt("tournament_event_id")
		tournament_team_id, _ := item.GetInt("tournament_team_id")
		tournament_player_id, _ := item.GetInt("tournament_player_id")

		base_name := marketname.GenerateMarketHashName(in.T, item_name, nil, "graffiti")

		// Consumer-grade graffiti come in every available tint colour.
		// Each tint produces its own market listing so we embed the full
		// tint list with per-tint market hash names.
		var resolvedTints []models.GraffitiTint
		if item_rarity == "common" {
			resolvedTints = make([]models.GraffitiTint, len(tints))
			for i, t := range tints {
				resolvedTints[i] = models.GraffitiTint{
					ID:             t.id,
					Name:           t.name,
					Hex:            t.hex,
					MarketHashName: base_name + " (" + t.name + ")",
				}
			}
		}

		out = append(out, models.Graffiti{
			DefinitionIndex: definition_index,
			Name:            name,
			MarketHashName:  base_name,
			StickerMaterial: sticker_material,
			Image:           "econ/stickers/" + sticker_material,
			Rarity:          item_rarity,
			Tints:           resolvedTints,
			Tournament:      i18n.GetTournamentData(in.T, tournament_event_id),
			Team:            i18n.GetTournamentTeamData(in.T, tournament_team_id),
			Player:          itemsgame.GetPlayerByAccountId(in.IG, tournament_player_id),
		})
	}

	return out, nil
}

// graffitiTint is an internal representation used during parse only.
type graffitiTint struct {
	id   int
	name string
	hex  string
}

// loadGraffitiTints reads the graffiti_tints block from items_game and
// resolves tint names via the translator.
func loadGraffitiTints(in *pipeline.Inputs) []graffitiTint {
	tintNode, err := in.IG.Get("graffiti_tints")
	if err != nil {
		return nil
	}

	var tints []graffitiTint
	for _, node := range tintNode.GetChilds() {
		id, _ := node.GetInt("id")
		if id <= 0 {
			continue
		}
		hex, _ := node.GetString("hex_color")

		// Tint display name is localised as "Attrib_SprayTintValue_{id}"
		locKey := "#Attrib_SprayTintValue_" + strconv.Itoa(id)
		tintName, err := in.T.GetValueByKey(locKey)
		if err != nil || tintName == "" {
			tintName = node.Key // fall back to internal key (e.g. "brick_red")
		}

		tints = append(tints, graffitiTint{id: id, name: tintName, hex: hex})
	}
	return tints
}

func (g *Graffiti) Commit(in *pipeline.Inputs, result any) {
	if r, ok := result.([]models.Graffiti); ok {
		in.GraffitiKits = r
	}
}
