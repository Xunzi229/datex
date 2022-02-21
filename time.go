package datex

import (
    "database/sql/driver"
    "fmt"
    "time"
)

const (
    DateLayout = "2006-01-02"
    TimeLayout = "2006-01-02 15:04:05"
)

type TimeX struct {
    Valid bool
    
    At time.Time
}

// Scan
// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *TimeX) Scan(value interface{}) error {
    var err error
    
    if value == nil {
        t.At, t.Valid = time.Time{}, false
        return nil
    }
    t.At, t.Valid = value.(time.Time)
    if !t.Valid {
        err = fmt.Errorf("value is not time.Time, value: %v", value)
    }
    return err
}

// Value
// 实现 driver.Valuer 接口，Value 返回 json value
// 写入数据库之前，对数据做类型转换
func (t TimeX) Value() (driver.Value, error) {
    if !t.Valid {
        return nil, nil
    }
    
    stamp := fmt.Sprintf("%s", t.At.Format(TimeLayout))
    return stamp, nil
}

func NewTimeX(t time.Time) TimeX {
    return TimeX{At: t, Valid: true}
}

func LoadTimeByLayout(t string, layout string) TimeX {
    tx, err := time.ParseInLocation(layout, t, time.Local)
    if err != nil {
        return TimeX{Valid: false}
    }
    return NewTimeX(tx)
}

func LoadTimeByYmd(t string) TimeX {
    tx, err := time.ParseInLocation(TimeLayout, t, time.Local)
    if err != nil {
        return TimeX{Valid: false}
    }
    return NewTimeX(tx)
}

func (t TimeX) MarshalJSON() ([]byte, error) {
    stamp := fmt.Sprintf("\"%s\"", t.At.Format(TimeLayout))
    return []byte(stamp), nil
}

func (t *TimeX) UnmarshalJSON(data []byte) error {
    var err error
    
    t.At, err = time.ParseInLocation("\"2006-01-02 15:04:05\"", string(data), time.Local)
    if err != nil {
        return err
    }
    t.Valid = true
    return nil
}

func (t *TimeX) AddDate(years int, months int, days int) time.Time {
    return t.At.AddDate(years, months, days)
}

func (t *TimeX) Format(layout string) string {
    return t.At.Format(layout)
}

func (t *TimeX) Equal(t1 TimeX) bool {
    if (t.At.Format(TimeLayout) != t1.At.Format(TimeLayout)) || t.Valid != t1.Valid {
        return false
    }
    return true
}
