package parsers

import (
	"go-csitems-parser/internal/parsers/containers"
	"go-csitems-parser/internal/parsers/cosmetics"
	"go-csitems-parser/internal/parsers/meta"
	"go-csitems-parser/internal/parsers/skins"
	"go-csitems-parser/internal/parsers/stickers"
	"go-csitems-parser/internal/parsers/weapons"
	"go-csitems-parser/internal/skinmap"
)

// Default is the canonical schema-gen pipeline. Tiers run sequentially;
// parsers within a tier run concurrently.
var Default = Pipeline{
	// Tier 0 — read only IG + T. Outputs feed later tiers via Commit.
	{Parsers: []Parser{
		cosmetics.NewAgents(),
		cosmetics.NewMusicKits(),
		cosmetics.NewCollectibles(),
		cosmetics.NewKeychains(),
		cosmetics.NewHighlightReels(),

		weapons.NewWeapons(),
		weapons.NewGloves(),
		weapons.NewKnives(),

		containers.NewWeaponCases(),
		containers.NewSouvenirPackages(),
		containers.NewStickerCapsules(),
		containers.NewMiscCapsules(),

		meta.NewRarities(),
		meta.NewPaintKits(),
	}},

	// Tier 1 — reads Tier-0 outputs.
	{
		Parsers: []Parser{
			stickers.NewKits(),
			meta.NewItemSets(),
			meta.NewCollections(),
		},
		// Apply item-set metadata to paint kits before Tier 2 reads them
		// concurrently — mutation-through-reads is unsafe.
		AfterCommit: func(in *Inputs) {
			skinmap.Default.EnrichPaintKits(in.PaintKits, in.ItemSets, in.Weapons)
		},
	},

	// Tier 2 — reads Tier-1 outputs.
	{Parsers: []Parser{
		stickers.NewCustom(),
		stickers.NewSlabs(),
		meta.NewArmoryRewards(),
		skins.NewKnife(),
		skins.NewWeapon(),
		skins.NewGlove(),
	}},
}
