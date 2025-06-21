package types_test

import (
	"testing"

	"github.com/khalid64/kongdb/internal/types"
	"github.com/stretchr/testify/assert"
)

type MockType struct {
	metadata *types.TypeMetadata
}

func (mt *MockType) Metadata() *types.TypeMetadata {
	return mt.metadata
}

func (mt *MockType) Size(value interface{}) int {
	if mt.metadata.Size > 0 {
		return mt.metadata.Size
	}

	if str, ok := value.(string); ok {
		return len(str)
	}

	return 8 // Default mock size
}

func (mt *MockType) Validate(value interface{}) error {
	if value == nil && mt.metadata.Size > 0 {
		return types.NewValidationError("value cannot be null")
	}
	return nil
}

func (mt *MockType) IsNull(value interface{}) bool {
	return value == nil
}

func (mt *MockType) Serialize(value interface{}) ([]byte, error) {
	if mt.IsNull(value) {
		return []byte{0xFF}, nil
	}

	if str, ok := value.(string); ok {
		return []byte(str), nil
	}

	return []byte("mock"), nil
}

func (mt *MockType) NullValue() interface{} {
	return nil
}

func (mt *MockType) Deserialize(data []byte) (interface{}, error) {
	if len(data) == 1 && data[0] == 0xFF {
		return mt.NullValue(), nil
	}
	return string(data), nil
}

func (mt *MockType) SerializeJSON(value interface{}) ([]byte, error) {
	if mt.IsNull(value) {
		return []byte("null"), nil
	}

	return []byte(`"` + value.(string) + `"`), nil
}

func (mt *MockType) DeserializeJSON(data []byte) (interface{}, error) {
	if string(data) == "null" {
		return mt.NullValue(), nil
	}

	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		return string(data[1 : len(data)-1]), nil
	}

	return string(data), nil
}

func (mt *MockType) SerializeText(value interface{}) (string, error) {
	if mt.IsNull(value) {
		return "NULL", nil
	}
	return value.(string), nil
}

func (mt *MockType) DeserializeText(text string) (interface{}, error) {
	if text == "NULL" {
		return mt.NullValue(), nil
	}
	return text, nil
}

func (mt *MockType) Compare(a, b interface{}) int {
	if mt.IsNull(a) && mt.IsNull(b) {
		return 0
	}

	if mt.IsNull(a) {
		return -1
	}

	if mt.IsNull(b) {
		return 1
	}

	strA := a.(string)
	strB := b.(string)
	if strA < strB {
		return -1
	}
	if strA > strB {
		return 1
	}
	return 0
}

func (mt *MockType) String() string {
	return mt.metadata.Name
}

// Test TypeRegistry
func TestTypeRegistry(t *testing.T) {
	registry := types.NewTypeRegistry()

	assert.Equal(t, 0, len(registry.ListTypes()))
	assert.Equal(t, 0, len(registry.ListTypeNames()))

	mockType1 := &MockType{
		metadata: &types.TypeMetadata{
			ID:        types.TypeIDInt32,
			Name:      "INT32",
			Size:      4,
			Alignment: 4,
		},
	}

	mockType2 := &MockType{
		metadata: &types.TypeMetadata{
			ID:        types.TypeIDVarchar,
			Name:      "VARCHAR",
			Size:      -1,
			Alignment: 1,
		},
	}

	err := registry.Register(mockType1)
	assert.NoError(t, err)

	err = registry.Register(mockType2)
	assert.NoError(t, err)

	typeList := registry.ListTypes()
	assert.Equal(t, 2, len(typeList))

	names := registry.ListTypeNames()
	assert.Equal(t, 2, len(names))
	assert.Contains(t, names, "INT32")
	assert.Contains(t, names, "VARCHAR")

	typ, exists := registry.GetByID(types.TypeIDInt32)
	assert.True(t, exists)
	assert.Equal(t, mockType1, typ)

	typ, exists = registry.GetByID(types.TypeIDVarchar)
	assert.True(t, exists)
	assert.Equal(t, mockType2, typ)

	typ, exists = registry.GetByName("INT32")
	assert.True(t, exists)
	assert.Equal(t, typ, mockType1)

	typ, exists = registry.GetByName("VARCHAR")
	assert.True(t, exists)
	assert.Equal(t, mockType2, typ)

	typ, exists = registry.GetByID(999)
	assert.False(t, exists)
	assert.Nil(t, typ)

	typ, exists = registry.GetByName("NONEXISTENT")
	assert.False(t, exists)
	assert.Nil(t, typ)
}

// TestTypeRegistryDuplicateRegistration tests duplicate registration handling
func TestTypeRegistryDuplicateRegistration(t *testing.T) {
	registry := types.NewTypeRegistry()

	mockType1 := &MockType{
		metadata: &types.TypeMetadata{
			ID:   types.TypeIDInt32,
			Name: "INT32",
			Size: 4,
		},
	}

	mockType2 := &MockType{
		metadata: &types.TypeMetadata{
			ID:   types.TypeIDInt32,
			Name: "INT32_DUPLICATE",
			Size: 4,
		},
	}

	mockType3 := &MockType{
		metadata: &types.TypeMetadata{
			ID:   types.TypeIDInt64,
			Name: "INT32",
			Size: 8,
		},
	}

	err := registry.Register(mockType1)
	assert.NoError(t, err)

	err = registry.Register(mockType2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")

	err = registry.Register(mockType3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

// TestTypeRegistryUnregister tests unregistering types
func TestTypeRegistryUnregister(t *testing.T) {
	registry := types.NewTypeRegistry()

	mockType := &MockType{
		metadata: &types.TypeMetadata{
			ID:   types.TypeIDInt32,
			Name: "INT32",
			Size: 4,
		},
	}

	err := registry.Register(mockType)
	assert.NoError(t, err)

	typ, exists := registry.GetByID(types.TypeIDInt32)
	assert.True(t, exists)
	assert.Equal(t, mockType, typ)

	err = registry.Unregister(types.TypeIDInt32)
	assert.NoError(t, err)

	typ, exists = registry.GetByID(types.TypeIDInt32)
	assert.False(t, exists)
	assert.Nil(t, typ)

	err = registry.Unregister(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestGlobalRegistry tests the global registry functionality
func TestGlobalRegistry(t *testing.T) {
	mockType := &MockType{
		metadata: &types.TypeMetadata{
			ID:   types.TypeIDBool,
			Name: "BOOLEAN",
			Size: 1,
		},
	}

	err := types.RegisterGlobal(mockType)
	assert.NoError(t, err)

	typ, exists := types.GetGlobalByID(types.TypeIDBool)
	assert.True(t, exists)
	assert.Equal(t, mockType, typ)

	typ, exists = types.GetGlobalByName("BOOLEAN")
	assert.True(t, exists)
	assert.Equal(t, mockType, typ)

	typeList := types.ListGlobalTypes()
	assert.GreaterOrEqual(t, len(typeList), 1)

	names := types.ListGlobalNames()
	assert.GreaterOrEqual(t, len(names), 1)
	assert.Contains(t, names, "BOOLEAN")
}

// TestTypeUtilityFunctions tests utility functions
func TestTypeUtilityFunctions(t *testing.T) {
	mockType := &MockType{
		metadata: &types.TypeMetadata{
			ID:   types.TypeIDInt32,
			Name: "INT32",
			Size: 4,
		},
	}

	assert.True(t, types.IsFixedSize(mockType))

	variableType := &MockType{
		metadata: &types.TypeMetadata{
			ID:   types.TypeIDVarchar,
			Name: "VARCHAR",
			Size: -1,
		},
	}
	assert.True(t, types.IsVariableSize(variableType))

	assert.Equal(t, 4, types.GetFixedSize(mockType))
	assert.Equal(t, -1, types.GetFixedSize(variableType))
}

// TestSerializeValue tests the SerializeValue convenience function
func TestSerializeValue(t *testing.T) {
	mockType := &MockType{
		metadata: &types.TypeMetadata{
			ID:   types.TypeIDInt32,
			Name: "INT32",
			Size: 4,
		},
	}

	data, err := types.SerializeValue(nil, mockType)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0xFF}, data)

	data, err = types.SerializeValue("test", mockType)
	assert.NoError(t, err)
	assert.Equal(t, []byte("test"), data)
}

// TestDeserializeValue tests the DeserializeValue convenience function
func TestDeserializeValue(t *testing.T) {
	mockType := &MockType{
		metadata: &types.TypeMetadata{
			ID:   types.TypeIDInt32,
			Name: "INT32",
			Size: 4,
		},
	}

	value, err := types.DeserializeValue([]byte{0xFF}, mockType)
	assert.NoError(t, err)
	assert.Nil(t, value)

	value, err = types.DeserializeValue([]byte("test"), mockType)
	assert.NoError(t, err)
	assert.Equal(t, "test", value)
}

// TestTypeIDSerialization tests TypeID serialization/deserialization
func TestTypeIDSerialization(t *testing.T) {
	testCases := []types.TypeID{
		types.TypeIDUnknown,
		types.TypeIDBool,
		types.TypeIDInt32,
		types.TypeIDVarchar,
		types.TypeIDText,
	}

	for _, tc := range testCases {
		data := types.SerializeTypeID(tc)
		assert.Equal(t, 4, len(data))

		id, err := types.DeserializeTypeID(data)
		assert.NoError(t, err)
		assert.Equal(t, tc, id)
	}
}
