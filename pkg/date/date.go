package date

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату задачи на основе правила повторения.
// now - время, от которого ищется ближайшая дата.
// dstart - исходное время в формате 20060102, от которого начинается отсчёт повторений.
// repeat - правило повторения.
// Возвращает следующую дату в формате 20060102 и ошибку.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не может быть пустым")
	}

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты dstart: %w", err)
	}

	now = resetTime(now)
	date = resetTime(date)

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("неверный формат правила: %s", repeat)
	}

	ruleType := parts[0]

	switch ruleType {
	case "d":
		return nextDateByDays(now, date, parts)
	case "y":
		return nextDateByYear(now, date)
	case "w":
		return nextDateByWeekdays(now, date, parts)
	case "m":
		return nextDateByMonthDays(now, date, parts)
	default:
		return "", fmt.Errorf("неподдерживаемый тип правила: %s", ruleType)
	}
}

// resetTime обнуляет время, оставляя только дату
func resetTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// afterNow проверяет, что date строго больше now
func afterNow(date, now time.Time) bool {
	return date.After(now)
}

// nextDateByDays обрабатывает правило "d <число>"
func nextDateByDays(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("неверный формат правила d: ожидается 'd <число>'")
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("неверное число в правиле d: %w", err)
	}

	if interval <= 0 || interval > 400 {
		return "", fmt.Errorf("интервал в днях должен быть от 1 до 400")
	}

	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format("20060102"), nil
}

// nextDateByYear обрабатывает правило "y"
func nextDateByYear(now, date time.Time) (string, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format("20060102"), nil
}

// nextDateByWeekdays обрабатывает правило "w <дни_недели>"
func nextDateByWeekdays(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("неверный формат правила w: ожидается 'w <дни>'")
	}

	dayStrs := strings.Split(parts[1], ",")
	weekdays := make(map[int]bool)

	for _, dayStr := range dayStrs {
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return "", fmt.Errorf("неверный день недели: %s", dayStr)
		}
		if day < 1 || day > 7 {
			return "", fmt.Errorf("день недели должен быть от 1 до 7")
		}
		// Конвертируем: в правиле 1=пн, 7=вс; в Go 0=вс, 1=пн, ..., 6=сб
		goWeekday := day % 7
		weekdays[goWeekday] = true
	}

	currentDate := now.AddDate(0, 0, 1)

	for i := 0; i < 7; i++ {
		if weekdays[int(currentDate.Weekday())] {
			return currentDate.Format("20060102"), nil
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return "", fmt.Errorf("не удалось найти подходящий день недели")
}

// nextDateByMonthDays обрабатывает правило "m <дни_месяца> [месяцы]"
func nextDateByMonthDays(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("неверный формат правила m: ожидается 'm <дни> [месяцы]'")
	}

	dayStrs := strings.Split(parts[1], ",")
	monthDays := make([]int, 0, len(dayStrs))

	for _, dayStr := range dayStrs {
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return "", fmt.Errorf("неверный день месяца: %s", dayStr)
		}
		if day == 0 || day > 31 || day < -2 {
			return "", fmt.Errorf("день месяца должен быть от 1 до 31 или -1, -2")
		}
		monthDays = append(monthDays, day)
	}

	months := make(map[int]bool)
	if len(parts) == 3 {
		monthStrs := strings.Split(parts[2], ",")
		for _, monthStr := range monthStrs {
			month, err := strconv.Atoi(monthStr)
			if err != nil {
				return "", fmt.Errorf("неверный месяц: %s", monthStr)
			}
			if month < 1 || month > 12 {
				return "", fmt.Errorf("месяц должен быть от 1 до 12")
			}
			months[month] = true
		}
	}

	// Начинаем от месяца dstart, а не от now
	currentDate := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, now.Location())

	for i := 0; i < 48; i++ {
		if len(months) == 0 || months[int(currentDate.Month())] {
			candidateDates := make([]time.Time, 0)

			for _, day := range monthDays {
				targetDay := day
				if day < 0 {
					lastDay := getLastDayOfMonth(currentDate.Year(), currentDate.Month())
					targetDay = lastDay + day + 1
				}

				if targetDay >= 1 && targetDay <= getLastDayOfMonth(currentDate.Year(), currentDate.Month()) {
					targetDate := time.Date(currentDate.Year(), currentDate.Month(), targetDay, 0, 0, 0, 0, now.Location())
					candidateDates = append(candidateDates, targetDate)
				}
			}

			// Находим минимальную дату, которая > now
			var minValidDate time.Time
			found := false

			for _, candidate := range candidateDates {
				if afterNow(candidate, now) {
					if !found || candidate.Before(minValidDate) {
						minValidDate = candidate
						found = true
					}
				}
			}

			if found {
				return minValidDate.Format("20060102"), nil
			}
		}

		currentDate = currentDate.AddDate(0, 1, 0)
	}

	return "", fmt.Errorf("не удалось найти подходящую дату в правиле m")
}

// getLastDayOfMonth возвращает последний день месяца
func getLastDayOfMonth(year int, month time.Month) int {
	firstDayNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDayNextMonth.AddDate(0, 0, -1)
	return lastDay.Day()
}
