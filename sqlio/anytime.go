package sqlio

import (
	"database/sql/driver"
	"strings"
	"time"

	"github.com/gurparit/go-common/logio"
)

type AnyTime struct {
	expect time.Time
}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	expectedValue := a.expect.String()[:19]
	actualValue := v.(time.Time).String()[:19]

	actualValue = strings.Replace(actualValue, "T", " ", 1)

	if expectedValue == actualValue {
		return true
	}

	logio.Println("AnyTime{}: expect [%s], actual [%s]", expectedValue, actualValue)

	return false
}
