// Package parsers is the top-level parser registry. It exposes [Default] — a
// ready-to-run [Pipeline] that fans every CS2 item-type parser across three
// tiers — and the type aliases that let callers use "parsers.Inputs",
// "parsers.Parser" etc. without a second import.
//
// # Layout
//
// The actual parsers live in domain sub-packages:
//
//   - parsers/pipeline    — Parser interface, Inputs, Tier, Pipeline infra
//   - parsers/base        — embeddable base.Parser for shared methods
//   - parsers/cosmetics   — agents, music_kits, collectibles, keychains, highlight_reels
//   - parsers/weapons     — weapons, knives, gloves (BaseWeapon family)
//   - parsers/containers  — weapon_cases, souvenir_packages, sticker_capsules, misc_capsules
//   - parsers/stickers    — sticker_kits, custom_stickers, sticker_slabs
//   - parsers/meta        — rarities, paint_kits, item_sets, collections, armory_rewards +
//     the loot-list index builders (SkinWeaponRarityMap, StickerItemSetMap)
//   - parsers/skins       — Parser wrappers around schema.Get{Weapon,Knife,Glove}PaintKits
//
// # Adding a parser
//
// 1. Pick (or create) the domain sub-package.
// 2. Write a struct that embeds base.Parser:
//
//	type MyThing struct{ base.Parser }
//
//	func NewMyThing() *MyThing { return &MyThing{Parser: base.New("my_thing")} }
//
//	func (m *MyThing) Parse(ctx context.Context, in *pipeline.Inputs) (any, error) {
//	    var out []models.MyThing
//	    defer m.LogCount(ctx, "my things", func() int { return len(out) })()
//	    // ...
//	    return out, nil
//	}
//
// 3. If the output feeds a later tier, override Commit on the struct:
//
//	func (m *MyThing) Commit(in *pipeline.Inputs, result any) {
//	    if r, ok := result.([]models.MyThing); ok {
//	        in.MyThings = r
//	    }
//	}
//
// 4. Register the struct in [Default]'s appropriate tier in registry.go.
package parsers
