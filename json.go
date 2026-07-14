package helper

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// json 解析前，可以是字符串形式，也可能是数字形式
// 都会转换为 int 类型
type UnstableInt int

func (u *UnstableInt) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	// 尝试移除首尾的引号
	if data[0] == '"' && data[len(data)-1] == '"' {
		if len(data) == 1 {
			return errors.New("invalid JSON string")
		}
		if len(data) == 2 {
			return nil
		}
		data = data[1 : len(data)-1]
	}
	temp, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}
	*u = UnstableInt(temp)
	return nil
}

func (u UnstableInt) MarshalJSON() ([]byte, error) {
	num := strconv.Itoa(int(u))
	return []byte(num), nil
}

// 接口返回值有时会返回字符串形式的浮点数值
type UnstableFloat float64

// 实现 json 解析接口
func (u *UnstableFloat) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if data[0] == '"' && data[len(data)-1] == '"' {
		if len(data) == 1 {
			return errors.New("invalid JSON string")
		}
		if len(data) == 2 {
			return nil
		}
		data = data[1 : len(data)-1]
	}
	temp, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return err
	}
	*u = UnstableFloat(temp)
	return nil
}

// 实现 json 序列化接口
func (u UnstableFloat) MarshalJSON() ([]byte, error) {
	num := strconv.FormatFloat(float64(u), 'f', -1, 64)
	return []byte(num), nil
}

// Date 表示一个仅包含“年月日”的日期，可用于映射数据库 DATE 字段。
//
// 与 time.Time 不同，Date 不表示一个时间点（Instant），而是一个没有时区、没有时分秒概念的日期值（Value Object）。
//
// 内部实现
//
//   - 内部借用 time.Time 作为存储结构。
//   - 所有合法的 Date 都统一使用 UTC 00:00:00.000000000 保存。
//   - UTC 仅作为内部规范化表示，不具有任何业务上的时区含义。
//
// 使用约定
//
//   - 应通过 NewDate、DateFromTime、Scan、UnmarshalJSON 等方法创建 Date。
//   - 不应直接通过 Date(time.Time) 构造，否则可能破坏内部约束。
//   - 需要进行时间计算时，应先调用 UTC() 或者 InLocation 转为 time.Time，计算完成后再使用
//     DateFromTime() 转回 Date。
//   - 数据库存储时输出 "2006-01-02" 格式，不参与任何时区转换。
//   - JSON 序列化格式为 "2006-01-02"，零值序列化为 null。
//
// 不变式（Invariant）
//
// 一个合法的 Date 必须始终满足：
//
//   - Location == UTC
//   - Hour == Minute == Second == Nanosecond == 0
//
// 因此，Date 可以安全地进行值比较（==），也能够保证在不同机器、不同时区下具有一致的内部表示。
type Date time.Time

// NewDate 创建一个日期。
func NewDate(year int, month time.Month, day int) Date {
	return Date(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

// DateFromTime 从 time.Time 提取日期部分。
func DateFromTime(t time.Time) Date {
	y, m, d := t.Date()
	return NewDate(y, m, d)
}

// UTC 返回底层 time.UTC（UTC 零点）。
//
// 返回值仅用于时间计算或与标准库交互。
// 如果计算结果仍表示一个日期，应使用 DateFromTime() 转换回 Date，
// 以保证 Date 的内部约束不被破坏。
func (d Date) UTC() time.Time {
    return time.Time(d)
}

// InLocation 返回本地零点的时间
func (d Date) InLocation(loc *time.Location) time.Time {
	t := d.UTC()

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0,
		0,
		0,
		0,
		loc,
	)
}

// String 返回 yyyy-MM-dd。
func (d Date) String() string {
	t := time.Time(d)
	if t.IsZero() {
		return ""
	}
	return t.Format(time.DateOnly)
}

// Year 返回年份。
func (d Date) Year() int {
	return time.Time(d).Year()
}

// Month 返回月份。
func (d Date) Month() time.Month {
	return time.Time(d).Month()
}

// Day 返回日期。
func (d Date) Day() int {
	return time.Time(d).Day()
}

func (j Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j).Format(time.DateOnly))
}

func (d *Date) UnmarshalJSON(data []byte) error {
    if string(data) == "null" {
        *d = Date{}
        return nil
    }

    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }

    t, err := time.Parse(time.DateOnly, s)
    if err != nil {
        return err
    }

    *d = DateFromTime(t)
    return nil
}

// 只响应一个日期
func (j Date) Value() (driver.Value, error) {
	t := time.Time(j)
	return t.Format(time.DateOnly), nil
}

func (d *Date) Scan(value any) error {
	if value == nil {
		*d = Date{}
		return nil
	}

	switch v := value.(type) {

	case time.Time:
		*d = DateFromTime(v)
		return nil

	case string:
		t, err := time.Parse(time.DateOnly, v)
		if err != nil {
			return err
		}
		*d = DateFromTime(t)
		return nil

	case []byte:
		t, err := time.Parse(time.DateOnly, string(v))
		if err != nil {
			return err
		}
		*d = DateFromTime(t)
		return nil
	}

	return fmt.Errorf("cannot scan %T into Date", value)
}

type DateTime time.Time

func (t *DateTime) GetFormat() string {
	return time.DateTime
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(t.GetFormat())+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, t.GetFormat())
	b = append(b, '"')
	return b, nil
}

func (t *DateTime) UnmarshalJSON(bs []byte) error {
	var s string
	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	tt, err := time.ParseInLocation(t.GetFormat(), s, time.Local)
	if err != nil {
		return err
	}
	*t = DateTime(tt)
	return nil
}
