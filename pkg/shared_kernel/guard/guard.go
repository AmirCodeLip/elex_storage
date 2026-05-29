package guard

import "time"

/* This method will check strings are null and empty or not if one of them have value it will skip  */
func AgainstPNullStrs(params ...*string) bool {
	oneValid := false
	for _, param := range params {
		if param != nil && *param != "" {
			oneValid = true
		}
	}
	return oneValid
}

/* This method will check string is null and empty or not if is null it's return false  */
func AgainstPNullStr(param *string) bool {
	return !(param == nil || *param == "")
}

func AgainstEmptyStr(param string) bool {
	return !(param == "")
}

func AgainstTimeDurationPtr(param *time.Duration) bool {
	return !(param == nil || *param == 0)
}
