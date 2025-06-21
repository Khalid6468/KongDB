package integration_test

import (
	"testing"

	"github.com/khalid64/kongdb/internal/types"
)

func TestStringTypes_Integration(t *testing.T) {
	// Test type registry integration
	varchar := types.NewVarcharType(255, "utf8_general_ci")
	text := types.NewTextType("utf8_general_ci")

	err := types.RegisterGlobal(varchar)
	if err != nil {
		t.Fatalf("failed to register varchar: %v", err)
	}

	err = types.RegisterGlobal(text)
	if err != nil {
		t.Fatalf("failed to register text: %v", err)
	}

	// Test retrieval
	retrievedVarchar, exists := types.GetGlobalByID(types.TypeIDVarchar)
	if !exists {
		t.Fatal("varchar type not found in registry")
	}

	if retrievedVarchar.String() != varchar.String() {
		t.Errorf("retrieved varchar doesn't match original")
	}

	retrievedText, exists := types.GetGlobalByID(types.TypeIDText)
	if !exists {
		t.Fatal("text type not found in registry")
	}

	if retrievedText.String() != text.String() {
		t.Errorf("retrieved text doesn't match original")
	}

	// Test round-trip serialization
	testValue := "hello world"

	data, err := varchar.Serialize(testValue)
	if err != nil {
		t.Fatalf("serialization failed: %v", err)
	}

	result, err := varchar.Deserialize(data)
	if err != nil {
		t.Fatalf("deserialization failed: %v", err)
	}

	if result != testValue {
		t.Errorf("round-trip failed: expected %v, got %v", testValue, result)
	}
}

func TestStringTypes_CrossTypeCompatibility(t *testing.T) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")
	text := types.NewTextType("utf8_general_ci")

	// Test that both types can handle the same data
	testValues := []string{
		"",
		"hello",
		"café",
		"Hello World",
		"string with 'quotes'",
		"unicode: 你好世界",
	}

	for _, value := range testValues {
		t.Run(value, func(t *testing.T) {
			// Test VARCHAR
			if err := varchar.Validate(value); err != nil {
				t.Errorf("varchar validation failed: %v", err)
			}

			varcharData, err := varchar.Serialize(value)
			if err != nil {
				t.Errorf("varchar serialization failed: %v", err)
				return
			}

			varcharResult, err := varchar.Deserialize(varcharData)
			if err != nil {
				t.Errorf("varchar deserialization failed: %v", err)
				return
			}

			if varcharResult != value {
				t.Errorf("varchar round-trip failed: expected %v, got %v", value, varcharResult)
			}

			// Test TEXT
			if err := text.Validate(value); err != nil {
				t.Errorf("text validation failed: %v", err)
			}

			textData, err := text.Serialize(value)
			if err != nil {
				t.Errorf("text serialization failed: %v", err)
				return
			}

			textResult, err := text.Deserialize(textData)
			if err != nil {
				t.Errorf("text deserialization failed: %v", err)
				return
			}

			if textResult != value {
				t.Errorf("text round-trip failed: expected %v, got %v", value, textResult)
			}
		})
	}
}

func TestStringTypes_CollationCompatibility(t *testing.T) {
	// Test different collations
	collations := []string{
		"utf8_general_ci",
		"utf8_bin",
		"utf8_unicode_ci",
	}

	testPairs := []struct {
		a, b     string
		expected int
	}{
		{"hello", "HELLO", 0}, // Should be equal in case-insensitive
		{"hello", "world", -1},
		{"world", "hello", 1},
	}

	for _, collation := range collations {
		t.Run(collation, func(t *testing.T) {
			varchar := types.NewVarcharType(255, collation)
			text := types.NewTextType(collation)

			for _, pair := range testPairs {
				varcharResult := varchar.Compare(pair.a, pair.b)
				textResult := text.Compare(pair.a, pair.b)

				// Both types should behave the same with same collation
				if varcharResult != textResult {
					t.Errorf("collation mismatch: varchar=%d, text=%d for %s vs %s",
						varcharResult, textResult, pair.a, pair.b)
				}

				// For case-insensitive collations, "hello" and "HELLO" should be equal
				if collation != "utf8_bin" && pair.a == "hello" && pair.b == "HELLO" {
					if varcharResult != 0 {
						t.Errorf("case-insensitive comparison failed: %d", varcharResult)
					}
				}
			}
		})
	}
}

func TestStringTypes_JSONCompatibility(t *testing.T) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")
	text := types.NewTextType("utf8_general_ci")

	testValues := []interface{}{
		"hello",
		"café",
		"string with \"quotes\"",
		nil,
	}

	for _, value := range testValues {
		t.Run("varchar", func(t *testing.T) {
			jsonData, err := varchar.SerializeJSON(value)
			if err != nil {
				t.Errorf("varchar JSON serialization failed: %v", err)
				return
			}

			result, err := varchar.DeserializeJSON(jsonData)
			if err != nil {
				t.Errorf("varchar JSON deserialization failed: %v", err)
				return
			}

			if value != result {
				t.Errorf("varchar JSON round-trip failed: expected %v, got %v", value, result)
			}
		})

		t.Run("text", func(t *testing.T) {
			jsonData, err := text.SerializeJSON(value)
			if err != nil {
				t.Errorf("text JSON serialization failed: %v", err)
				return
			}

			result, err := text.DeserializeJSON(jsonData)
			if err != nil {
				t.Errorf("text JSON deserialization failed: %v", err)
				return
			}

			if value != result {
				t.Errorf("text JSON round-trip failed: expected %v, got %v", value, result)
			}
		})
	}
}

func TestStringTypes_TextSerializationCompatibility(t *testing.T) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")
	text := types.NewTextType("utf8_general_ci")

	testValues := []interface{}{
		"hello",
		"café",
		"string with 'quotes'",
		nil,
	}

	for _, value := range testValues {
		t.Run("varchar", func(t *testing.T) {
			textStr, err := varchar.SerializeText(value)
			if err != nil {
				t.Errorf("varchar text serialization failed: %v", err)
				return
			}

			result, err := varchar.DeserializeText(textStr)
			if err != nil {
				t.Errorf("varchar text deserialization failed: %v", err)
				return
			}

			if value != result {
				t.Errorf("varchar text round-trip failed: expected %v, got %v", value, result)
			}
		})

		t.Run("text", func(t *testing.T) {
			textStr, err := text.SerializeText(value)
			if err != nil {
				t.Errorf("text text serialization failed: %v", err)
				return
			}

			result, err := text.DeserializeText(textStr)
			if err != nil {
				t.Errorf("text text deserialization failed: %v", err)
				return
			}

			if value != result {
				t.Errorf("text text round-trip failed: expected %v, got %v", value, result)
			}
		})
	}
}

func TestStringTypes_PerformanceComparison(t *testing.T) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")
	text := types.NewTextType("utf8_general_ci")

	// Test that both types have similar performance characteristics
	testValue := "hello world"

	// Measure serialization performance
	varcharData, err := varchar.Serialize(testValue)
	if err != nil {
		t.Fatalf("varchar serialization failed: %v", err)
	}

	textData, err := text.Serialize(testValue)
	if err != nil {
		t.Fatalf("text serialization failed: %v", err)
	}

	// Both should produce similar sized data for same input
	if len(varcharData) != len(textData) {
		t.Errorf("data size mismatch: varchar=%d, text=%d", len(varcharData), len(textData))
	}

	// Both should have same size calculation
	varcharSize := varchar.Size(testValue)
	textSize := text.Size(testValue)

	if varcharSize != textSize {
		t.Errorf("size calculation mismatch: varchar=%d, text=%d", varcharSize, textSize)
	}
}
