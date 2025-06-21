package types

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
)

type Int8Type struct{}
type Int16Type struct{}
type Int32Type struct{}
type Int64Type struct{}
type Uint8Type struct{}
type Uint16Type struct{}
type Uint32Type struct{}
type Uint64Type struct{}

func (t *Int8Type) Metadata() *TypeMetadata {
	return &TypeMetadata{
		ID:        TypeIDInt8,
		Name:      "INT8",
		Size:      1,
		Alignment: 1,
	}
}

func (t *Int8Type) Size(value interface{}) int {
	return 1
}

func (t *Int8Type) Validate(value interface{}) error {
	if value == nil {
		return nil
	}
	_, ok := value.(int8)
	if !ok {
		return NewValidationError("value is not int8")
	}
	// No need to check range, int8 covers all possible values
	return nil
}

func (t *Int8Type) Serialize(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte{0xFF}, nil
	}
	v, ok := value.(int8)
	if !ok {
		return nil, fmt.Errorf("value is not int8")
	}
	return []byte{byte(v)}, nil
}

func (t *Int8Type) Deserialize(data []byte) (interface{}, error) {
	if len(data) == 1 && data[0] == 0xFF {
		return t.NullValue(), nil
	}
	if len(data) < 1 {
		return nil, fmt.Errorf("not enough data for int8")
	}
	return int8(data[0]), nil
}

func (t *Int8Type) SerializeJSON(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(value)
}

func (t *Int8Type) DeserializeJSON(data []byte) (interface{}, error) {
	if string(data) == "null" {
		return t.NullValue(), nil
	}
	var v int8
	err := json.Unmarshal(data, &v)
	return v, err
}

func (t *Int8Type) SerializeText(value interface{}) (string, error) {
	if value == nil {
		return "NULL", nil
	}
	return strconv.FormatInt(int64(value.(int8)), 10), nil
}

func (t *Int8Type) DeserializeText(text string) (interface{}, error) {
	if text == "NULL" {
		return t.NullValue(), nil
	}
	v, err := strconv.ParseInt(text, 10, 8)
	return int8(v), err
}

func (t *Int8Type) Compare(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}
	ai := a.(int8)
	bi := b.(int8)
	if ai < bi {
		return -1
	}
	if ai > bi {
		return 1
	}
	return 0
}

func (t *Int8Type) IsNull(value interface{}) bool {
	return value == nil
}

func (t *Int8Type) NullValue() interface{} {
	return nil
}

func (t *Int8Type) String() string {
	return "INT8"
}

func (t *Int16Type) Metadata() *TypeMetadata {
	return &TypeMetadata{
		ID:        TypeIDInt16,
		Name:      "INT16",
		Size:      2,
		Alignment: 2,
	}
}

func (t *Int16Type) Size(value interface{}) int {
	return 2
}

func (t *Int16Type) Validate(value interface{}) error {
	if value == nil {
		return nil
	}
	_, ok := value.(int16)
	if !ok {
		return NewValidationError("value is not int16")
	}
	// No need to check range, int16 covers all possible values
	return nil
}

func (t *Int16Type) Serialize(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte{0xFF, 0xFF}, nil
	}
	v, ok := value.(int16)
	if !ok {
		return nil, fmt.Errorf("value is not int16")
	}
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(v))
	return buf, nil
}

func (t *Int16Type) Deserialize(data []byte) (interface{}, error) {
	if len(data) == 2 && data[0] == 0xFF && data[1] == 0xFF {
		return t.NullValue(), nil
	}
	if len(data) < 2 {
		return nil, fmt.Errorf("not enough data for int16")
	}
	return int16(binary.BigEndian.Uint16(data)), nil
}

func (t *Int16Type) SerializeJSON(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(value)
}

func (t *Int16Type) DeserializeJSON(data []byte) (interface{}, error) {
	if string(data) == "null" {
		return t.NullValue(), nil
	}
	var v int16
	err := json.Unmarshal(data, &v)
	return v, err
}

func (t *Int16Type) SerializeText(value interface{}) (string, error) {
	if value == nil {
		return "NULL", nil
	}
	return strconv.FormatInt(int64(value.(int16)), 10), nil
}

func (t *Int16Type) DeserializeText(text string) (interface{}, error) {
	if text == "NULL" {
		return t.NullValue(), nil
	}
	v, err := strconv.ParseInt(text, 10, 16)
	return int16(v), err
}

func (t *Int16Type) Compare(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}
	ai := a.(int16)
	bi := b.(int16)
	if ai < bi {
		return -1
	}
	if ai > bi {
		return 1
	}
	return 0
}

func (t *Int16Type) IsNull(value interface{}) bool {
	return value == nil
}

func (t *Int16Type) NullValue() interface{} {
	return nil
}

func (t *Int16Type) String() string {
	return "INT16"
}

func (t *Int32Type) Metadata() *TypeMetadata {
	return &TypeMetadata{
		ID:        TypeIDInt32,
		Name:      "INT32",
		Size:      4,
		Alignment: 4,
	}
}

func (t *Int32Type) Size(value interface{}) int {
	return 4
}

func (t *Int32Type) Validate(value interface{}) error {
	if value == nil {
		return nil
	}
	_, ok := value.(int32)
	if !ok {
		return NewValidationError("value is not int32")
	}
	// No need to check range, int32 covers all possible values
	return nil
}

func (t *Int32Type) Serialize(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte{0xFF, 0xFF, 0xFF, 0xFF}, nil
	}
	v, ok := value.(int32)
	if !ok {
		return nil, fmt.Errorf("value is not int32")
	}
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(v))
	return buf, nil
}

func (t *Int32Type) Deserialize(data []byte) (interface{}, error) {
	if len(data) == 4 && data[0] == 0xFF && data[1] == 0xFF && data[2] == 0xFF && data[3] == 0xFF {
		return t.NullValue(), nil
	}
	if len(data) < 4 {
		return nil, fmt.Errorf("not enough data for int32")
	}
	return int32(binary.BigEndian.Uint32(data)), nil
}

func (t *Int32Type) SerializeJSON(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(value)
}

func (t *Int32Type) DeserializeJSON(data []byte) (interface{}, error) {
	if string(data) == "null" {
		return t.NullValue(), nil
	}
	var v int32
	err := json.Unmarshal(data, &v)
	return v, err
}

func (t *Int32Type) SerializeText(value interface{}) (string, error) {
	if value == nil {
		return "NULL", nil
	}
	return strconv.FormatInt(int64(value.(int32)), 10), nil
}

func (t *Int32Type) DeserializeText(text string) (interface{}, error) {
	if text == "NULL" {
		return t.NullValue(), nil
	}
	v, err := strconv.ParseInt(text, 10, 32)
	return int32(v), err
}

func (t *Int32Type) Compare(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}
	ai := a.(int32)
	bi := b.(int32)
	if ai < bi {
		return -1
	}
	if ai > bi {
		return 1
	}
	return 0
}

func (t *Int32Type) IsNull(value interface{}) bool {
	return value == nil
}

func (t *Int32Type) NullValue() interface{} {
	return nil
}

func (t *Int32Type) String() string {
	return "INT32"
}

func (t *Int64Type) Metadata() *TypeMetadata {
	return &TypeMetadata{
		ID:        TypeIDInt64,
		Name:      "INT64",
		Size:      8,
		Alignment: 8,
	}
}

func (t *Int64Type) Size(value interface{}) int {
	return 8
}

func (t *Int64Type) Validate(value interface{}) error {
	if value == nil {
		return nil
	}
	_, ok := value.(int64)
	if !ok {
		return NewValidationError("value is not int64")
	}
	// No need to check range, int64 covers all possible values
	return nil
}

func (t *Int64Type) Serialize(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, nil
	}
	v, ok := value.(int64)
	if !ok {
		return nil, fmt.Errorf("value is not int64")
	}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(v))
	return buf, nil
}

func (t *Int64Type) Deserialize(data []byte) (interface{}, error) {
	if len(data) == 8 &&
		data[0] == 0xFF && data[1] == 0xFF && data[2] == 0xFF && data[3] == 0xFF &&
		data[4] == 0xFF && data[5] == 0xFF && data[6] == 0xFF && data[7] == 0xFF {
		return t.NullValue(), nil
	}
	if len(data) < 8 {
		return nil, fmt.Errorf("not enough data for int64")
	}
	return int64(binary.BigEndian.Uint64(data)), nil
}

func (t *Int64Type) SerializeJSON(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(value)
}

func (t *Int64Type) DeserializeJSON(data []byte) (interface{}, error) {
	if string(data) == "null" {
		return t.NullValue(), nil
	}
	var v int64
	err := json.Unmarshal(data, &v)
	return v, err
}

func (t *Int64Type) SerializeText(value interface{}) (string, error) {
	if value == nil {
		return "NULL", nil
	}
	return strconv.FormatInt(value.(int64), 10), nil
}

func (t *Int64Type) DeserializeText(text string) (interface{}, error) {
	if text == "NULL" {
		return t.NullValue(), nil
	}
	v, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return nil, err
	}
	return int64(v), nil
}

func (t *Int64Type) Compare(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}
	ai := a.(int64)
	bi := b.(int64)
	if ai < bi {
		return -1
	}
	if ai > bi {
		return 1
	}
	return 0
}

func (t *Int64Type) IsNull(value interface{}) bool {
	return value == nil
}

func (t *Int64Type) NullValue() interface{} {
	return nil
}

func (t *Int64Type) String() string {
	return "INT64"
}

func (t *Uint8Type) Metadata() *TypeMetadata {
	return &TypeMetadata{
		ID:        TypeIDUint8,
		Name:      "UINT8",
		Size:      1,
		Alignment: 1,
	}
}

func (t *Uint8Type) Size(value interface{}) int { return 1 }

func (t *Uint8Type) Validate(value interface{}) error {
	if value == nil {
		return nil
	}
	_, ok := value.(uint8)
	if !ok {
		return NewValidationError("value is not uint8")
	}
	// No need to check range, uint8 covers all possible values
	return nil
}

func (t *Uint8Type) Serialize(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte{0xFF}, nil
	}
	v, ok := value.(uint8)
	if !ok {
		return nil, fmt.Errorf("value is not uint8")
	}
	return []byte{byte(v)}, nil
}

func (t *Uint8Type) Deserialize(data []byte) (interface{}, error) {
	if len(data) == 1 && data[0] == 0xFF {
		return t.NullValue(), nil
	}
	if len(data) < 1 {
		return nil, fmt.Errorf("not enough data for uint8")
	}
	return uint8(data[0]), nil
}

func (t *Uint8Type) SerializeJSON(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(value)
}

func (t *Uint8Type) DeserializeJSON(data []byte) (interface{}, error) {
	if string(data) == "null" {
		return t.NullValue(), nil
	}
	var v uint8
	err := json.Unmarshal(data, &v)
	return v, err
}

func (t *Uint8Type) SerializeText(value interface{}) (string, error) {
	if value == nil {
		return "NULL", nil
	}
	return strconv.FormatUint(uint64(value.(uint8)), 10), nil
}

func (t *Uint8Type) DeserializeText(text string) (interface{}, error) {
	if text == "NULL" {
		return t.NullValue(), nil
	}
	v, err := strconv.ParseUint(text, 10, 8)
	if err != nil {
		return nil, err
	}
	return uint8(v), nil
}

func (t *Uint8Type) Compare(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}
	ai := a.(uint8)
	bi := b.(uint8)
	if ai < bi {
		return -1
	}
	if ai > bi {
		return 1
	}
	return 0
}

func (t *Uint8Type) IsNull(value interface{}) bool { return value == nil }
func (t *Uint8Type) NullValue() interface{}        { return nil }
func (t *Uint8Type) String() string                { return "UINT8" }

func (t *Uint16Type) Metadata() *TypeMetadata {
	return &TypeMetadata{
		ID:        TypeIDUint16,
		Name:      "UINT16",
		Size:      2,
		Alignment: 2,
	}
}

func (t *Uint16Type) Size(value interface{}) int { return 2 }

func (t *Uint16Type) Validate(value interface{}) error {
	if value == nil {
		return nil
	}
	_, ok := value.(uint16)
	if !ok {
		return NewValidationError("value is not uint16")
	}
	// No need to check range, uint16 covers all possible values
	return nil
}

func (t *Uint16Type) Serialize(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte{0xFF, 0xFF}, nil
	}
	v, ok := value.(uint16)
	if !ok {
		return nil, fmt.Errorf("value is not uint16")
	}
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, v)
	return buf, nil
}

func (t *Uint16Type) Deserialize(data []byte) (interface{}, error) {
	if len(data) == 2 && data[0] == 0xFF && data[1] == 0xFF {
		return t.NullValue(), nil
	}
	if len(data) < 2 {
		return nil, fmt.Errorf("not enough data for uint16")
	}
	return binary.BigEndian.Uint16(data), nil
}

func (t *Uint16Type) SerializeJSON(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(value)
}

func (t *Uint16Type) DeserializeJSON(data []byte) (interface{}, error) {
	if string(data) == "null" {
		return t.NullValue(), nil
	}
	var v uint16
	err := json.Unmarshal(data, &v)
	return v, err
}

func (t *Uint16Type) SerializeText(value interface{}) (string, error) {
	if value == nil {
		return "NULL", nil
	}
	return strconv.FormatUint(uint64(value.(uint16)), 10), nil
}

func (t *Uint16Type) DeserializeText(text string) (interface{}, error) {
	if text == "NULL" {
		return t.NullValue(), nil
	}
	v, err := strconv.ParseUint(text, 10, 16)
	if err != nil {
		return nil, err
	}
	return uint16(v), nil
}

func (t *Uint16Type) Compare(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}
	ai := a.(uint16)
	bi := b.(uint16)
	if ai < bi {
		return -1
	}
	if ai > bi {
		return 1
	}
	return 0
}

func (t *Uint16Type) IsNull(value interface{}) bool { return value == nil }
func (t *Uint16Type) NullValue() interface{}        { return nil }
func (t *Uint16Type) String() string                { return "UINT16" }

func (t *Uint32Type) Metadata() *TypeMetadata {
	return &TypeMetadata{
		ID:        TypeIDUint32,
		Name:      "UINT32",
		Size:      4,
		Alignment: 4,
	}
}

func (t *Uint32Type) Size(value interface{}) int { return 4 }

func (t *Uint32Type) Validate(value interface{}) error {
	if value == nil {
		return nil
	}
	_, ok := value.(uint32)
	if !ok {
		return NewValidationError("value is not uint32")
	}
	// No need to check range, uint32 covers all possible values
	return nil
}

func (t *Uint32Type) Serialize(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte{0xFF, 0xFF, 0xFF, 0xFF}, nil
	}
	v, ok := value.(uint32)
	if !ok {
		return nil, fmt.Errorf("value is not uint32")
	}
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	return buf, nil
}

func (t *Uint32Type) Deserialize(data []byte) (interface{}, error) {
	if len(data) == 4 && data[0] == 0xFF && data[1] == 0xFF && data[2] == 0xFF && data[3] == 0xFF {
		return t.NullValue(), nil
	}
	if len(data) < 4 {
		return nil, fmt.Errorf("not enough data for uint32")
	}
	return binary.BigEndian.Uint32(data), nil
}

func (t *Uint32Type) SerializeJSON(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(value)
}

func (t *Uint32Type) DeserializeJSON(data []byte) (interface{}, error) {
	if string(data) == "null" {
		return t.NullValue(), nil
	}
	var v uint32
	err := json.Unmarshal(data, &v)
	return v, err
}

func (t *Uint32Type) SerializeText(value interface{}) (string, error) {
	if value == nil {
		return "NULL", nil
	}
	return strconv.FormatUint(uint64(value.(uint32)), 10), nil
}

func (t *Uint32Type) DeserializeText(text string) (interface{}, error) {
	if text == "NULL" {
		return t.NullValue(), nil
	}
	v, err := strconv.ParseUint(text, 10, 32)
	if err != nil {
		return nil, err
	}
	return uint32(v), nil
}

func (t *Uint32Type) Compare(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}
	ai := a.(uint32)
	bi := b.(uint32)
	if ai < bi {
		return -1
	}
	if ai > bi {
		return 1
	}
	return 0
}

func (t *Uint32Type) IsNull(value interface{}) bool { return value == nil }
func (t *Uint32Type) NullValue() interface{}        { return nil }
func (t *Uint32Type) String() string                { return "UINT32" }

func (t *Uint64Type) Metadata() *TypeMetadata {
	return &TypeMetadata{
		ID:        TypeIDUint64,
		Name:      "UINT64",
		Size:      8,
		Alignment: 8,
	}
}

func (t *Uint64Type) Size(value interface{}) int { return 8 }

func (t *Uint64Type) Validate(value interface{}) error {
	if value == nil {
		return nil
	}
	_, ok := value.(uint64)
	if !ok {
		return NewValidationError("value is not uint64")
	}
	// No need to check range, uint64 covers all possible values
	return nil
}

func (t *Uint64Type) Serialize(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, nil
	}
	v, ok := value.(uint64)
	if !ok {
		return nil, fmt.Errorf("value is not uint64")
	}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, v)
	return buf, nil
}

func (t *Uint64Type) Deserialize(data []byte) (interface{}, error) {
	if len(data) == 8 &&
		data[0] == 0xFF && data[1] == 0xFF && data[2] == 0xFF && data[3] == 0xFF &&
		data[4] == 0xFF && data[5] == 0xFF && data[6] == 0xFF && data[7] == 0xFF {
		return t.NullValue(), nil
	}
	if len(data) < 8 {
		return nil, fmt.Errorf("not enough data for uint64")
	}
	return binary.BigEndian.Uint64(data), nil
}

func (t *Uint64Type) SerializeJSON(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(value)
}

func (t *Uint64Type) DeserializeJSON(data []byte) (interface{}, error) {
	if string(data) == "null" {
		return t.NullValue(), nil
	}
	var v uint64
	err := json.Unmarshal(data, &v)
	return v, err
}

func (t *Uint64Type) SerializeText(value interface{}) (string, error) {
	if value == nil {
		return "NULL", nil
	}
	return strconv.FormatUint(value.(uint64), 10), nil
}

func (t *Uint64Type) DeserializeText(text string) (interface{}, error) {
	if text == "NULL" {
		return t.NullValue(), nil
	}
	v, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		return nil, err
	}
	return uint64(v), nil
}

func (t *Uint64Type) Compare(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}
	ai := a.(uint64)
	bi := b.(uint64)
	if ai < bi {
		return -1
	}
	if ai > bi {
		return 1
	}
	return 0
}

func (t *Uint64Type) IsNull(value interface{}) bool { return value == nil }
func (t *Uint64Type) NullValue() interface{}        { return nil }
func (t *Uint64Type) String() string                { return "UINT64" }
