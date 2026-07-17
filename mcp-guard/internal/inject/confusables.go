package inject

// confusableMap maps visually-confusable non-ASCII characters to their ASCII equivalents.
// These are characters from Cyrillic, Greek, and other scripts that appear identical
// to ASCII letters in most fonts, enabling homoglyph injection attacks.
var confusableMap = map[rune]rune{
	// Cyrillic
	0x0430: 'a', // а
	0x0435: 'e', // е
	0x043E: 'o', // о
	0x0440: 'p', // р
	0x0441: 'c', // с
	0x0445: 'x', // х
	0x0443: 'y', // у
	0x0432: 'b', // в
	0x043A: 'k', // к
	0x043C: 'm', // м
	0x043D: 'h', // н
	0x0442: 't', // т
	0x0438: 'i', // и (looks like Latin i in some fonts)

	// Greek
	0x03BF: 'o', // ο (omicron)
	0x03B5: 'e', // ε (epsilon)
	0x03C1: 'p', // ρ (rho)
	0x03C4: 't', // τ (tau)
	0x03B1: 'a', // α (alpha)
	0x03B7: 'n', // η (eta)
	0x03BA: 'k', // κ (kappa)
	0x03BD: 'v', // ν (nu)
	0x03C3: 'c', // σ (sigma)
	0x03C7: 'x', // χ (chi)
	0x03B9: 'i', // ι (iota)
	0x03BC: 'm', // μ (mu)

	// Latin extended
	0x00E1: 'a', // á
	0x00E9: 'e', // é
	0x00ED: 'i', // í
	0x00F3: 'o', // ó
	0x00FA: 'u', // ú

	// Full-width ASCII (used in CJK environments)
	0xFF21: 'A', // Ａ
	0xFF22: 'B', // Ｂ
	0xFF23: 'C', // Ｃ
	0xFF24: 'D', // Ｄ
	0xFF25: 'E', // Ｅ
	0xFF26: 'F', // Ｆ
	0xFF27: 'G', // Ｇ
	0xFF28: 'H', // Ｈ
	0xFF29: 'I', // Ｉ
	0xFF2A: 'J', // Ｊ
	0xFF2B: 'K', // Ｋ
	0xFF2C: 'L', // Ｌ
	0xFF2D: 'M', // Ｍ
	0xFF2E: 'N', // Ｎ
	0xFF2F: 'O', // Ｏ
	0xFF30: 'P', // Ｐ
	0xFF31: 'Q', // Ｑ
	0xFF32: 'R', // Ｒ
	0xFF33: 'S', // Ｓ
	0xFF34: 'T', // Ｔ
	0xFF35: 'U', // Ｕ
	0xFF36: 'V', // Ｖ
	0xFF37: 'W', // Ｗ
	0xFF38: 'X', // Ｘ
	0xFF39: 'Y', // Ｙ
	0xFF3A: 'Z', // Ｚ
	0xFF41: 'a', // ａ
	0xFF42: 'b', // ｂ
	0xFF43: 'c', // ｃ
	0xFF44: 'd', // ｄ
	0xFF45: 'e', // ｅ
	0xFF46: 'f', // ｆ
	0xFF47: 'g', // ｇ
	0xFF48: 'h', // ｈ
	0xFF49: 'i', // ｉ
	0xFF4A: 'j', // ｊ
	0xFF4B: 'k', // ｋ
	0xFF4C: 'l', // ｌ
	0xFF4D: 'm', // ｍ
	0xFF4E: 'n', // ｎ
	0xFF4F: 'o', // ｏ
	0xFF50: 'p', // ｐ
	0xFF51: 'q', // ｑ
	0xFF52: 'r', // ｒ
	0xFF53: 's', // ｓ
	0xFF54: 't', // ｔ
	0xFF55: 'u', // ｕ
	0xFF56: 'v', // ｖ
	0xFF57: 'w', // ｗ
	0xFF58: 'x', // ｘ
	0xFF59: 'y', // ｙ
	0xFF5A: 'z', // ｚ
}

// isConfusableRune checks if a rune is visually confusable with ASCII.
// Returns the ASCII equivalent and true if confusable.
func isConfusableRune(r rune) (ascii rune, isConfusable bool) {
	ascii, ok := confusableMap[r]
	return ascii, ok
}
