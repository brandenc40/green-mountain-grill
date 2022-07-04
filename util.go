package gmg

func tempWithHighValue(low, high uint8) int {
	return ((int(high) & 255) << 8) + (int(low) & 255)
}

func fourByteConversion(val1, val2, val3, val4 uint8) int {
	return ((int(val4) & 255) << 24) + ((int(val3) & 255) << 16) + ((int(val2) & 255) << 8) + (int(val1) & 255)
}

func curveValue(fourBytesConverted int) [3]int {
	var (
		val1 = 0
		val2 = 0
		val3 = fourBytesConverted
	)
	for {
		i := val3
		if i < 60 {
			break
		}
		val3 = i - 60
		val2++
	}
	for {
		i2 := val3
		if i2 >= 0 {
			break
		}
		val3 = i2 + 60
		val2--
	}
	for {
		i3 := val2
		if i3 < 60 {
			break
		}
		val2 = i3 - 60
		val1++
	}
	for {
		i4 := val2
		if i4 < 0 {
			val2 = i4 + 60
			val1--
		} else {
			return [3]int{val1, val2, val3}
		}
	}
}
