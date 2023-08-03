package utils

func IsChanClosable(ch chan struct{}) bool {
	if ch == nil {
		return false
	}
	select {
	case _, ok := <-ch:
		if ok {
			return true
		}
		return false
	default:
		return true
	}
}
