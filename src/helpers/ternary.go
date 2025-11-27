package helpers

func Ternary(condition bool, match any, fallback any) any {
	if condition {
		return match
	} else {
		return fallback
	}
}
