package service

import (
	"go_todo_list/errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, startDate time.Time, repeat string) (string, error) {
	switch {
	case strings.HasPrefix(repeat, "d "):
		days, err := strconv.Atoi(strings.TrimSpace(repeat[2:]))
		if err != nil || days < 1 || days > 400 {
			return "", &errors.ValidationError{Message: "invalid repeat interval for days"}
		}
		return calcNextDateForDays(now, startDate, days), nil
	case repeat == "y":
		return calcNextDateForYearly(now, startDate), nil
	case strings.HasPrefix(repeat, "w "):
		daysOfWeek, err := parseDaysOfWeek(repeat[2:])
		if err != nil {
			return "", err
		}
		return calcNextDateForWeekdays(now, startDate, daysOfWeek), nil
	case strings.HasPrefix(repeat, "m "):
		daysOfMonth, months, err := parseDaysAndMonths(repeat[2:])
		if err != nil {
			return "", err
		}
		return calcNextDateForMonthly(now, startDate, daysOfMonth, months), nil
	default:
		return "", &errors.ValidationError{Message: "invalid repeat format"}
	}

}

func calcNextDateForDays(now time.Time, startDate time.Time, days int) string {
	nextDate := startDate
	for !nextDate.After(now) || !nextDate.After(startDate) {
		nextDate = nextDate.AddDate(0, 0, days)
	}
	return nextDate.Format("20060102")
}

func calcNextDateForYearly(now, startDate time.Time) string {
	nextDate := startDate
	for !nextDate.After(now) || !nextDate.After(startDate) {
		nextDate = nextDate.AddDate(1, 0, 0)
	}
	return nextDate.Format("20060102")
}

func parseDaysOfWeek(s string) (map[time.Weekday]bool, error) {
	parts := strings.Split(s, ",")
	daysOfWeek := make(map[time.Weekday]bool)
	for _, part := range parts {
		dayNum, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || dayNum < 1 || dayNum > 7 {
			return nil, &errors.ValidationError{Message: "invalid day of week"}
		}
		if dayNum == 7 {
			daysOfWeek[time.Sunday] = true
		} else {
			daysOfWeek[time.Weekday(dayNum)] = true
		}
	}
	return daysOfWeek, nil
}

func calcNextDateForWeekdays(now, startDate time.Time, daysOfWeek map[time.Weekday]bool) string {
	currentDate := startDate
	for {
		currentDate = currentDate.AddDate(0, 0, 1)
		if daysOfWeek[currentDate.Weekday()] && currentDate.After(now) {
			return currentDate.Format("20060102")
		}
	}
}

func parseDaysAndMonths(s string) (map[int]bool, map[int]bool, error) {
	parts := strings.Split(s, " ")
	dayParts := strings.Split(parts[0], ",")
	days := make(map[int]bool)
	for _, part := range dayParts {
		day, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || day > 31 || day < -2 {
			return nil, nil, &errors.ValidationError{Message: "invalid day of month"}
		}
		days[day] = true
	}

	months := make(map[int]bool)
	if len(parts) > 1 {
		monthParts := strings.Split(parts[1], ",")
		for _, part := range monthParts {
			month, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil || month < 1 || month > 12 {
				return nil, nil, &errors.ValidationError{Message: "invalid month"}
			}
			months[month] = true
		}
	}

	if len(months) == 0 {
		for m := 1; m <= 12; m++ {
			months[m] = true
		}
	}

	return days, months, nil
}

func daysInMonth(currentDay time.Time) int {
	return time.Date(currentDay.Year(), currentDay.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func calcNextDateForMonthly(now, startDate time.Time, daysOfMonth, months map[int]bool) string {
	currentDate := startDate
	for {
		currentDate = currentDate.AddDate(0, 0, 1)
		if _, monthOk := months[int(currentDate.Month())]; !monthOk {
			continue
		}
		day := currentDate.Day()
		_, dayOk := daysOfMonth[day]
		_, dayReverseOk := daysOfMonth[day-daysInMonth(currentDate)-1]
		if (dayOk || dayReverseOk) && currentDate.After(now) {
			return currentDate.Format("20060102")
		}
	}
}
