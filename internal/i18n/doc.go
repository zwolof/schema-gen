// Package i18n resolves Valve localisation keys (e.g. "#SFUI_Scoreboard_Team")
// to their localised strings.
//
// Callers depend on [Translator] (the interface). [FileTranslator] is the
// default file-backed implementation that [Load] and [LoadAll] return. The
// [Factory] groups translators by language name so callers can pick one
// without handling raw file paths.
package i18n
