package optional

import (
	"time"
)

func FromString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func String(str string) *string {
	return &str
}

func FromBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func Bool(b bool) *bool {
	return &b
}

func Int64(v int64) *int64 {
	return &v
}

func Int(v int) *int {
	return &v
}

func Time(v time.Time) *time.Time {
	return &v
}

func FromTime(v *time.Time) time.Time {
	if v == nil {
		return time.Time{}
	}
	return *v
}

func FromDuration(v *time.Duration) (t time.Duration) {
	if v == nil {
		return t
	}
	return *v
}

func Duration(v time.Duration) *time.Duration {
	return &v
}

func FromBytes(v *[]byte) []byte {
	if v == nil {
		return nil
	}
	return *v
}

func Bytes(v []byte) *[]byte {
	return &v
}

func Int64ToDuration(v *int64) *time.Duration {
	if v == nil {
		return nil
	}
	return Duration(time.Duration(*v) * time.Millisecond)
}

func DurationToInt64(v *time.Duration) *int64 {
	if v == nil {
		return nil
	}
	ms := v.Milliseconds()
	return &ms
}
