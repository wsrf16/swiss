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

type TimeFormat string

const (
	NormalFormat      TimeFormat = "2006-01-02 15:04:05"
	DateFormat                   = "2006-01-02"
	TightNormalFormat            = "20060102150405"
	TightDateFormat              = "2006-01-02 15:04:05"
)

func Format(t time.Time, f TimeFormat) string {
	return t.Format(string(f))
}
