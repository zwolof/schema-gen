package i18n

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/baldurstod/vdf"
	"github.com/rs/zerolog"
)

// Load reads a single csgo_<lang>.txt file from folderPath. Use this when
// only one language is needed — it avoids the ~150 MB of unused parsing that
// [LoadAll] does.
func Load(ctx context.Context, folderPath, lang string) (*Factory, error) {
	logger := zerolog.Ctx(ctx)
	start := time.Now()

	fileName := fmt.Sprintf("csgo_%s.txt", strings.ToLower(lang))
	path := fmt.Sprintf("%s/%s", folderPath, fileName)

	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	parser := vdf.VDF{}
	kv := parser.Parse(removeBOM(fileData))
	if kv.Value == nil {
		return nil, fmt.Errorf("parse %s: empty VDF", path)
	}

	t, langName := loadLanguage(&kv)
	if t == nil {
		return nil, fmt.Errorf("load language from %s", path)
	}

	logger.Info().Msgf("Loaded language '%s' in %s", langName, time.Since(start))

	return &Factory{
		Translators: map[string]*FileTranslator{langName: t},
	}, nil
}

// LoadAll reads every csgo_*.txt file in folderPath into a [Factory].
// Per-file parse failures are logged and skipped; the only returned error
// is from reading the directory itself.
func LoadAll(ctx context.Context, folderPath string) (*Factory, error) {
	logger := zerolog.Ctx(ctx)

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("read dir %s: %w", folderPath, err)
	}
	logger.Info().Msgf("Found '%d' language files", len(files))

	start := time.Now()
	langMap := make(map[string]*FileTranslator)
	parser := vdf.VDF{}

	for _, entry := range files {
		if entry.IsDir() {
			logger.Info().Msgf("Skipping directory %s", entry.Name())
			continue
		}

		fileName := entry.Name()
		if !strings.HasPrefix(fileName, "csgo_") || !strings.HasSuffix(fileName, ".txt") {
			logger.Info().Msgf("Skipping file %s", fileName)
			continue
		}

		path := fmt.Sprintf("%s/%s", folderPath, fileName)
		fileData, err := os.ReadFile(path)
		if err != nil {
			logger.Error().Msgf("Error reading file %s: %v", path, err)
			continue
		}

		kv := parser.Parse(removeBOM(fileData))
		if kv.Value == nil {
			logger.Error().Msgf("Error parsing file %s", path)
			continue
		}

		t, langName := loadLanguage(&kv)
		if t == nil {
			logger.Error().Msgf("Error loading language from file %s", path)
			continue
		}
		langMap[langName] = t
	}

	logger.Info().Msgf("Parsed '%d' language files in %s", len(files), time.Since(start))
	return &Factory{Translators: langMap}, nil
}

func loadLanguage(kv *vdf.KeyValue) (*FileTranslator, string) {
	if kv == nil {
		panic("vdf is nil")
	}

	inner, _ := kv.Get("lang")
	langName, _ := inner.GetString("Language")

	tokens, _ := inner.Get("Tokens")
	if inner == nil || tokens == nil {
		panic("translation file does not contain 'lang.Tokens' section")
	}

	tokenMap, err := tokens.ToStringMap()
	if err != nil {
		panic(fmt.Sprintf("Error parsing tokens: %v", err))
	}

	// Case-fold all keys for consistent lookup.
	for key, value := range *tokenMap {
		lower := strings.ToLower(key)
		if lower != key {
			(*tokenMap)[lower] = value
			delete(*tokenMap, key)
		}
	}

	return &FileTranslator{Language: langName, Tokens: *tokenMap}, langName
}

func removeBOM(b []byte) []byte {
	return bytes.Trim(b, "\xef\xbb\xbf")
}
