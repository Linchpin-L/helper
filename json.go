package helper

import (
	"errors"
	"strconv"
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
