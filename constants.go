package crudex

const (
	CmdArgStrategyAlways      = "always"
	CmdArgStrategyIfNotExists = "newonly"
	CmdArgStrategyNever       = "never"
)

//go:generate stringer -type=ScaffoldStrategy
type ScaffoldStrategy int

const (
	// SCAFFOLD_ALWAYS will always scaffold the model templates
	SCAFFOLD_ALWAYS ScaffoldStrategy = iota
	// SCAFFOLD_IF_NOT_EXISTS will only scaffold the model templates if they do not exist
	SCAFFOLD_IF_NOT_EXISTS
	// SCAFFOLD_NEVER will never scaffold the model templates
	SCAFFOLD_NEVER
)
