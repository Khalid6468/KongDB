package types_test

import (
	"testing"

	"github.com/khalid64/kongdb/internal/types"
)

func TestVarcharType_Metadata(t *testing.T) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")
	metadata := varchar.Metadata()

	if metadata.ID != types.TypeIDVarchar {
		t.Errorf("expected TypeIDVarchar, got %d", metadata.ID)
	}

	if metadata.Name != "VARCHAR(255)" {
		t.Errorf("expected VARCHAR(255), got %s", metadata.Name)
	}

	if metadata.Size != -1 {
		t.Errorf("expected -1 for variable size, got %d", metadata.Size)
	}
}

func TestVarcharType_Validate(t *testing.T) {
	varchar := types.NewVarcharType(10, "utf8_general_ci")

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"valid string", "hello", false},
		{"empty string", "", false},
		{"max length", "1234567890", false},
		{"too long", "12345678901", true},
		{"null value", nil, false},
		{"wrong type", 123, true},
		{"unicode string", "café", false},
		{"unicode too long", "caféééééééééé", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := varchar.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVarcharType_Serialize(t *testing.T) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"normal string", "hello", false},
		{"empty string", "", false},
		{"unicode string", "café", false},
		{"null value", nil, false},
		{"invalid value", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := varchar.Serialize(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Test deserialization
				result, err := varchar.Deserialize(data)
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

func TestVarcharType_JSON(t *testing.T) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")

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
			data, err := varchar.SerializeJSON(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				result, err := varchar.DeserializeJSON(data)
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

func TestVarcharType_Compare(t *testing.T) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := varchar.Compare(tt.a, tt.b)
			if result != tt.want {
				t.Errorf("Compare() = %d, want %d", result, tt.want)
			}
		})
	}
}

func TestVarcharType_Size(t *testing.T) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")

	tests := []struct {
		name  string
		value interface{}
		want  int
	}{
		{"empty string", "", 4},
		{"hello", "hello", 9},
		{"unicode", "café", 9},
		{"null", nil, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := varchar.Size(tt.value)
			if result != tt.want {
				t.Errorf("Size() = %d, want %d", result, tt.want)
			}
		})
	}
}
