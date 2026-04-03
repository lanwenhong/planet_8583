package utils

func PtrString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func PtrInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}
