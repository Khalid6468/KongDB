package types

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const MAX_TEXT_LENGTH = 65535

type TextType struct {
	collation string
}

func (t *TextType) Metadata() *TypeMetadata {
	return &TypeMetadata{
		ID:        TypeIDText,
		Name:      "TEXT",
		Size:      -1,
		Alignment: 1,
	}
}

func NewTextType(collation string) *TextType {
	if collation == "" {
		collation = "utf8_general_ci"
	}
	return &TextType{
		collation: collation,
	}
}

type TextValue struct {
	Value  string
	IsNull bool
}

func (t *TextType) Collation() string {
	return t.collation
}

func (t *TextType) Validate(value interface{}) error {
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

	if len(str) > MAX_TEXT_LENGTH {
		return NewValidationError("text too large (max 65535 bytes)")
	}
	return nil
}

func (t *TextType) IsNull(value interface{}) bool {
	return value == nil
}

func (t *TextType) NullValue() interface{} {
	return nil
}

func (t *TextType) Size(value interface{}) int {
	if t.IsNull(value) {
		return 1
	}

	str, ok := value.(string)
	if !ok {
		return 0
	}
	// Todo:  Variable length encoding
	return 4 + len(str)
}

func (t *TextType) String() string {
	return "TEXT"
}

func (t *TextType) Serialize(value interface{}) ([]byte, error) {
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

func (t *TextType) Deserialize(data []byte) (interface{}, error) {
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

func (t *TextType) SerializeJSON(value interface{}) ([]byte, error) {
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

func (t *TextType) DeserializeJSON(data []byte) (interface{}, error) {
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

func (t *TextType) SerializeText(value interface{}) (string, error) {
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

func (t *TextType) DeserializeText(text string) (interface{}, error) {
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

func (t *TextType) unicodeCompare(a, b string) int {
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

func (t *TextType) Compare(a, b interface{}) int {
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

func (t *TextType) compareWithCollation(a, b string) int {
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

func (t *TextType) SetCollation(collation string) {
	t.collation = collation
}

func (t *TextType) GetCollation() string {
	return t.collation
}

func (t *TextType) IsCaseSensitive() bool {
	return strings.Contains(t.collation, "_bin") || t.collation == "binary"
}
