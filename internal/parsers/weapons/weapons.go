// Package weapons holds the BaseWeapon-producing parsers: Weapons (rifle/
// pistol prefabs), Knives (melee_unusual), and Gloves (hands_paintable).
// They share the same output type and publish into Inputs.{Weapons,
// Knives, Gloves}.
package weapons

import (
	"context"
	"strconv"
	"strings"

	"go-csitems-parser/internal/models"
	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"

	"github.com/rs/zerolog"
)

// Weapons extracts core weapon prefabs (rifles, pistols, smgs, snipers).
type Weapons struct{ base.Parser }

func NewWeapons() *Weapons { return &Weapons{Parser: base.New("weapons")} }

func (w *Weapons) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	logger := zerolog.Ctx(ctx)

	prefabs, err := in.IG.Get("prefabs")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get prefabs")
		return nil, nil
	}

	gameInfo, err := in.IG.Get("game_info")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get game_info")
		return nil, nil
	}

	var out []models.BaseWeapon
	defer w.LogCount(ctx, "weapons", func() int { return len(out) })()

	maxStickers, _ := gameInfo.GetInt("max_num_stickers")

	for _, p := range prefabs.GetChilds() {
		if !strings.HasPrefix(p.Key, "weapon_") || !strings.HasSuffix(p.Key, "_prefab") {
			continue
		}

		itemClass := strings.TrimSuffix(p.Key, "_prefab")
		defIdx := lookupDefinitionIndex(itemClass, in.IG)
		if defIdx == -1 {
			logger.Error().Msgf("Failed to get definition index for weapon class '%s'", itemClass)
			continue
		}

		if _, err := p.Get("paint_data"); err != nil {
			continue
		}

		itemName, _ := p.GetString("item_name")
		itemDesc, _ := p.GetString("item_description")
		imageInventory, _ := p.GetString("image_inventory")

		translatedName, err := in.T.GetValueByKey(itemName)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to translate item name for weapon %s", itemName)
			translatedName = itemName
		}

		translatedDesc, err := in.T.GetValueByKey(itemDesc)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to translate item description for weapon %s", itemDesc)
			translatedDesc = itemDesc
		}

		out = append(out, models.BaseWeapon{
			DefinitionIndex: defIdx,
			Name:            translatedName,
			Description:     translatedDesc,
			ClassName:       itemClass,
			ImageInventory:  imageInventory,
			NumStickers:     maxStickers,
		})
	}

	return out, nil
}

func (w *Weapons) Commit(in *pipeline.Inputs, result any) {
	if r, ok := result.([]models.BaseWeapon); ok {
		in.Weapons = r
	}
}

// lookupDefinitionIndex finds the items[].Key for the given classname.
// Linear scan: called ~20 times per run (one per weapon prefab).
func lookupDefinitionIndex(class string, ig *models.ItemsGame) int {
	items, err := ig.Get("items")
	if err != nil {
		return -1
	}

	for _, w := range items.GetChilds() {
		name, _ := w.GetString("name")
		if name == class {
			idx, _ := strconv.Atoi(w.Key)
			return idx
		}
	}
	return -1
}

// prefab constants kept for clarity of intent.
const (
	prefabMeleeUnusual   = "melee_unusual"
	prefabHandsPaintable = "hands_paintable"
)
