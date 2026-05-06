package tests

// Package tests contains integration tests that validate the exported JSON
// output produced by the parser against known ground-truth values.
//
// Run with:
//
//	go test ./tests/...
//
// The tests read files from the ../exported/ directory relative to this file,
// so they must be run after the parser has generated its output (make dev or
// go run . from the root).

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// exportedDir returns the absolute path to the exported/ directory regardless
// of where `go test` is invoked from.
func exportedDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(file), "..", "exported")
}

// readJSON unmarshals a JSON file from the exported directory into v.
func readJSON(t *testing.T, name string, v any) {
	t.Helper()
	path := filepath.Join(exportedDir(t), name)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readJSON %s: %v — have you run the parser yet?", name, err)
	}
	if err := json.Unmarshal(b, v); err != nil {
		t.Fatalf("readJSON %s: unmarshal: %v", name, err)
	}
}

// ---------------------------------------------------------------------------
// shared types
// ---------------------------------------------------------------------------

// WeaponsFile represents the top-level weapon_skins.json map:
//
//	weapon_def_index → { name, image, sticker_count, type, paints: { paint_def_index → PaintEntry } }
type WeaponsFile map[string]WeaponEntry

type WeaponEntry struct {
	Name         string                `json:"name"`
	Image        string                `json:"image"`
	StickerCount int                   `json:"sticker_count"`
	Type         string                `json:"type"`
	Paints       map[string]PaintEntry `json:"paints"`
}

type PaintEntry struct {
	DefinitionIndex int    `json:"definition_index"`
	Name            string `json:"name"`
	ItemSetID       string `json:"item_set_id"`
	Image           string `json:"image"`
	Rarity          string `json:"rarity"`
	Float           struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"float"`
	Souvenir bool `json:"souvenir"`
	StatTrak  bool `json:"stattrak"`
}

// findSkin searches all weapons in the file for a paint whose name matches
// skinName under weapon weaponName (case-sensitive). Returns the entry and a
// descriptive key, or an error when not found.
func findSkin(wf WeaponsFile, weaponName, skinName string) (PaintEntry, string, error) {
	for weaponDefIdx, weapon := range wf {
		if weapon.Name != weaponName {
			continue
		}
		for paintDefIdx, paint := range weapon.Paints {
			if paint.Name == skinName {
				return paint, fmt.Sprintf("%s (weapon %s, paint %s)", skinName, weaponDefIdx, paintDefIdx), nil
			}
		}
	}
	return PaintEntry{}, "", fmt.Errorf("skin %q not found on weapon %q", skinName, weaponName)
}

// ---------------------------------------------------------------------------
// weapon_skins tests
// ---------------------------------------------------------------------------

func TestWeaponSkins(t *testing.T) {
	var wf WeaponsFile
	readJSON(t, "weapon_skins.json", &wf)

	tests := []struct {
		weapon   string
		skin     string
		rarity   string
		souvenir bool
		stattrak bool
	}{
		// AK-47 | Olive Polycam – part of set_realism_camo_uncommon loot list.
		// Historically mis-classified as "common" due to HasSuffix("uncommon", "common").
		{weapon: "AK-47", skin: "Olive Polycam", rarity: "uncommon", souvenir: false, stattrak: false},

		// AK-47 | B the Monster – Overpass 2024 collection, covert, souvenir only.
		{weapon: "AK-47", skin: "B the Monster", rarity: "ancient", souvenir: true, stattrak: false},
	}

	for _, tc := range tests {
		t.Run(tc.weapon+"|"+tc.skin, func(t *testing.T) {
			entry, key, err := findSkin(wf, tc.weapon, tc.skin)
			if err != nil {
				t.Fatal(err)
			}

			if entry.Rarity != tc.rarity {
				t.Errorf("%s: rarity = %q, want %q", key, entry.Rarity, tc.rarity)
			}
			if entry.Souvenir != tc.souvenir {
				t.Errorf("%s: souvenir = %v, want %v", key, entry.Souvenir, tc.souvenir)
			}
			if entry.StatTrak != tc.stattrak {
				t.Errorf("%s: stattrak = %v, want %v", key, entry.StatTrak, tc.stattrak)
			}
		})
	}
}
