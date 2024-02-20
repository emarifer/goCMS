package timezone_conversion

import "time"

func ConvertDateTime(tz string, dt time.Time) string {
	loc, _ := time.LoadLocation(tz)

	return dt.In(loc).Format(time.RFC822Z)
}
