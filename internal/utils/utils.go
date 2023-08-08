package utils

func IsChanClosable(ch chan struct{}) bool {
	if ch == nil {
		return false
	}
	select {
	case _, ok := <-ch:
		return ok
	default:
		return true
	}
}
