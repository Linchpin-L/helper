package helper

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// 类似于: "1,2,3,4,5"
type MutipleUint64 []uint64

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *MutipleUint64) Scan(value any) error {
	bytess, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to parse MutipleUint64:", value))
	}

	temp := make([]uint64, 0)

	if len(bytess) > 0 { // 避免产生一个长度为 1 仅包含 空串 的数组
		plain := string(bytess)
		plains := strings.Split(plain, ",")
		for _, v := range plains {
			t, _ := strconv.ParseUint(v, 10, 64)
			temp = append(temp, t)
		}
	}

	*j = temp
	return nil
}

// Value return json value, implement driver.Valuer interface
func (j MutipleUint64) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "", nil
	}

	temp := make([]string, len(j))
	for i := range j {
		temp[i] = strconv.FormatUint(j[i], 10)
	}

	return strings.Join(temp, ","), nil
}

// 计算平均分
// 索引视为分数+1, 值视为人数. 求人均分
func (j MutipleUint64) AverageScore() float64 {
	var total, count uint64
	for score, num := range j {
		total += num * uint64(score+1)
		count += num
	}
	if count > 0 {
		return float64(total) / float64(count)
	}

	return 0
}

// 计算和
// 在某些场景下, 即统计总PV
func (j MutipleUint64) Sum() (total uint64) {
	for i := range j {
		total += j[i]
	}
	return
}