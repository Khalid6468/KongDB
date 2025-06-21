package types_test

import (
	"testing"

	"github.com/khalid64/kongdb/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegerTypes_Validation(t *testing.T) {
	tests := []struct {
		name    string
		typeObj types.Type
		valid   []interface{}
		invalid []interface{}
	}{
		{
			name:    "INT8",
			typeObj: &types.Int8Type{},
			valid:   []interface{}{int8(-128), int8(0), int8(127), nil},
			invalid: []interface{}{int16(-129), int16(128), int16(0), "foo"},
		},
		{
			name:    "INT16",
			typeObj: &types.Int16Type{},
			valid:   []interface{}{int16(-32768), int16(0), int16(32767), nil},
			invalid: []interface{}{int32(-32769), int32(32768), int32(0), "foo"},
		},
		{
			name:    "INT32",
			typeObj: &types.Int32Type{},
			valid:   []interface{}{int32(-2147483648), int32(0), int32(2147483647), nil},
			invalid: []interface{}{int64(-2147483649), int64(2147483648), int64(0), "foo"},
		},
		{
			name:    "INT64",
			typeObj: &types.Int64Type{},
			valid:   []interface{}{int64(-9223372036854775808), int64(0), int64(9223372036854775807), nil},
			invalid: []interface{}{float64(0), "foo"},
		},
		{
			name:    "UINT8",
			typeObj: &types.Uint8Type{},
			valid:   []interface{}{uint8(0), uint8(255), nil},
			invalid: []interface{}{uint16(256), int8(-1), "foo"},
		},
		{
			name:    "UINT16",
			typeObj: &types.Uint16Type{},
			valid:   []interface{}{uint16(0), uint16(65535), nil},
			invalid: []interface{}{uint32(65536), int16(-1), "foo"},
		},
		{
			name:    "UINT32",
			typeObj: &types.Uint32Type{},
			valid:   []interface{}{uint32(0), uint32(4294967295), nil},
			invalid: []interface{}{uint64(4294967296), int32(-1), "foo"},
		},
		{
			name:    "UINT64",
			typeObj: &types.Uint64Type{},
			valid:   []interface{}{uint64(0), uint64(18446744073709551615), nil},
			invalid: []interface{}{int64(-1), "foo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			for _, v := range tt.valid {
				assert.NoError(t, tt.typeObj.Validate(v), "should be valid: %v", v)
			}
		})
		t.Run(tt.name+"_invalid", func(t *testing.T) {
			for _, v := range tt.invalid {
				assert.Error(t, tt.typeObj.Validate(v), "should be invalid: %v", v)
			}
		})
	}
}

func TestIntegerTypes_Serialization(t *testing.T) {
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
		t.Run(tt.name+"_binary", func(t *testing.T) {
			data, err := tt.typeObj.Serialize(tt.value)
			require.NoError(t, err)
			v2, err := tt.typeObj.Deserialize(data)
			require.NoError(t, err)
			assert.Equal(t, tt.value, v2)
		})
		t.Run(tt.name+"_null", func(t *testing.T) {
			data, err := tt.typeObj.Serialize(nil)
			require.NoError(t, err)
			v2, err := tt.typeObj.Deserialize(data)
			require.NoError(t, err)
			assert.Nil(t, v2)
		})
		t.Run(tt.name+"_wrong_type", func(t *testing.T) {
			_, err := tt.typeObj.Serialize("wrong type")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "value is not")
		})
	}
}

func TestIntegerTypes_JSONTextSerialization(t *testing.T) {
	tests := []struct {
		name    string
		typeObj types.Type
		value   interface{}
		jsonStr string
		textStr string
	}{
		{"INT8", &types.Int8Type{}, int8(42), "42", "42"},
		{"INT16", &types.Int16Type{}, int16(12345), "12345", "12345"},
		{"INT32", &types.Int32Type{}, int32(123456789), "123456789", "123456789"},
		{"INT64", &types.Int64Type{}, int64(1234567890123456789), "1234567890123456789", "1234567890123456789"},
		{"UINT8", &types.Uint8Type{}, uint8(200), "200", "200"},
		{"UINT16", &types.Uint16Type{}, uint16(40000), "40000", "40000"},
		{"UINT32", &types.Uint32Type{}, uint32(3000000000), "3000000000", "3000000000"},
		{"UINT64", &types.Uint64Type{}, uint64(9000000000000000000), "9000000000000000000", "9000000000000000000"},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_json", func(t *testing.T) {
			data, err := tt.typeObj.SerializeJSON(tt.value)
			require.NoError(t, err)
			assert.JSONEq(t, tt.jsonStr, string(data))
			v2, err := tt.typeObj.DeserializeJSON(data)
			require.NoError(t, err)
			assert.Equal(t, tt.value, v2)
		})
		t.Run(tt.name+"_text", func(t *testing.T) {
			text, err := tt.typeObj.SerializeText(tt.value)
			require.NoError(t, err)
			assert.Equal(t, tt.textStr, text)
			v2, err := tt.typeObj.DeserializeText(text)
			require.NoError(t, err)
			assert.Equal(t, tt.value, v2)
		})
		t.Run(tt.name+"_json_null", func(t *testing.T) {
			data, err := tt.typeObj.SerializeJSON(nil)
			require.NoError(t, err)
			assert.Equal(t, "null", string(data))
			v2, err := tt.typeObj.DeserializeJSON(data)
			require.NoError(t, err)
			assert.Nil(t, v2)
		})
		t.Run(tt.name+"_text_null", func(t *testing.T) {
			text, err := tt.typeObj.SerializeText(nil)
			require.NoError(t, err)
			assert.Equal(t, "NULL", text)
			v2, err := tt.typeObj.DeserializeText(text)
			require.NoError(t, err)
			assert.Nil(t, v2)
		})
	}
}

func TestIntegerTypes_Compare(t *testing.T) {
	tests := []struct {
		name    string
		typeObj types.Type
		a, b    interface{}
		expect  int
	}{
		// INT8 tests
		{"INT8_eq", &types.Int8Type{}, int8(1), int8(1), 0},
		{"INT8_lt", &types.Int8Type{}, int8(1), int8(2), -1},
		{"INT8_gt", &types.Int8Type{}, int8(2), int8(1), 1},
		{"INT8_null", &types.Int8Type{}, nil, int8(1), -1},
		{"INT8_null2", &types.Int8Type{}, int8(1), nil, 1},
		{"INT8_null_eq", &types.Int8Type{}, nil, nil, 0},
		{"INT8_min", &types.Int8Type{}, int8(-128), int8(-127), -1},
		{"INT8_max", &types.Int8Type{}, int8(127), int8(126), 1},

		// INT16 tests
		{"INT16_eq", &types.Int16Type{}, int16(1), int16(1), 0},
		{"INT16_lt", &types.Int16Type{}, int16(1), int16(2), -1},
		{"INT16_gt", &types.Int16Type{}, int16(2), int16(1), 1},
		{"INT16_null", &types.Int16Type{}, nil, int16(1), -1},
		{"INT16_null2", &types.Int16Type{}, int16(1), nil, 1},
		{"INT16_null_eq", &types.Int16Type{}, nil, nil, 0},
		{"INT16_min", &types.Int16Type{}, int16(-32768), int16(-32767), -1},
		{"INT16_max", &types.Int16Type{}, int16(32767), int16(32766), 1},

		// INT32 tests
		{"INT32_eq", &types.Int32Type{}, int32(1), int32(1), 0},
		{"INT32_lt", &types.Int32Type{}, int32(1), int32(2), -1},
		{"INT32_gt", &types.Int32Type{}, int32(2), int32(1), 1},
		{"INT32_null", &types.Int32Type{}, nil, int32(1), -1},
		{"INT32_null2", &types.Int32Type{}, int32(1), nil, 1},
		{"INT32_null_eq", &types.Int32Type{}, nil, nil, 0},
		{"INT32_min", &types.Int32Type{}, int32(-2147483648), int32(-2147483647), -1},
		{"INT32_max", &types.Int32Type{}, int32(2147483647), int32(2147483646), 1},

		// INT64 tests
		{"INT64_eq", &types.Int64Type{}, int64(1), int64(1), 0},
		{"INT64_lt", &types.Int64Type{}, int64(1), int64(2), -1},
		{"INT64_gt", &types.Int64Type{}, int64(2), int64(1), 1},
		{"INT64_null", &types.Int64Type{}, nil, int64(1), -1},
		{"INT64_null2", &types.Int64Type{}, int64(1), nil, 1},
		{"INT64_null_eq", &types.Int64Type{}, nil, nil, 0},
		{"INT64_min", &types.Int64Type{}, int64(-9223372036854775808), int64(-9223372036854775807), -1},
		{"INT64_max", &types.Int64Type{}, int64(9223372036854775807), int64(9223372036854775806), 1},

		// UINT8 tests
		{"UINT8_eq", &types.Uint8Type{}, uint8(1), uint8(1), 0},
		{"UINT8_lt", &types.Uint8Type{}, uint8(1), uint8(2), -1},
		{"UINT8_gt", &types.Uint8Type{}, uint8(2), uint8(1), 1},
		{"UINT8_null", &types.Uint8Type{}, nil, uint8(1), -1},
		{"UINT8_null2", &types.Uint8Type{}, uint8(1), nil, 1},
		{"UINT8_null_eq", &types.Uint8Type{}, nil, nil, 0},
		{"UINT8_min", &types.Uint8Type{}, uint8(0), uint8(1), -1},
		{"UINT8_max", &types.Uint8Type{}, uint8(255), uint8(254), 1},

		// UINT16 tests
		{"UINT16_eq", &types.Uint16Type{}, uint16(1), uint16(1), 0},
		{"UINT16_lt", &types.Uint16Type{}, uint16(1), uint16(2), -1},
		{"UINT16_gt", &types.Uint16Type{}, uint16(2), uint16(1), 1},
		{"UINT16_null", &types.Uint16Type{}, nil, uint16(1), -1},
		{"UINT16_null2", &types.Uint16Type{}, uint16(1), nil, 1},
		{"UINT16_null_eq", &types.Uint16Type{}, nil, nil, 0},
		{"UINT16_min", &types.Uint16Type{}, uint16(0), uint16(1), -1},
		{"UINT16_max", &types.Uint16Type{}, uint16(65535), uint16(65534), 1},

		// UINT32 tests
		{"UINT32_eq", &types.Uint32Type{}, uint32(1), uint32(1), 0},
		{"UINT32_lt", &types.Uint32Type{}, uint32(1), uint32(2), -1},
		{"UINT32_gt", &types.Uint32Type{}, uint32(2), uint32(1), 1},
		{"UINT32_null", &types.Uint32Type{}, nil, uint32(1), -1},
		{"UINT32_null2", &types.Uint32Type{}, uint32(1), nil, 1},
		{"UINT32_null_eq", &types.Uint32Type{}, nil, nil, 0},
		{"UINT32_min", &types.Uint32Type{}, uint32(0), uint32(1), -1},
		{"UINT32_max", &types.Uint32Type{}, uint32(4294967295), uint32(4294967294), 1},

		// UINT64 tests
		{"UINT64_eq", &types.Uint64Type{}, uint64(1), uint64(1), 0},
		{"UINT64_lt", &types.Uint64Type{}, uint64(1), uint64(2), -1},
		{"UINT64_gt", &types.Uint64Type{}, uint64(2), uint64(1), 1},
		{"UINT64_null", &types.Uint64Type{}, nil, uint64(1), -1},
		{"UINT64_null2", &types.Uint64Type{}, uint64(1), nil, 1},
		{"UINT64_null_eq", &types.Uint64Type{}, nil, nil, 0},
		{"UINT64_min", &types.Uint64Type{}, uint64(0), uint64(1), -1},
		{"UINT64_max", &types.Uint64Type{}, uint64(18446744073709551615), uint64(18446744073709551614), 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.typeObj.Compare(tt.a, tt.b))
		})
	}
}

func TestIntegerTypes_UtilityMethods(t *testing.T) {
	tests := []struct {
		name    string
		typeObj types.Type
		value   interface{}
		size    int
	}{
		{"INT8", &types.Int8Type{}, int8(42), 1},
		{"INT16", &types.Int16Type{}, int16(12345), 2},
		{"INT32", &types.Int32Type{}, int32(123456789), 4},
		{"INT64", &types.Int64Type{}, int64(1234567890123456789), 8},
		{"UINT8", &types.Uint8Type{}, uint8(200), 1},
		{"UINT16", &types.Uint16Type{}, uint16(40000), 2},
		{"UINT32", &types.Uint32Type{}, uint32(3000000000), 4},
		{"UINT64", &types.Uint64Type{}, uint64(9000000000000000000), 8},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_isnull", func(t *testing.T) {
			assert.False(t, tt.typeObj.IsNull(tt.value))
			assert.True(t, tt.typeObj.IsNull(nil))
		})
		t.Run(tt.name+"_nullvalue", func(t *testing.T) {
			assert.Nil(t, tt.typeObj.NullValue())
		})
		t.Run(tt.name+"_string", func(t *testing.T) {
			assert.Equal(t, tt.name, tt.typeObj.String())
		})
		t.Run(tt.name+"_size", func(t *testing.T) {
			assert.Equal(t, tt.size, tt.typeObj.Size(tt.value))
		})
		t.Run(tt.name+"_metadata", func(t *testing.T) {
			metadata := tt.typeObj.Metadata()
			assert.NotNil(t, metadata)
			assert.Equal(t, tt.name, metadata.Name)
			assert.Equal(t, tt.size, metadata.Size)
		})
	}
}

func TestIntegerTypes_ErrorCases(t *testing.T) {
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
		t.Run(tt.name+"_invalid_json", func(t *testing.T) {
			_, err := tt.typeObj.DeserializeJSON([]byte("invalid json"))
			assert.Error(t, err)
		})
		t.Run(tt.name+"_invalid_text", func(t *testing.T) {
			_, err := tt.typeObj.DeserializeText("not a number")
			assert.Error(t, err)
		})
		t.Run(tt.name+"_insufficient_data", func(t *testing.T) {
			_, err := tt.typeObj.Deserialize([]byte{})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not enough data")
		})
	}
}
