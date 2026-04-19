// Package skins wraps the internal/skinmap Builder as Parser implementations
// so the pipeline can fan out weapon/knife/glove skin construction alongside
// the item parsers.
package skins

import (
	"context"

	"go-csitems-parser/internal/parsers/base"
	"go-csitems-parser/internal/parsers/pipeline"
	"go-csitems-parser/internal/skinmap"
)

// Weapon builds the weapon skin map via skinmap.Default.
type Weapon struct{ base.Parser }

func NewWeapon() *Weapon { return &Weapon{Parser: base.New("weapon_skins")} }

func (w *Weapon) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	return skinmap.Default.Weapons(in.Weapons, in.PaintKits, in.ItemSets, in.SkinRarityMap), nil
}

// Knife builds the knife skin map via skinmap.Default.
type Knife struct{ base.Parser }

func NewKnife() *Knife { return &Knife{Parser: base.New("knife_skins")} }

func (k *Knife) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	return skinmap.Default.Knives(in.Knives, in.PaintKits, in.KnifeSkinMap, in.SkinRarityMap), nil
}

// Glove builds the glove skin map via skinmap.Default.
type Glove struct{ base.Parser }

func NewGlove() *Glove { return &Glove{Parser: base.New("glove_skins")} }

func (g *Glove) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
	return skinmap.Default.Gloves(in.Gloves, in.PaintKits, in.KnifeSkinMap, in.SkinRarityMap), nil
}
