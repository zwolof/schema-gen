package i18n

// Factory holds translators keyed by language name (e.g. "English", "French").
type Factory struct {
	Translators map[string]*FileTranslator
}

// Get returns the translator for the given language name, or nil if absent.
// Returned as [Translator] so callers can swap in alternative implementations
// (tests, stubs) without casting.
func (f *Factory) Get(language string) Translator {
	if t, ok := f.Translators[language]; ok {
		return t
	}

	return nil
}
