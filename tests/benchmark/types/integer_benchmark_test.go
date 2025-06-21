package types_test

import (
	"testing"

	"github.com/khalid64/kongdb/internal/types"
)

func BenchmarkIntegerSerialization(b *testing.B) {
	tests := []struct {
		name    string
		typeObj types.Type
		value   interface{}
	}{
		{"INT8", &types.Int8Type{}, int8(42)},
		{"INT16", &types.Int16Type{}, int16(12345)},
		{"INT32", &types.Int32Type{}, int32(123456789)},
		{"INT64", &types.Int64Type{}, int64(1234567890123456789)},
		{"UINT8", &types.Uint8Type{}, uint8(200)},
		{"UINT16", &types.Uint16Type{}, uint16(40000)},
		{"UINT32", &types.Uint32Type{}, uint32(3000000000)},
		{"UINT64", &types.Uint64Type{}, uint64(9000000000000000000)},
	}

	for _, tt := range tests {
		b.Run(tt.name+"_serialize", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = tt.typeObj.Serialize(tt.value)
			}
		})
		b.Run(tt.name+"_deserialize", func(b *testing.B) {
			data, _ := tt.typeObj.Serialize(tt.value)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = tt.typeObj.Deserialize(data)
			}
		})
	}
}

func BenchmarkIntegerJSONSerialization(b *testing.B) {
	tests := []struct {
		name    string
		typeObj types.Type
		value   interface{}
	}{
		{"INT8", &types.Int8Type{}, int8(42)},
		{"INT16", &types.Int16Type{}, int16(12345)},
		{"INT32", &types.Int32Type{}, int32(123456789)},
		{"INT64", &types.Int64Type{}, int64(1234567890123456789)},
		{"UINT8", &types.Uint8Type{}, uint8(200)},
		{"UINT16", &types.Uint16Type{}, uint16(40000)},
		{"UINT32", &types.Uint32Type{}, uint32(3000000000)},
		{"UINT64", &types.Uint64Type{}, uint64(9000000000000000000)},
	}

	for _, tt := range tests {
		b.Run(tt.name+"_json_serialize", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = tt.typeObj.SerializeJSON(tt.value)
			}
		})
		b.Run(tt.name+"_json_deserialize", func(b *testing.B) {
			data, _ := tt.typeObj.SerializeJSON(tt.value)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = tt.typeObj.DeserializeJSON(data)
			}
		})
	}
}

func BenchmarkIntegerTextSerialization(b *testing.B) {
	tests := []struct {
		name    string
		typeObj types.Type
		value   interface{}
	}{
		{"INT8", &types.Int8Type{}, int8(42)},
		{"INT16", &types.Int16Type{}, int16(12345)},
		{"INT32", &types.Int32Type{}, int32(123456789)},
		{"INT64", &types.Int64Type{}, int64(1234567890123456789)},
		{"UINT8", &types.Uint8Type{}, uint8(200)},
		{"UINT16", &types.Uint16Type{}, uint16(40000)},
		{"UINT32", &types.Uint32Type{}, uint32(3000000000)},
		{"UINT64", &types.Uint64Type{}, uint64(9000000000000000000)},
	}

	for _, tt := range tests {
		b.Run(tt.name+"_text_serialize", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = tt.typeObj.SerializeText(tt.value)
			}
		})
		b.Run(tt.name+"_text_deserialize", func(b *testing.B) {
			text, _ := tt.typeObj.SerializeText(tt.value)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = tt.typeObj.DeserializeText(text)
			}
		})
	}
}

func BenchmarkIntegerValidation(b *testing.B) {
	tests := []struct {
		name    string
		typeObj types.Type
		value   interface{}
	}{
		{"INT8", &types.Int8Type{}, int8(42)},
		{"INT16", &types.Int16Type{}, int16(12345)},
		{"INT32", &types.Int32Type{}, int32(123456789)},
		{"INT64", &types.Int64Type{}, int64(1234567890123456789)},
		{"UINT8", &types.Uint8Type{}, uint8(200)},
		{"UINT16", &types.Uint16Type{}, uint16(40000)},
		{"UINT32", &types.Uint32Type{}, uint32(3000000000)},
		{"UINT64", &types.Uint64Type{}, uint64(9000000000000000000)},
	}

	for _, tt := range tests {
		b.Run(tt.name+"_validate", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tt.typeObj.Validate(tt.value)
			}
		})
		b.Run(tt.name+"_validate_null", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tt.typeObj.Validate(nil)
			}
		})
	}
}

func BenchmarkIntegerComparison(b *testing.B) {
	tests := []struct {
		name    string
		typeObj types.Type
		a, b    interface{}
	}{
		{"INT8", &types.Int8Type{}, int8(42), int8(100)},
		{"INT16", &types.Int16Type{}, int16(12345), int16(20000)},
		{"INT32", &types.Int32Type{}, int32(123456789), int32(987654321)},
		{"INT64", &types.Int64Type{}, int64(1234567890123456789), int64(987654321098765432)},
		{"UINT8", &types.Uint8Type{}, uint8(200), uint8(100)},
		{"UINT16", &types.Uint16Type{}, uint16(40000), uint16(20000)},
		{"UINT32", &types.Uint32Type{}, uint32(3000000000), uint32(2000000000)},
		{"UINT64", &types.Uint64Type{}, uint64(9000000000000000000), uint64(8000000000000000000)},
	}

	for _, tt := range tests {
		b.Run(tt.name+"_compare", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tt.typeObj.Compare(tt.a, tt.b)
			}
		})
		b.Run(tt.name+"_compare_eq", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tt.typeObj.Compare(tt.a, tt.a)
			}
		})
		b.Run(tt.name+"_compare_null", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tt.typeObj.Compare(nil, tt.b)
			}
		})
	}
}

func BenchmarkIntegerUtilityMethods(b *testing.B) {
	tests := []struct {
		name    string
		typeObj types.Type
		value   interface{}
	}{
		{"INT8", &types.Int8Type{}, int8(42)},
		{"INT16", &types.Int16Type{}, int16(12345)},
		{"INT32", &types.Int32Type{}, int32(123456789)},
		{"INT64", &types.Int64Type{}, int64(1234567890123456789)},
		{"UINT8", &types.Uint8Type{}, uint8(200)},
		{"UINT16", &types.Uint16Type{}, uint16(40000)},
		{"UINT32", &types.Uint32Type{}, uint32(3000000000)},
		{"UINT64", &types.Uint64Type{}, uint64(9000000000000000000)},
	}

	for _, tt := range tests {
		b.Run(tt.name+"_isnull", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tt.typeObj.IsNull(tt.value)
			}
		})
		b.Run(tt.name+"_isnull_null", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tt.typeObj.IsNull(nil)
			}
		})
		b.Run(tt.name+"_size", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tt.typeObj.Size(tt.value)
			}
		})
		b.Run(tt.name+"_string", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tt.typeObj.String()
			}
		})
		b.Run(tt.name+"_metadata", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tt.typeObj.Metadata()
			}
		})
	}
}

func BenchmarkIntegerBoundaryValues(b *testing.B) {
	tests := []struct {
		name    string
		typeObj types.Type
		value   interface{}
	}{
		{"INT8_MIN", &types.Int8Type{}, int8(-128)},
		{"INT8_MAX", &types.Int8Type{}, int8(127)},
		{"INT16_MIN", &types.Int16Type{}, int16(-32768)},
		{"INT16_MAX", &types.Int16Type{}, int16(32767)},
		{"INT32_MIN", &types.Int32Type{}, int32(-2147483648)},
		{"INT32_MAX", &types.Int32Type{}, int32(2147483647)},
		{"INT64_MIN", &types.Int64Type{}, int64(-9223372036854775808)},
		{"INT64_MAX", &types.Int64Type{}, int64(9223372036854775807)},
		{"UINT8_MAX", &types.Uint8Type{}, uint8(255)},
		{"UINT16_MAX", &types.Uint16Type{}, uint16(65535)},
		{"UINT32_MAX", &types.Uint32Type{}, uint32(4294967295)},
		{"UINT64_MAX", &types.Uint64Type{}, uint64(18446744073709551615)},
	}

	for _, tt := range tests {
		b.Run(tt.name+"_serialize", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = tt.typeObj.Serialize(tt.value)
			}
		})
		b.Run(tt.name+"_validate", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tt.typeObj.Validate(tt.value)
			}
		})
	}
}

func BenchmarkIntegerNullHandling(b *testing.B) {
	tests := []struct {
		name    string
		typeObj types.Type
	}{
		{"INT8", &types.Int8Type{}},
		{"INT16", &types.Int16Type{}},
		{"INT32", &types.Int32Type{}},
		{"INT64", &types.Int64Type{}},
		{"UINT8", &types.Uint8Type{}},
		{"UINT16", &types.Uint16Type{}},
		{"UINT32", &types.Uint32Type{}},
		{"UINT64", &types.Uint64Type{}},
	}

	for _, tt := range tests {
		b.Run(tt.name+"_serialize_null", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = tt.typeObj.Serialize(nil)
			}
		})
		b.Run(tt.name+"_json_serialize_null", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = tt.typeObj.SerializeJSON(nil)
			}
		})
		b.Run(tt.name+"_text_serialize_null", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = tt.typeObj.SerializeText(nil)
			}
		})
	}
}
