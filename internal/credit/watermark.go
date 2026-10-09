package credit

import "strings"

// Zero-width characters
const (
	zwsp = "\u200B" // Zero-width space
	zwnj = "\u200C" // Zero-width non-joiner
	zwj  = "\u200D" // Zero-width joiner
	bom  = "\uFEFF" // Byte order mark
)

// EncodeToZeroWidth converts text to zero-width Unicode characters
func EncodeToZeroWidth(text string) string {
	var encoded strings.Builder
	for _, char := range []byte(text) {
		for i := 7; i >= 0; i-- {
			bit := (char >> i) & 1
			if bit == 1 {
				encoded.WriteString(zwsp)
			} else {
				encoded.WriteString(zwnj)
			}
		}
		encoded.WriteString(zwj) // Character separator
	}
	return encoded.String()
}

// DecodeFromZeroWidth extracts hidden text
func DecodeFromZeroWidth(watermarked string) string {
	var decoded []byte
	var currentChar byte
	var bitCount int

	chars := []rune(watermarked)
	for _, char := range chars {
		switch string(char) {
		case zwsp:
			currentChar = (currentChar << 1) | 1
			bitCount++
		case zwnj:
			currentChar = (currentChar << 1) | 0
			bitCount++
		case zwj:
			if bitCount == 8 {
				decoded = append(decoded, currentChar)
			}
			currentChar = 0
			bitCount = 0
		}
	}
	return string(decoded)
}

// WatermarkMessage adds invisible "STDBOTS|deepanshu.in" watermark to any text
func WatermarkMessage(text ...string) string {
	base := ""
	if len(text) > 0 {
		base = text[0]
	}
	watermark := EncodeToZeroWidth("STDBOTS|deepanshu.in")
	return base + watermark
}

// HasWatermark checks if text has STD BOTS watermark
func HasWatermark(text string) bool {
	extracted := DecodeFromZeroWidth(text)
	return strings.Contains(extracted, "STDBOTS|deepanshu.in")
}

// StripWatermark removes watermark (for internal use)
func StripWatermark(text string) string {
	result := strings.ReplaceAll(text, zwsp, "")
	result = strings.ReplaceAll(result, zwnj, "")
	result = strings.ReplaceAll(result, zwj, "")
	result = strings.ReplaceAll(result, bom, "")
	return result
}
