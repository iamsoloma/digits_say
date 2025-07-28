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
