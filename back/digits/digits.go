package digits

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetConsciousnessNumber(birthdate string) (int, error) {
	parts := strings.Split(birthdate, "-")
	v1, err := strconv.Atoi(string(parts[2][0]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}
	v2, err := strconv.Atoi(string(parts[2][1]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}

	resp, err := reduction(v1 + v2)
	if err != nil {
		return 0, err
	}
	return resp, nil

}

func GetActionNumber(birthdate string) (int, error) {
	parts := strings.Split(birthdate, "-")
	d1, err := strconv.Atoi(string(parts[2][0]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}
	d2, err := strconv.Atoi(string(parts[2][1]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}

	m1, err := strconv.Atoi(string(parts[1][0]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}
	m2, err := strconv.Atoi(string(parts[1][1]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}

	y1, err := strconv.Atoi(string(parts[0][0]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}
	y2, err := strconv.Atoi(string(parts[0][1]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}
	y3, err := strconv.Atoi(string(parts[0][2]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}
	y4, err := strconv.Atoi(string(parts[0][3]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}

	resp, err := reduction(d1 + d2 + m1 + m2 + y1 + y2 + y3 + y4)
	if err != nil {
		return 0, err
	}
	return resp, nil
}

func GetKarmaNumber(birthdate string) (int, error) {
	parts := strings.Split(birthdate, "-")
	d1, err := strconv.Atoi(string(parts[2][0]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}
	d2, err := strconv.Atoi(string(parts[2][1]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}

	m1, err := strconv.Atoi(string(parts[1][0]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}
	m2, err := strconv.Atoi(string(parts[1][1]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}

	resp, err := reduction(d1 + d2 + m1 + m2)
	if err != nil {
		return 0, err
	}

	return resp, nil
}

func GetMonthNumber(birthdate string) (int, error) {
	parts := strings.Split(birthdate, "-")
	m1, err := strconv.Atoi(string(parts[1][0]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}
	m2, err := strconv.Atoi(string(parts[1][1]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}

	resp, err := reduction(m1 + m2)
	if err != nil {
		return 0, err
	}

	return resp, nil
}

func GetYearNumber(birthdate string) (int, error) {
	parts := strings.Split(birthdate, "-")
	currentYear := strconv.Itoa(time.Now().Year())

	d1, err := strconv.Atoi(string(parts[2][0]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}
	d2, err := strconv.Atoi(string(parts[2][1]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}

	m1, err := strconv.Atoi(string(parts[1][0]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}
	m2, err := strconv.Atoi(string(parts[1][1]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}

	y1, err := strconv.Atoi(string(currentYear[0]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}
	y2, err := strconv.Atoi(string(currentYear[1]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}
	y3, err := strconv.Atoi(string(currentYear[2]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}
	y4, err := strconv.Atoi(string(currentYear[3]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}

	resp, err := reduction(d1 + d2 + m1 + m2 + y1 + y2 + y3 + y4)
	if err != nil {
		return 0, err
	}

	return resp, nil
}

func GetCommonDay() (int, error) {
	strtime := time.Now().String()
	resp, err := strconv.Atoi(string(strtime[9]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}

	return resp, nil
}

func GetPrivateDay(birthdate string) (int, error) {
	strtime := time.Now().String()
	d1, err := strconv.Atoi(string(strtime[8]))
	if err != nil {
		return 0, fmt.Errorf("can`t parce current time: %w", err)
	}
	d2, err := strconv.Atoi(string(strtime[9]))
	if err != nil {
		return 0, fmt.Errorf("can`t parce current time: %w", err)
	}

	m1, err := strconv.Atoi(string(strtime[5]))
	if err != nil {
		return 0, fmt.Errorf("can`t parce current time: %w", err)
	}
	m2, err := strconv.Atoi(string(strtime[6]))
	if err != nil {
		return 0, fmt.Errorf("can`t parce current time: %w", err)
	}

	y, err := GetYearNumber(birthdate)
	if err != nil {
		return 0, err
	}
	resp, err := reduction(d1 + d2 + m1 + m2 + y)
	if err != nil {
		return 0, err
	}

	return resp, nil

}

func reduction(resp int) (int, error) {
	for resp > 9 {
		r1, err := strconv.Atoi(string(strconv.Itoa(resp)[0]))
		if err != nil {
			return 0, fmt.Errorf("error on reduction: %w", err)
		}
		r2, err := strconv.Atoi(string(strconv.Itoa(resp)[1]))
		if err != nil {
			return 0, fmt.Errorf("error on reduction: %w", err)
		}
		resp = r1 + r2
	}
	return resp, nil
}
