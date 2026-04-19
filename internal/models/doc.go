// Package models defines the pure data types shared between the loader,
// the parsers, and the schema exporters. Types here carry no logic — they
// exist to be unmarshalled from Valve's items_game.txt VDF tree and
// re-marshalled into the JSON files written to ./exported/.
//
// Core types:
//
//   - ItemsGame wraps *vdf.KeyValue and is the root of every parser's input.
//   - BaseWeapon, PaintKit, StickerKit, ItemSet, Rarity carry the per-item
//     metadata extracted from items_game.txt.
//   - Schema* types are the shape of the JSON that csfloat.com consumes.
//
// See internal/parsers for the logic that produces these values and
// internal/schema/mapper.go for the types that combine them.
package models
