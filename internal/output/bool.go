package output

// FormatBool renders a boolean for human-readable output as "yes" / "no".
// Machine output (--json/--ndjson) always emits the raw true/false instead.
func FormatBool(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// FormatBoolPtr renders a nullable boolean as "yes" / "no", and a JSON null (or
// an absent field) as the placeholder — an unknown is not a false.
func FormatBoolPtr(b *bool) string {
	if b == nil {
		return placeholder
	}
	return FormatBool(*b)
}
