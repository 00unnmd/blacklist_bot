package validation

import (
	"fmt"
	"unicode/utf8"
)

func ValidateDescriptionStr(input string) error {
	count := utf8.RuneCountInString(input)

	if count < 7 {
		return fmt.Errorf("поле должно содержать минимум 7 символов")
	}

	return nil
}
