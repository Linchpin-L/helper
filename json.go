package helper

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
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

// 只是用 年月日 的日期格式，即：2006-01-02
// 当从字符解析到时间格式时，时区将被设置为 Local
type Date time.Time

func (j Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j).Format("2006-01-02"))
}

func (j *Date) UnmarshalJSON(data []byte) error {
	if len(data) <= 2 {
		return errors.New("invalid time format")
	}
	if slices.Equal(data, []byte("null")) {
		return nil
	}
	if data[0] != '"' || data[len(data)-1] != '"' { // 检测数据格式
		return errors.New("invalid time format")
	}
	t, err := time.ParseInLocation("2006-01-02", string(data[1:len(data)-1]), time.Local)
	// t, err := time.Parse("2006-01-02", string(data[1:len(data)-1]))
	if err != nil {
		return err
	}
	*j = Date(t)
	return nil
}

func (j Date) Value() (driver.Value, error) {
	return time.Time(j), nil
}

func (j *Date) Scan(value interface{}) error {
	if value == nil {
		*j = Date(time.Time{})
		return nil
	}

	if t, ok := value.(time.Time); ok {
		*j = Date(t)
		return nil
	}
	return fmt.Errorf("无法扫描类型 %T 到 JSONDate", value)
}
