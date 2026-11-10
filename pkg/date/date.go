package date

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не может быть пустым")
	}

	date, err := time.Parse(DateFormat, dstart)
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

func resetTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func nextDateByDays(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("неверный формат правила d")
	}
	interval, err := strconv.Atoi(parts[1])
	if err != nil || interval <= 0 || interval > 400 {
		return "", fmt.Errorf("интервал в днях должен быть от 1 до 400")
	}

	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(DateFormat), nil
}

func nextDateByYear(now, date time.Time) (string, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(DateFormat), nil
}

func nextDateByWeekdays(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("неверный формат правила w")
	}

	dayStrs := strings.Split(parts[1], ",")
	weekdays := make(map[int]bool)

	for _, dayStr := range dayStrs {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < 1 || day > 7 {
			return "", fmt.Errorf("день недели должен быть от 1 до 7")
		}
		weekdays[day%7] = true
	}

	currentDate := date
	for i := 0; i < 365*5; i++ { // лимит 5 лет
		// Дата должна быть >= dstart и строго > now
		if !currentDate.Before(date) && afterNow(currentDate, now) && weekdays[int(currentDate.Weekday())] {
			return currentDate.Format(DateFormat), nil
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}
	return "", fmt.Errorf("не удалось найти подходящий день недели")
}

func nextDateByMonthDays(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("неверный формат правила m")
	}

	dayStrs := strings.Split(parts[1], ",")
	monthDays := make([]int, 0, len(dayStrs))

	for _, dayStr := range dayStrs {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day == 0 || day > 31 || day < -2 {
			return "", fmt.Errorf("день месяца должен быть от 1 до 31 или -1, -2")
		}
		monthDays = append(monthDays, day)
	}

	months := make(map[int]bool)
	if len(parts) == 3 {
		monthStrs := strings.Split(parts[2], ",")
		for _, monthStr := range monthStrs {
			month, err := strconv.Atoi(monthStr)
			if err != nil || month < 1 || month > 12 {
				return "", fmt.Errorf("месяц должен быть от 1 до 12")
			}
			months[month] = true
		}
	}

	currentDate := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())

	for i := 0; i < 120; i++ {
		if len(months) == 0 || months[int(currentDate.Month())] {
			var candidates []time.Time
			for _, day := range monthDays {
				targetDay := day
				if day < 0 {
					lastDay := getLastDayOfMonth(currentDate.Year(), currentDate.Month())
					targetDay = lastDay + day + 1
				}
				if targetDay >= 1 && targetDay <= getLastDayOfMonth(currentDate.Year(), currentDate.Month()) {
					t := time.Date(currentDate.Year(), currentDate.Month(), targetDay, 0, 0, 0, 0, date.Location())

					// Дата не может быть раньше dstart и должна быть > now
					if !t.Before(date) && afterNow(t, now) {
						candidates = append(candidates, t)
					}
				}
			}

			if len(candidates) > 0 {
				// Берем минимальную подходящую дату в этом месяце
				minDate := candidates[0]
				for _, c := range candidates {
					if c.Before(minDate) {
						minDate = c
					}
				}
				return minDate.Format(DateFormat), nil
			}
		}
		currentDate = currentDate.AddDate(0, 1, 0)
	}

	return "", fmt.Errorf("не удалось найти подходящую дату")
}

func getLastDayOfMonth(year int, month time.Month) int {
	firstDayNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDayNextMonth.AddDate(0, 0, -1)
	return lastDay.Day()
}
