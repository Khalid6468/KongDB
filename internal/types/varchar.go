package types

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const MAX_VARCHAR_LENGTH int = 255

type VarcharType struct {
	maxLength int
	collation string
}

func NewVarcharType(maxLength int, collation string) *VarcharType {
	if maxLength <= 0 {
		maxLength = MAX_VARCHAR_LENGTH
	}
	if collation == "" {
		collation = "utf8_general_ci"
	}
	return &VarcharType{
		maxLength: maxLength,
		collation: collation,
	}
}

func (t *VarcharType) Metadata() *TypeMetadata {
	return &TypeMetadata{
		ID:        TypeIDVarchar,
		Name:      fmt.Sprintf("VARCHAR(%d)", t.maxLength),
		Size:      -1,
		Alignment: 1,
	}
}

type VarcharValue struct {
	Value  string
	Length int
	IsNull bool
}

func (t *VarcharType) MaxLength() int {
	return t.maxLength
}

func (t *VarcharType) Collation() string {
	return t.collation
}

func (t *VarcharType) Validate(value interface{}) error {
	if value == nil {
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return NewValidationError(fmt.Sprintf("expected string, got %T", value))
	}

	if !utf8.ValidString(str) {
		return NewValidationError("invalid UTF-8 sequence")
	}

	runeCount := utf8.RuneCountInString(str)
	if runeCount > t.maxLength {
		return NewValidationError(fmt.Sprintf("string length %d exceeds maximum %d", runeCount, t.maxLength))
	}

	return nil
}

func (t *VarcharType) IsNull(value interface{}) bool {
	return value == nil
}

func (t *VarcharType) NullValue() interface{} {
	return nil
}

func (t *VarcharType) Size(value interface{}) int {
	if t.IsNull(value) {
		return 1
	}

	str, ok := value.(string)
	if !ok {
		return 0
	}
	// Todo: Variable length encoding
	return 4 + len(str)
}

func (t *VarcharType) String() string {
	return fmt.Sprintf("VARCHAR(%d)", t.maxLength)
}

func (t *VarcharType) Serialize(value interface{}) ([]byte, error) {
	if t.IsNull(value) {
		return []byte{0xFF}, nil
	}

	str, ok := value.(string)
	if !ok {
		return nil, NewValidationError(fmt.Sprintf("expected string, got %T", value))
	}

	if err := t.Validate(value); err != nil {
		return nil, err
	}

	data := make([]byte, 4+len(str))

	binary.BigEndian.PutUint32(data[:4], uint32(len(str)))
	copy(data[4:], str)

	return data, nil
}

func (t *VarcharType) Deserialize(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, NewValidationError("empty data")
	}

	if len(data) == 1 && data[0] == 0xFF {
		return nil, nil
	}

	if len(data) < 4 {
		return nil, NewValidationError("insufficient data for length prefix")
	}

	length := binary.BigEndian.Uint32(data[:4])

	if len(data) < int(4+length) {
		return nil, NewValidationError("data shorter than declared length")
	}

	str := string(data[4 : 4+length])
	if err := t.Validate(str); err != nil {
		return nil, err
	}

	return str, nil
}

func (t *VarcharType) SerializeJSON(value interface{}) ([]byte, error) {
	if t.IsNull(value) {
		return []byte("null"), nil
	}

	str, ok := value.(string)
	if !ok {
		return nil, NewValidationError(fmt.Sprintf("expected string, got %T", value))
	}

	if err := t.Validate(value); err != nil {
		return nil, err
	}

	return json.Marshal(str)
}

func (t *VarcharType) DeserializeJSON(data []byte) (interface{}, error) {
	var str *string

	if err := json.Unmarshal(data, &str); err != nil {
		return nil, NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
	}

	if str == nil {
		return nil, nil
	}

	if err := t.Validate(*str); err != nil {
		return nil, err
	}

	return *str, nil
}

func (t *VarcharType) SerializeText(value interface{}) (string, error) {
	if t.IsNull(value) {
		return "NULL", nil
	}

	str, ok := value.(string)
	if !ok {
		return "", NewValidationError(fmt.Sprintf("expected string, got %T", value))
	}

	if err := t.Validate(value); err != nil {
		return "", err
	}

	escaped := strings.ReplaceAll(str, "'", "''")

	return "'" + escaped + "'", nil
}

func (t *VarcharType) DeserializeText(text string) (interface{}, error) {
	text = strings.TrimSpace(text)

	if strings.ToUpper(text) == "NULL" {
		return nil, nil
	}

	if len(text) < 2 || !strings.HasPrefix(text, "'") || !strings.HasSuffix(text, "'") {
		return nil, NewValidationError("expected quoted string")
	}

	str := text[1 : len(text)-1]
	str = strings.ReplaceAll(str, "''", "'") // Unescape single quotes

	if err := t.Validate(str); err != nil {
		return nil, err
	}

	return str, nil
}

func (t *VarcharType) unicodeCompare(a, b string) int {
	runesA := []rune(a)
	runesB := []rune(b)

	minLen := len(runesA)
	if len(runesB) < minLen {
		minLen = len(runesB)
	}

	for i := 0; i < minLen; i++ {
		if runesA[i] != runesB[i] {
			normA := unicode.ToLower(runesA[i])
			normB := unicode.ToLower(runesB[i])

			if normA < normB {
				return -1
			}
			if normA > normB {
				return 1
			}
		}
	}

	if len(runesA) < len(runesB) {
		return -1
	}
	if len(runesA) > len(runesB) {
		return 1
	}

	return 0
}

func (t *VarcharType) Compare(a, b interface{}) int {
	aNull := t.IsNull(a)
	bNull := t.IsNull(b)

	if aNull && bNull {
		return 0
	}
	if aNull {
		return -1
	}
	if bNull {
		return 1
	}

	strA, okA := a.(string)
	strB, okB := b.(string)

	if !okA || !okB {
		strA = fmt.Sprintf("%v", a)
		strB = fmt.Sprintf("%v", b)
	}

	return t.compareWithCollation(strA, strB)
}

func (t *VarcharType) compareWithCollation(a, b string) int {
	switch t.collation {
	case "utf8_bin", "binary":
		return strings.Compare(a, b)

	case "utf8_general_ci", "utf8_unicode_ci":
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))

	case "utf8_unicode_520_ci":
		return t.unicodeCompare(a, b)

	default:
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))
	}
}

func (t *VarcharType) SetCollation(collation string) {
	t.collation = collation
}

func (t *VarcharType) GetCollation() string {
	return t.collation
}

func (t *VarcharType) IsCaseSensitive() bool {
	return strings.Contains(t.collation, "_bin") || t.collation == "binary"
}
