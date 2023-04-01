package timekit

import "time"

func TimeStampNanoSecond() int64 {
	return time.Now().UnixNano()
}

func TimeStampMilliSecond() int64 {
	return time.Now().UnixNano() / 1e6
}

func TimeStampSecond() int64 {
	return time.Now().Unix()
}

func ToMilliNanoSecond(t time.Time) int64 {
	return t.UnixNano()
}

func ToMilliSecond(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func ToSecond(t time.Time) int64 {
	return t.Unix()
}

func Watch() time.Duration {
	start := time.Now()
	enhause := time.Since(start)
	// enhause = time.Now().Sub(start)
	return enhause
}

func Time1Year() time.Time {
	return time.Now().Add(365 * 24 * time.Hour)
}

func Time3Minutes() time.Time {
	return time.Now().Add(3 * time.Minute)
}

func Time1Hour() time.Time {
	return time.Now().Add(1 * time.Hour)
}

type TimeFormat = string

const (
	DateTimeFormat              TimeFormat = "2006-01-02 15:04:05"
	DateTimeTightFormat                    = "20060102150405"
	DateTimeSlashFormat                    = "2006/01/02 15:04:05"
	DateOnlyFormat                         = "2006-01-02"
	TimeOnlyFormat                         = "15:04:05"
	DateTimeDateOnlySlashFormat            = "2006/01/02"
)

func Format(t time.Time, f TimeFormat) string {
	return t.Format(string(f))
}

type NormalTime time.Time

func (t NormalTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(DateTimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, DateTimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t *NormalTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+DateTimeFormat+`"`, string(data), time.Local)
	*t = NormalTime(now)
	return
}

func (t NormalTime) String() string {
	return time.Time(t).Format(DateTimeFormat)
}
