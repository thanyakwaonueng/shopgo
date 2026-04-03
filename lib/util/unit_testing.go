package util

var isTestMode = false

func IsTestMode() bool {
	return isTestMode
}

func TurnOnTestMode() {
	isTestMode = true
}

func TurnOffTestMode() {
	isTestMode = false
}
