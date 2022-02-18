package datex

import (
    "database/sql/driver"
    "fmt"
    "time"
)

type DateX struct {
    At    time.Time
    Valid bool
}

// Scan
// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *DateX) Scan(value interface{}) error {
    var err error
    
    if value == nil {
        return nil
    }
    t.At, t.Valid = value.(time.Time)
    if t.Valid {
        err = fmt.Errorf("value is not time.Time, value: %v", value)
    }
    return err
}

// Value
// 实现 driver.Valuer 接口，Value 返回 json value
// 写入数据库之前，对数据做类型转换
func (t DateX) Value() (driver.Value, error) {
    if !t.Valid {
        return nil, nil
    }
    
    stamp := fmt.Sprintf("%s", t.At.Format(DateLayout))
    return stamp, nil
}

func NewDateX(t time.Time) DateX {
    return DateX{At: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()), Valid: true}
}

func LoadDateByLayout(str string, layout string) DateX {
    tx, err := time.ParseInLocation(layout, str, time.Local)
    if err != nil {
        return DateX{Valid: false}
    }
    return NewDateX(tx)
}

func LoadDateByYmd(t string) DateX {
    tx, err := time.ParseInLocation(DateLayout, t, time.Local)
    if err != nil {
        return DateX{Valid: false}
    }
    return NewDateX(tx)
}

func (t DateX) MarshalJSON() ([]byte, error) {
    stamp := fmt.Sprintf("\"%s\"", t.At.Format(DateLayout))
    return []byte(stamp), nil
}

func (t *DateX) UnmarshalJSON(data []byte) error {
    var err error
    
    t.At, err = time.ParseInLocation(`"2006-01-02"`, string(data), time.Local)
    if err != nil {
        return err
    }
    t.Valid = true
    return nil
}

func (t *DateX) AddDate(years int, months int, days int) time.Time {
    return t.At.AddDate(years, months, days)
}

func (t *DateX) Add(years int, months int, days int) DateX {
    return NewDateX(t.AddDate(years, months, days))
}

func (t *DateX) Format(layout string) string {
    return t.At.Format(layout)
}

func (t DateX) Equal(t1 DateX) bool {
    return t.Format(DateLayout) == t1.Format(DateLayout)
}
