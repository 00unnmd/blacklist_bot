package validation

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

func ValidateCityStr(input string) error {
	if utf8.RuneCountInString(input) < 3 {
		return fmt.Errorf("название города должно содержать минимум 3 символа")
	}

	hasDigits, _ := regexp.MatchString(`\d`, input)
	if hasDigits {
		return fmt.Errorf("название города не должно содержать цифр")
	}

	hasSpecialChars, _ := regexp.MatchString(`[!@#$%^&*()_+=\[\]{};':"\\|,.<>/?]`, input)
	if hasSpecialChars {
		return fmt.Errorf("название города содержит недопустимые символы")
	}

	return nil
}
