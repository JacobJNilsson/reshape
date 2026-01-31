package core

// WarningMessage returns the default message for a warning code.
func WarningMessage(code WarningCode) string {
	switch code {
	case WarningCodeJoinArray:
		return "joined array into string"
	case WarningCodeDropField:
		return "dropped field"
	case WarningCodeCoerceType:
		return "coerced type"
	default:
		return "warning"
	}
}

func WarningFor(code WarningCode, path string) Warning {
	return Warning{
		Code:    code,
		Path:    path,
		Message: WarningMessage(code),
	}
}
