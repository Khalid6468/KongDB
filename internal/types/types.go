package types

import (
	"encoding/binary"
	"fmt"
)

type TypeID uint32

const (
	TypeIDUnknown TypeID = iota
	TypeIDNull
	TypeIDBool
	TypeIDInt8
	TypeIDInt16
	TypeIDInt32
	TypeIDInt64
	TypeIDUint8
	TypeIDUint16
	TypeIDUint32
	TypeIDUint64
	TypeIDFloat32
	TypeIDFloat64
	TypeIDText
	TypeIDVarChar
	TypeIDTimestamp
	TypeIDDate
	TypeIDTime
	TypeIDJson
	TypeIDBlob
)

type TypeMetadata struct {
	ID        TypeID
	Name      string
	Size      int
	Alignment int
}

type Type interface {
	Metadata() *TypeMetadata
	Size(value interface{}) int
	Validate(value interface{}) error
	Serialize(value interface{}) ([]byte, error)
	Deserialize(data []byte) (interface{}, error)
	SerializeJSON(value interface{}) ([]byte, error)
	DeserializeJSON(data []byte) (interface{}, error)
	SerializeText(value interface{}) (string, error)
	DeserializeText(text string) (interface{}, error)
	Compare(a, b interface{}) int
	IsNull(value interface{}) bool
	NullValue() interface{}
	String() string
}

type TypeRegistry struct {
	types map[TypeID]Type
	names map[string]TypeID
}

func NewTypeRegistry() *TypeRegistry {
	return &TypeRegistry{
		types: make(map[TypeID]Type),
		names: make(map[string]TypeID),
	}
}

func (tr *TypeRegistry) Register(t Type) error {
	metadata := t.Metadata()

	if existing, exists := tr.types[metadata.ID]; exists {
		return fmt.Errorf("type ID %d already registered as %s", metadata.ID, existing.String())
	}

	if existingID, exists := tr.names[metadata.Name]; exists {
		return fmt.Errorf("type name %s already registered with ID %d", metadata.Name, existingID)
	}

	tr.types[metadata.ID] = t
	tr.names[metadata.Name] = metadata.ID

	return nil
}

func (tr *TypeRegistry) GetByID(id TypeID) (Type, bool) {
	t, exists := tr.types[id]
	return t, exists
}

func (tr *TypeRegistry) GetByName(name string) (Type, bool) {
	id, exists := tr.names[name]
	if !exists {
		return nil, false
	}

	t, exists := tr.types[id]
	return t, exists
}

func (tr *TypeRegistry) ListTypes() []Type {
	types := make([]Type, 0, len(tr.names))
	for _, t := range tr.types {
		types = append(types, t)
	}
	return types
}

func (tr *TypeRegistry) ListTypeNames() []string {
	names := make([]string, 0, len(tr.names))
	for name := range tr.names {
		names = append(names, name)
	}
	return names
}

func (tr *TypeRegistry) Unregister(id TypeID) error {
	t, exists := tr.types[id]
	if !exists {
		return fmt.Errorf("type ID %d not found", id)
	}
	metadata := t.Metadata()
	delete(tr.types, id)
	delete(tr.names, metadata.Name)
	return nil
}

var globalRegistry = NewTypeRegistry()

func RegisterGlobal(t Type) error {
	return globalRegistry.Register(t)
}

func GetGlobalByID(id TypeID) (Type, bool) {
	return globalRegistry.GetByID(id)
}

func GetGlobalByName(name string) (Type, bool) {
	return globalRegistry.GetByName(name)
}

func ListGlobalTypes() []Type {
	return globalRegistry.ListTypes()
}

func ListGlobalNames() []string {
	return globalRegistry.ListTypeNames()
}

func IsFixedSize(t Type) bool {
	return t.Metadata().Size > 0
}

func IsVariableSize(t Type) bool {
	return t.Metadata().Size == -1
}

func GetFixedSize(t Type) int {
	return t.Metadata().Size
}

func SerializeValue(value interface{}, t Type) ([]byte, error) {
	if t.IsNull(value) {
		return []byte{0xFF}, nil
	}

	return t.Serialize(value)
}

func DeserializeValue(data []byte, t Type) (interface{}, error) {
	if len(data) == 1 && data[0] == 0xFF {
		return t.NullValue(), nil
	}

	return t.Deserialize(data)
}

func SerializeTypeID(id TypeID) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(id))
	return buf
}

func DeserializeTypeID(data []byte) (TypeID, error) {
	if len(data) < 4 {
		return TypeIDUnknown, fmt.Errorf("insufficient data for TypeID: need 4 bytes, got %d", len(data))
	}

	return TypeID(binary.BigEndian.Uint32(data[:4])), nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) error {
	return &ValidationError{Message: message}
}
