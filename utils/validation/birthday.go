package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ValidateBirthdayStr(input string) error {
	matched, _ := regexp.MatchString(`^\d{2}\.\d{2}\.\d{4}$`, input)
	if !matched {
		return fmt.Errorf("неверный формат даты. Используйте ДД.ММ.ГГГГ")
	}

	parts := strings.Split(input, ".")
	if len(parts) != 3 {
		return fmt.Errorf("неверный формат даты. Используйте ДД.ММ.ГГГГ")
	}

	day, err1 := strconv.Atoi(parts[0])
	month, err2 := strconv.Atoi(parts[1])
	year, err3 := strconv.Atoi(parts[2])

	if err1 != nil || err2 != nil || err3 != nil {
		return fmt.Errorf("дата должна содержать только цифры")
	}
	if year < 1900 || year > time.Now().Year() {
		return fmt.Errorf("некорректный год рождения")
	}
	if month < 1 || month > 12 {
		return fmt.Errorf("некорректный месяц")
	}

	maxDays := 31
	switch month {
	case 4, 6, 9, 11:
		maxDays = 30
	case 2:
		if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
			maxDays = 29
		} else {
			maxDays = 28
		}
	}

	if day < 1 || day > maxDays {
		return fmt.Errorf("некорректный день для указанного месяца и года")
	}

	birthDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if birthDate.After(time.Now()) {
		return fmt.Errorf("дата рождения не может быть в будущем")
	}

	return nil
}
