package generate

import (
	"testing"
)

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HelloWorld"},
		{"hello", "Hello"},
		{"hello_world_check", "HelloWorldCheck"},
		{"my_app", "MyApp"},
		{"get_state", "GetState"},
		{"app_id", "AppID"},
		{"get_url", "GetURL"},
		{"xml_parser", "XMLParser"},
		{"doNothing", "DoNothing"},
		{"appEquals", "AppEquals"},
		{"camelCase", "CamelCase"},
		{"", ""},
		{"subscribe_xgov", "SubscribeXgov"},
		{"config_xgov_registry", "ConfigXgovRegistry"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToPascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToPascalCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "helloWorld"},
		{"Hello", "hello"},
		{"ID", "id"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToCamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToPackageName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "helloworld"},
		{"My-Contract", "mycontract"},
		{"state_decoding", "statedecoding"},
		{"XGovRegistry", "xgovregistry"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToPackageName(tt.input)
			if result != tt.expected {
				t.Errorf("ToPackageName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSafeGoName(t *testing.T) {
	if SafeGoName("type") != "typeVal" {
		t.Errorf("expected typeVal, got %s", SafeGoName("type"))
	}
	if SafeGoName("name") != "name" {
		t.Errorf("expected name, got %s", SafeGoName("name"))
	}
}
