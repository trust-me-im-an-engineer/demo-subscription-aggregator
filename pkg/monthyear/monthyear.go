package monthyear

import (
	"encoding/json"
	"fmt"
	"time"
)

const DateLayout = "01-2006"

type MonthYear time.Time

func (my *MonthYear) UnmarshalJSON(b []byte) error {
	s := string(b)
	// Remove quotes from the string if present
	if len(s) > 0 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	t, err := time.Parse(DateLayout, s)
	if err != nil {
		return fmt.Errorf("invalid date format, expected MM-YYYY: %w", err)
	}
	*my = MonthYear(t)
	return nil
}

func (my *MonthYear) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*my).Format(DateLayout))
}
