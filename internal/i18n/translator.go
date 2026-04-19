package i18n

import (
	"errors"
	"fmt"
	"strings"
)

// Translator resolves a localisation key to its translated string. The "#"
// prefix on Valve keys is stripped by implementations; lookups are
// case-insensitive.
type Translator interface {
	GetValueByKey(key string) (string, error)
}

// keyAliases maps items_game.txt keys to the actual (different) keys found
// in the language files. Covers cases where Valve shipped mismatched tokens.
var keyAliases = map[string]string{
	"stickerkit_dhw2014_dignitas_gold": "stickerkit_dhw2014_teamdignitas_gold",
}

// FileTranslator is a [Translator] backed by a single parsed language file.
// Safe for concurrent reads after construction.
type FileTranslator struct {
	Language string
	Tokens   map[string]string
}

// GetValueByKey implements [Translator].
func (t *FileTranslator) GetValueByKey(key string) (string, error) {
	key = strings.Replace(key, "#", "", 1)
	key = strings.ToLower(key)

	if t == nil {
		return key, errors.New("translator is nil")
	}

	if alias, ok := keyAliases[key]; ok {
		key = alias
	}

	value, ok := t.Tokens[key]
	if !ok || value == "" {
		return key, fmt.Errorf("key not found: %s", key)
	}
	return value, nil
}
