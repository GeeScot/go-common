package sqlio

import (
	"database/sql/driver"
	"strings"
	"time"

	"fmt"
	"errors"
)

type AnyTime struct {
	expect time.Time
}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) (bool, error) {
	expectedValue := a.expect.String()[:19]
	actualValue := v.(time.Time).String()[:19]

	actualValue = strings.Replace(actualValue, "T", " ", 1)

	if expectedValue == actualValue {
		return true, nil
	}

	return false, errors.New(fmt.Sprintf("AnyTime{}: expect [%s], actual [%s]", expectedValue, actualValue))
}
