package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func generateAccountNumber() (string, error) {
	const prefix = "8801"
	const digits = 12

	max := new(big.Int).Exp(
		big.NewInt(10),
		big.NewInt(digits),
		nil,
	)

	number, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	randomDigits := number.String()
	randomDigits = strings.Repeat(
		"0",
		digits-len(randomDigits),
	) + randomDigits

	return fmt.Sprintf("%s%s", prefix, randomDigits), nil
}

func isAccountNumberDuplicate(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	var mysqlErr *mysqlDriver.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return strings.Contains(
			strings.ToLower(err.Error()),
			"account_number",
		)
	}

	return false
}
