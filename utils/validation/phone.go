package validation

import (
	"fmt"
	"strings"
)

func IsPhoneNumber(input string) bool {
	for _, r := range input {
		switch {
		case r >= '0' && r <= '9', r == ' ', r == '-', r == '(', r == ')', r == '+':
			continue
		default:
			return false
		}
	}

	return true
}

func ValidateAndNormalizePhone(input string) (string, error) {
	if !strings.HasPrefix(input, "+7") && !strings.HasPrefix(input, "7") && !strings.HasPrefix(input, "8") {
		return "", fmt.Errorf("номер должен начинаться с +7, 7 ли или 8")
	}

	var builder strings.Builder

	for _, r := range input {
		switch {
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
		case r == '+':
			if builder.Len() > 0 {
				return "", fmt.Errorf("символ '+' должен быть только в начале номера")
			}
			continue
		case r == ' ', r == '-', r == '(', r == ')':
			continue
		default:
			return "", fmt.Errorf("номер содержит недопустимые символы")
		}
	}

	phone := builder.String()
	if len(phone) != 11 {
		return "", fmt.Errorf("номер должен содержать 11 цифр")
	}

	if strings.HasPrefix(phone, "8") {
		phone = "7" + phone[1:]
	}

	return phone, nil
}
