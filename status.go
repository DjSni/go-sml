package sml

// StatusParse preserves all 64 bits of the SML_Status unsigned choice.
func StatusParse(buf *Buffer) (uint64, error) {
	if BufOptionalIsSkipped(buf) {
		return 0, nil
	}
	return U64Parse(buf)
}
