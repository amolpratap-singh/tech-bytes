package iteration

func Repeat(character string, times ...int) string {
	var repeated string
	noRepeation := 3

	if len(times) > 0 {
		noRepeation = times[0]
	}

	for i := 0; i < noRepeation; i++ {
		repeated += character
	}
	return repeated
}

func CharacterCount(str, char string) (int) {
	count := 0

	if len(char) == 0 {
		return 0
	}
	
	for i := 0; i < len(str); i++ {
		if str[i] == []byte(char)[0]  {
			count += 1
		}
	}

	return count
}