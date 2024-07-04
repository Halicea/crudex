package crudex

const (
	CmdArgStrategyAlways      = "always"
	CmdArgStrategyIfNotExists = "newonly"
	CmdArgStrategyNever       = "never"
)

//go:generate stringer -type=ScaffoldStrategy
type ScaffoldStrategy int

const (
	// ScaffoldStrategyAlways will always scaffold the model templates
	ScaffoldStrategyAlways ScaffoldStrategy = iota
	// ScaffoldStrategyIfNotExists will only scaffold the model templates if they do not exist
	ScaffoldStrategyIfNotExists
	// ScaffoldStrategyNever will never scaffold the model templates
	ScaffoldStrategyNever
)
