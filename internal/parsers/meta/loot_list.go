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

// lootListRarity returns the rarity of a sub-list key by checking whether its
// last underscore-delimited segment is a known rarity token derived from the
// parsed rarities block. Returns "default" when no match is found so callers
// can skip aggregator keys (e.g. "crate_abc").
func lootListRarity(name string, rarityKeys map[string]bool) string {
	if idx := strings.LastIndex(name, "_"); idx >= 0 {
		if suffix := name[idx+1:]; rarityKeys[suffix] {
			return suffix
		}
	}
	return "default"
}
