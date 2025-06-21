package types_test

import (
	"strings"
	"testing"

	"github.com/khalid64/kongdb/internal/types"
)

func BenchmarkVarcharType_Serialize(b *testing.B) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")
	value := "hello world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := varchar.Serialize(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVarcharType_Deserialize(b *testing.B) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")
	value := "hello world"
	data, _ := varchar.Serialize(value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := varchar.Deserialize(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVarcharType_Compare(b *testing.B) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")
	a := "hello world"
	b_val := "hello universe"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		varchar.Compare(a, b_val)
	}
}

func BenchmarkVarcharType_Validate(b *testing.B) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")
	value := "hello world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := varchar.Validate(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVarcharType_LargeString(b *testing.B) {
	varchar := types.NewVarcharType(1000, "utf8_general_ci")
	value := strings.Repeat("hello world ", 50) // ~600 characters

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := varchar.Serialize(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVarcharType_JSONSerialize(b *testing.B) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")
	value := "hello world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := varchar.SerializeJSON(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVarcharType_JSONDeserialize(b *testing.B) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")
	value := "hello world"
	data, _ := varchar.SerializeJSON(value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := varchar.DeserializeJSON(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVarcharType_TextSerialize(b *testing.B) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")
	value := "hello world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := varchar.SerializeText(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVarcharType_TextDeserialize(b *testing.B) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")
	value := "hello world"
	textStr, _ := varchar.SerializeText(value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := varchar.DeserializeText(textStr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVarcharType_UnicodeString(b *testing.B) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")
	value := "café über naïve"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := varchar.Serialize(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVarcharType_CaseInsensitiveCompare(b *testing.B) {
	varchar := types.NewVarcharType(255, "utf8_general_ci")
	a := "Hello World"
	b_val := "hello world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		varchar.Compare(a, b_val)
	}
}

func BenchmarkVarcharType_CaseSensitiveCompare(b *testing.B) {
	varchar := types.NewVarcharType(255, "utf8_bin")
	a := "Hello World"
	b_val := "hello world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		varchar.Compare(a, b_val)
	}
}
