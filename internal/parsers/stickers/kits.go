// Package stickers holds the sticker-related parsers. Kits builds the
// canonical StickerKit slice used by custom_stickers and sticker_slabs.
package stickers

import (
	"context"
	"fmt"
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

// Kits extracts sticker_kits from items_game.txt and publishes them as
// Inputs.StickerKits so CustomStickers and Slabs can read them.
type Kits struct{ base.Parser }

func NewKits() *Kits { return &Kits{Parser: base.New("sticker_kits")} }

func (k *Kits) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	kits, err := in.IG.Get("sticker_kits")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get sticker kits from items_game.txt")
		return nil, nil
	}

	var out []models.StickerKit
	defer k.LogCount(ctx, "sticker kits", func() int { return len(out) })()

	for _, item := range kits.GetChilds() {
		definition_index, _ := strconv.Atoi(item.Key)
		if definition_index <= 0 {
			logger.Debug().Msgf("Skipping invalid sticker kit definition index: %d", definition_index)
			continue
		}

		item_name, _ := item.GetString("item_name")
		name, _ := item.GetString("name")
		if strings.Contains(name, "patch_") || strings.Contains(name, "_graffiti") {
			continue
		}

		sticker_material, _ := item.GetString("sticker_material")
		item_rarity, _ := item.GetString("item_rarity")
		tournament_event_id, _ := item.GetInt("tournament_event_id")
		tournament_team_id, _ := item.GetInt("tournament_team_id")
		tournament_player_id, _ := item.GetInt("tournament_player_id")

		out = append(out, models.StickerKit{
			DefinitionIndex: definition_index,
			Name:            name,
			MarketHashName:  marketname.GenerateMarketHashName(in.T, item_name, nil, "sticker_kit"),
			StickerMaterial: sticker_material,
			Image:           fmt.Sprintf("econ/stickers/%s", sticker_material),
			Rarity:          item_rarity,
			Effect:          StickerEffect(sticker_material),
			Type:            StickerType(tournament_player_id, tournament_team_id, tournament_event_id),
			ItemSetId:       in.StickerItemSetMap[name],
			Tournament:      i18n.GetTournamentData(in.T, tournament_event_id),
			Team:            i18n.GetTournamentTeamData(in.T, tournament_team_id),
			Player:          itemsgame.GetPlayerByAccountId(in.IG, tournament_player_id),
		})
	}

	return out, nil
}

func (k *Kits) Commit(in *pipeline.Inputs, result any) {
	if r, ok := result.([]models.StickerKit); ok {
		in.StickerKits = r
	}
}

// StickerType classifies a kit by which tournament-ID field is populated.
//
// The upstream schema was originally generated from a function whose
// parameter labels were swapped (GetStickerType(player, event, team)) but
// whose sticker_kits call site passed (player_id, team_id, event_id). The
// net effect cemented into the downstream JSON:
//
//   - tournament_event_id set  → "team"
//   - tournament_team_id set   → "event"
//   - tournament_player_id set → "autograph"
//
// We replicate that mapping here so exported hashes stay byte-stable.
func StickerType(player, team, event int) string {
	switch {
	case player > 0:
		return "autograph"
	case event > 0:
		return "team"
	case team > 0:
		return "event"
	default:
		return "normal"
	}
}

// effectSuffixes lists the sticker_material suffixes in deterministic match
// order. First match wins.
var effectSuffixes = []struct {
	suffix, effect string
}{
	{"_glitter", "glitter"},
	{"_holo", "holo"},
	{"_foil", "foil"},
	{"_gold", "gold"},
	{"_lenticular", "lenticular"},
	{"_embroidered", "embroidered"},
}

// StickerEffect derives the visual effect from the sticker_material suffix.
func StickerEffect(material string) string {
	for _, e := range effectSuffixes {
		if strings.HasSuffix(material, e.suffix) {
			return e.effect
		}
	}
	return "normal"
}
