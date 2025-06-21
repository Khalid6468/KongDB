package types_test

import (
	"strings"
	"testing"

	"github.com/khalid64/kongdb/internal/types"
)

func TestTextType_Metadata(t *testing.T) {
	text := types.NewTextType("utf8_general_ci")
	metadata := text.Metadata()

	if metadata.ID != types.TypeIDText {
		t.Errorf("expected TypeIDText, got %d", metadata.ID)
	}

	if metadata.Name != "TEXT" {
		t.Errorf("expected TEXT, got %s", metadata.Name)
	}

	if metadata.Size != -1 {
		t.Errorf("expected -1 for variable size, got %d", metadata.Size)
	}
}

func TestTextType_Validate(t *testing.T) {
	text := types.NewTextType("utf8_general_ci")

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"valid string", "hello", false},
		{"empty string", "", false},
		{"large string", strings.Repeat("a", 1000), false},
		{"unicode string", "café", false},
		{"null value", nil, false},
		{"wrong type", 123, true},
		{"too large", strings.Repeat("a", 70000), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := text.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTextType_Serialize(t *testing.T) {
	text := types.NewTextType("utf8_general_ci")

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"normal string", "hello", false},
		{"empty string", "", false},
		{"unicode string", "café", false},
		{"large string", strings.Repeat("a", 1000), false},
		{"null value", nil, false},
		{"invalid value", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := text.Serialize(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Test deserialization
				result, err := text.Deserialize(data)
				if err != nil {
					t.Errorf("Deserialize() error = %v", err)
					return
				}

				if tt.value == nil {
					if result != nil {
						t.Errorf("expected nil, got %v", result)
					}
				} else {
					if result != tt.value {
						t.Errorf("expected %v, got %v", tt.value, result)
					}
				}
			}
		})
	}
}

func TestTextType_JSON(t *testing.T) {
	text := types.NewTextType("utf8_general_ci")

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"normal string", "hello", false},
		{"string with quotes", "hello 'world'", false},
		{"unicode string", "café", false},
		{"large string", strings.Repeat("a", 1000), false},
		{"null value", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := text.SerializeJSON(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				result, err := text.DeserializeJSON(data)
				if err != nil {
					t.Errorf("DeserializeJSON() error = %v", err)
					return
				}

				if tt.value == nil {
					if result != nil {
						t.Errorf("expected nil, got %v", result)
					}
				} else {
					if result != tt.value {
						t.Errorf("expected %v, got %v", tt.value, result)
					}
				}
			}
		})
	}
}

func TestTextType_Compare(t *testing.T) {
	text := types.NewTextType("utf8_general_ci")

	tests := []struct {
		name string
		a    interface{}
		b    interface{}
		want int
	}{
		{"hello < world", "hello", "world", -1},
		{"hello = hello", "hello", "hello", 0},
		{"world > hello", "world", "hello", 1},
		{"null < string", nil, "hello", -1},
		{"string > null", "hello", nil, 1},
		{"null = null", nil, nil, 0},
		{"case insensitive", "Hello", "hello", 0},
		{"large strings", strings.Repeat("a", 100), strings.Repeat("b", 100), -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := text.Compare(tt.a, tt.b)
			if result != tt.want {
				t.Errorf("Compare() = %d, want %d", result, tt.want)
			}
		})
	}
}

func TestTextType_Size(t *testing.T) {
	text := types.NewTextType("utf8_general_ci")

	tests := []struct {
		name  string
		value interface{}
		want  int
	}{
		{"empty string", "", 4},
		{"hello", "hello", 9},
		{"unicode", "café", 9},
		{"large string", strings.Repeat("a", 1000), 1004},
		{"null", nil, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := text.Size(tt.value)
			if result != tt.want {
				t.Errorf("Size() = %d, want %d", result, tt.want)
			}
		})
	}
}

func TestTextType_TextSerialization(t *testing.T) {
	text := types.NewTextType("utf8_general_ci")

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"normal string", "hello", false},
		{"string with quotes", "hello 'world'", false},
		{"unicode string", "café", false},
		{"null value", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			textStr, err := text.SerializeText(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				result, err := text.DeserializeText(textStr)
				if err != nil {
					t.Errorf("DeserializeText() error = %v", err)
					return
				}

				if tt.value == nil {
					if result != nil {
						t.Errorf("expected nil, got %v", result)
					}
				} else {
					if result != tt.value {
						t.Errorf("expected %v, got %v", tt.value, result)
					}
				}
			}
		})
	}
}
