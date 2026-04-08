package codec

import "testing"

func TestEncodeLen(t *testing.T) {
	id := uint64(123456789)
	code := Encode(id)
	if len(code) != codeLen {
		t.Errorf("expected code length %d, got %d", codeLen, len(code))
	}
}

func TestEncodeUnique(t *testing.T) {
	ids := []uint64{0, 1, 23, 67, 127, 123456789, 9876543210}
	codes := make(map[string]struct{})

	for _, id := range ids {
		code := Encode(id)
		if _, exists := codes[code]; exists {
			t.Errorf("duplicate code generated for ID %d: %s", id, code)
		}
		codes[code] = struct{}{}
	}
}

func TestEncodeConsistency(t *testing.T) {
	id := uint64(123456789)
	code1 := Encode(id)
	code2 := Encode(id)

	if code1 != code2 {
		t.Errorf("unconsistent for the same ID: %s vs %s", code1, code2)
	}
}
