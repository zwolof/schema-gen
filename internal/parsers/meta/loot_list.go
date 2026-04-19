// Package meta holds the "metadata" parsers and indexes: rarity + paint-kit +
// item-set + collection + armory-reward parsers plus the loot-list helpers
// used to build the cross-cutting SkinRarity and StickerItemSet index maps.
package meta

import (
	"regexp"
	"strings"
)

// lootListItemRegexp matches "[paintkit]weaponclass" entries in
// client_loot_lists sub-lists. Shared by SkinRarityMap and StickerItemSetMap.
var lootListItemRegexp = regexp.MustCompile(`^\[(.+?)\](.+)$`)

// lootListRarityEndings enumerates the known rarity suffixes that appear at
// the end of client_loot_lists sub-list keys. Order matters: longer/more
// specific suffixes are checked first so that, e.g. "mythical" isn't matched
// as "common".
var lootListRarityEndings = []string{
	"default",
	"common",
	"uncommon",
	"rare",
	"mythical",
	"legendary",
	"ancient",
	"immortal",
	"unusual",
}

// lootListRarity returns the rarity suffix of a sub-list key, or "default"
// if none is recognised. Aggregator keys (e.g. "crate_abc") typically return
// "default" and are skipped by callers.
func lootListRarity(name string) string {
	for _, ending := range lootListRarityEndings {
		if strings.HasSuffix(name, ending) {
			return ending
		}
	}
	return "default"
}
