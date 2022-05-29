package gmg

func getTempWithHighValue(low, high uint8) int {
	// a high value of 2 represents that the temp is not available
	if high == 2 {
		return 0
	}
	return int(low) + int(high)*256
}

func tempToLowHighValue(temp int) (low, high uint8) {
	return 0, 0
}
