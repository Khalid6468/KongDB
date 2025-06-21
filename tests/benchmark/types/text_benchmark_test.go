package types_test

import (
	"strings"
	"testing"

	"github.com/khalid64/kongdb/internal/types"
)

func BenchmarkTextType_Serialize(b *testing.B) {
	text := types.NewTextType("utf8_general_ci")
	value := "hello world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := text.Serialize(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextType_Deserialize(b *testing.B) {
	text := types.NewTextType("utf8_general_ci")
	value := "hello world"
	data, _ := text.Serialize(value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := text.Deserialize(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextType_Compare(b *testing.B) {
	text := types.NewTextType("utf8_general_ci")
	a := "hello world"
	b_val := "hello universe"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		text.Compare(a, b_val)
	}
}

func BenchmarkTextType_Validate(b *testing.B) {
	text := types.NewTextType("utf8_general_ci")
	value := "hello world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := text.Validate(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextType_LargeString(b *testing.B) {
	text := types.NewTextType("utf8_general_ci")
	value := strings.Repeat("hello world ", 1000) // ~12KB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := text.Serialize(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextType_VeryLargeString(b *testing.B) {
	text := types.NewTextType("utf8_general_ci")
	value := strings.Repeat("hello world ", 5000) // ~60KB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := text.Serialize(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextType_JSONSerialize(b *testing.B) {
	text := types.NewTextType("utf8_general_ci")
	value := "hello world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := text.SerializeJSON(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextType_JSONDeserialize(b *testing.B) {
	text := types.NewTextType("utf8_general_ci")
	value := "hello world"
	data, _ := text.SerializeJSON(value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := text.DeserializeJSON(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextType_TextSerialize(b *testing.B) {
	text := types.NewTextType("utf8_general_ci")
	value := "hello world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := text.SerializeText(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextType_TextDeserialize(b *testing.B) {
	text := types.NewTextType("utf8_general_ci")
	value := "hello world"
	textStr, _ := text.SerializeText(value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := text.DeserializeText(textStr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextType_UnicodeString(b *testing.B) {
	text := types.NewTextType("utf8_general_ci")
	value := "café über naïve"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := text.Serialize(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextType_LargeUnicodeString(b *testing.B) {
	text := types.NewTextType("utf8_general_ci")
	value := strings.Repeat("café über naïve ", 1000) // ~18KB of Unicode

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := text.Serialize(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextType_CaseInsensitiveCompare(b *testing.B) {
	text := types.NewTextType("utf8_general_ci")
	a := "Hello World"
	b_val := "hello world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		text.Compare(a, b_val)
	}
}

func BenchmarkTextType_CaseSensitiveCompare(b *testing.B) {
	text := types.NewTextType("utf8_bin")
	a := "Hello World"
	b_val := "hello world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		text.Compare(a, b_val)
	}
}

func BenchmarkTextType_LargeStringCompare(b *testing.B) {
	text := types.NewTextType("utf8_general_ci")
	a := strings.Repeat("hello world ", 100)
	b_val := strings.Repeat("hello universe ", 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		text.Compare(a, b_val)
	}
}

func BenchmarkTextType_MemoryUsage(b *testing.B) {
	text := types.NewTextType("utf8_general_ci")
	value := strings.Repeat("hello world ", 1000)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := text.Serialize(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}
