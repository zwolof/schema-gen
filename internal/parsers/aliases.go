package parsers

import "go-csitems-parser/internal/parsers/pipeline"

// These aliases let callers (main.go, tests, etc.) talk to the pipeline
// infrastructure through the parsers package without a second import.
type (
	Inputs   = pipeline.Inputs
	Parser   = pipeline.Parser
	Tier     = pipeline.Tier
	Pipeline = pipeline.Pipeline
	Export   = pipeline.Export
)
