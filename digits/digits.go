package digits

import (
	"fmt"
	"strconv"
	"strings"
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

	return v1 + v2, nil

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

	pre := d1 + d2 + m1 + m2 + y1 + y2 + y3 + y4
	strpre := strconv.Itoa(pre)

	pre1, err := strconv.Atoi(string(strpre[0]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}
	pre2, err := strconv.Atoi(string(strpre[1]))
	if err != nil {
		return 0, fmt.Errorf("invalid day in birthdate: %w", err)
	}

	return pre1 + pre2, nil
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

	return d1 + d2 + m1 + m2, nil
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

	return m1 + m2, nil
}
