package generate

import (
	algokit "github.com/kylebeee/algokit-utils-go"
)

// MethodCallConfig describes how a method can be called.
type MethodCallConfig struct {
	CanCreate     bool
	CanCall       bool
	CanOptIn      bool
	CanCloseOut   bool
	CanUpdate     bool
	CanDelete     bool
	IsReadonly    bool
}

// AnalyzeCallConfig determines the call configuration for a method.
func AnalyzeCallConfig(method algokit.Arc56Method) MethodCallConfig {
	config := MethodCallConfig{
		IsReadonly: method.Readonly,
	}

	for _, action := range method.Actions.Create {
		switch action {
		case "NoOp":
			config.CanCreate = true
		case "OptIn":
			config.CanCreate = true
			config.CanOptIn = true
		}
	}

	for _, action := range method.Actions.Call {
		switch action {
		case "NoOp":
			config.CanCall = true
		case "OptIn":
			config.CanOptIn = true
		case "CloseOut":
			config.CanCloseOut = true
		case "UpdateApplication":
			config.CanUpdate = true
		case "DeleteApplication":
			config.CanDelete = true
		}
	}

	return config
}

// BareCallConfig describes what bare calls are allowed.
type BareCallConfig struct {
	CanCreate   bool
	CanCall     bool
	CanOptIn    bool
	CanCloseOut bool
	CanUpdate   bool
	CanDelete   bool
}

// AnalyzeBareConfig determines the bare call configuration.
func AnalyzeBareConfig(bareActions algokit.BareActions) BareCallConfig {
	config := BareCallConfig{}

	for _, action := range bareActions.Create {
		switch action {
		case "NoOp":
			config.CanCreate = true
		case "OptIn":
			config.CanCreate = true
			config.CanOptIn = true
		}
	}

	for _, action := range bareActions.Call {
		switch action {
		case "NoOp":
			config.CanCall = true
		case "OptIn":
			config.CanOptIn = true
		case "CloseOut":
			config.CanCloseOut = true
		case "UpdateApplication":
			config.CanUpdate = true
		case "DeleteApplication":
			config.CanDelete = true
		}
	}

	return config
}

// HasCreateMethod returns true if any method supports create.
func HasCreateMethod(methods []algokit.Arc56Method) bool {
	for _, m := range methods {
		config := AnalyzeCallConfig(m)
		if config.CanCreate {
			return true
		}
	}
	return false
}
