package helper

import (
	"bytes"
	"encoding/json"
	"math"
	"strconv"
	"testing"
	"time"
)

func Test_unstableInt(t *testing.T) {
	tests := []struct {
		name    string
		in      []byte
		out     int
		wantErr bool
	}{
		{"1", []byte(`"1"`), 1, false},
		{"2", []byte(`""`), 0, false},
		{"3", []byte(`"9223372036854775807"`), math.MaxInt, false},
		{"4", []byte(`"-9223372036854775808"`), math.MinInt, false},
		{"5", []byte(`"0"`), 0, false},
		{"6", []byte(`"9.9"`), 0, true},
		{"6", []byte(`"`), 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := new(UnstableInt)
			if err := temp.UnmarshalJSON(tt.in); (err != nil) != tt.wantErr {
				t.Errorf("unstableInt.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.out != int(*temp) {
				t.Errorf("unstableInt.UnmarshalJSON() = %v, want %v", *temp, tt.out)
			}
			out, err := temp.MarshalJSON()
			if err != nil {
				t.Errorf("unstableInt.MarshalJSON() = %v, want %v", out, tt.out)
			}
			if !bytes.Equal(out, []byte(strconv.Itoa(tt.out))) {
				t.Errorf("unstableInt.MarshalJSON() = %v, want %v", out, tt.in)
			}
		})
	}

	type f struct {
		in  []byte
		out []byte
	}
	tests2 := []f{
		{[]byte(`{"A":"1"}`), []byte(`{"A":1}`)},
		{[]byte(`{"A":""}`), []byte(`{"A":0}`)},
		{[]byte(`{}`), []byte(`{"A":null}`)},
		{[]byte(`{"A":"9223372036854775807"}`), []byte(`{"A":9223372036854775807}`)},
		{[]byte(`{"A":"-1"}`), []byte(`{"A":-1}`)},
	}
	for _, tt := range tests2 {
		t.Run(string(tt.in), func(t *testing.T) {
			plain, plains := tt.in, tt.out
			var s struct {
				A *UnstableInt
			}
			// fmt.Println(1, s.A)
			err := json.Unmarshal(plain, &s)
			if err != nil {
				t.Errorf("unstableInt.UnmarshalJSON() error = %v", err)
			}
			// fmt.Println(2, s.A)

			b, err := json.Marshal(s)
			if err != nil {
				t.Errorf("unstableInt.MarshalJSON() error = %v", err)
			}
			if !bytes.Equal(b, plains) {
				t.Errorf("unstableInt.MarshalJSON() = %v, want %v", b, plains)
			}
		})
	}
}

func Test_UnstableFloat(t *testing.T) {
	tests := []struct {
		name    string
		in      []byte
		out     float64
		wantErr bool
	}{
		{"1", []byte(`"1.0"`), 1, false},
		{"2", []byte(`""`), 0, false},
		{"3", []byte(`"1.79769e+308"`), 1.79769e+308, false},
		{"5", []byte(`"0.0"`), 0, false},
		{"6", []byte(`"9.99"`), 9.99, false},
		{"7", []byte(`"`), 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := new(UnstableFloat)
			if err := temp.UnmarshalJSON(tt.in); (err != nil) != tt.wantErr {
				t.Errorf("UnstableFloat.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.out != float64(*temp) {
				t.Errorf("UnstableFloat.UnmarshalJSON() = %v, want %v", *temp, tt.out)
			}
			out, err := temp.MarshalJSON()
			if err != nil {
				t.Errorf("UnstableFloat.MarshalJSON() = %v, want %v", out, tt.out)
			}
			if !bytes.Equal(out, []byte(strconv.FormatFloat(tt.out, 'f', -1, 64))) {
				t.Errorf("UnstableFloat.MarshalJSON() = %v, want %v", out, tt.in)
			}
		})
	}

	type f struct {
		in  []byte
		out []byte
	}
	tests2 := []f{
		{[]byte(`{"A":"1.1"}`), []byte(`{"A":1.1}`)},
		{[]byte(`{"A":""}`), []byte(`{"A":0}`)},
		{[]byte(`{}`), []byte(`{"A":null}`)},
		{[]byte(`{"A":"922.22"}`), []byte(`{"A":922.22}`)},
		{[]byte(`{"A":"-10.0"}`), []byte(`{"A":-10}`)},
	}
	for _, tt := range tests2 {
		t.Run(string(tt.in), func(t *testing.T) {
			plain, plains := tt.in, tt.out
			var s struct {
				A *UnstableFloat
			}
			// fmt.Println(1, s.A)
			err := json.Unmarshal(plain, &s)
			if err != nil {
				t.Errorf("UnstableFloat.UnmarshalJSON() error = %v", err)
			}
			// fmt.Println(2, s.A)

			b, err := json.Marshal(s)
			if err != nil {
				t.Errorf("UnstableFloat.MarshalJSON() error = %v", err)
			}
			if !bytes.Equal(b, plains) {
				t.Errorf("UnstableFloat.MarshalJSON() = %s, want %s", string(b), string(plains))
			}
		})
	}
}

// 测试 MarshalJSON 方法
func TestDate_MarshalJSON(t *testing.T) {
    tests := []struct {
        name     string
        date     Date
        expected string
        wantErr  bool
    }{
        {
            name:     "正常日期",
            date:     Date(time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)),
            expected: `"2023-12-25"`,
            wantErr:  false,
        },
        {
            name:     "零值时间",
            date:     Date(time.Time{}),
            expected: `"0001-01-01"`,
            wantErr:  false,
        },
        {
            name:     "闰年日期",
            date:     Date(time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)),
            expected: `"2024-02-29"`,
            wantErr:  false,
        },
        {
            name:     "LOCAL时间",
            date:     Date(time.Date(2024, 1, 1, 1, 1, 0, 0, time.Local)),
            expected: `"2024-01-01"`,
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := tt.date.MarshalJSON()
            if (err != nil) != tt.wantErr {
                t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if string(got) != tt.expected {
                t.Errorf("MarshalJSON() got = %s, want %s", got, tt.expected)
            }
        })
    }
}

// 测试 UnmarshalJSON 方法
func TestDate_UnmarshalJSON(t *testing.T) {
    tests := []struct {
        name     string
        jsonStr  string
        expected Date
        wantErr  bool
    }{
        {
            name:     "正常日期解析",
            jsonStr:  `"2023-12-25"`,
            expected: Date(time.Date(2023, 12, 25, 0, 0, 0, 0, time.Local)),
            wantErr:  false,
        },
        {
            name:     "无效日期格式",
            jsonStr:  `"2023-13-45"`,
            expected: Date{},
            wantErr:  true,
        },
        {
            name:     "空字符串",
            jsonStr:  `""`,
            expected: Date{},
            wantErr:  true,
        },
        {
            name:     "非字符串类型",
            jsonStr:  `12345`,
            expected: Date{},
            wantErr:  true,
        },
        {
            name:     "null值",
            jsonStr:  `null`,
            expected: Date{},
            wantErr:  false, // JSON null 会解析为零值
        },
        {
            name:     "带时间的日期字符串",
            jsonStr:  `"2023-12-25T15:04:05Z"`,
            expected: Date{},
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var d Date
            err := d.UnmarshalJSON([]byte(tt.jsonStr))
            
            if (err != nil) != tt.wantErr {
                t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !tt.wantErr {
                // 比较时间部分（忽略时区差异）
                expectedTime := time.Time(tt.expected).UTC()
                actualTime := time.Time(d).UTC()
                
                if !expectedTime.Equal(actualTime) {
                    t.Errorf("UnmarshalJSON() got = %v, want %v", 
                        actualTime.Format("2006-01-02"),
                        expectedTime.Format("2006-01-02"))
                }
            }
        })
    }
}

// 测试 Value 方法（用于数据库写入）
func TestDate_Value(t *testing.T) {
    tests := []struct {
        name     string
        date     Date
        expected interface{}
        wantErr  bool
    }{
        {
            name:     "正常日期",
            date:     Date(time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)),
            expected: time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC),
            wantErr:  false,
        },
        {
            name:     "零值时间",
            date:     Date(time.Time{}),
            expected: time.Time{},
            wantErr:  false,
        },
        {
            name:     "偏移时间",
            date:     Date(time.Date(2025, 12, 1, 23, 33, 33, 0, time.UTC)),
            expected: time.Date(2025, 12, 2, 7, 33, 33, 0, time.Local),
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := tt.date.Value()
            if (err != nil) != tt.wantErr {
                t.Errorf("Value() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !tt.wantErr {
                expectedTime := tt.expected.(time.Time)
                actualTime := got.(time.Time)
                
                if !expectedTime.Equal(actualTime) {
                    t.Errorf("Value() got = %v, want %v", actualTime, expectedTime)
                }
            }
        })
    }
}

// 测试 Scan 方法（用于数据库读取）
func TestDate_Scan(t *testing.T) {
    tests := []struct {
        name     string
        input    interface{}
        expected Date
        wantErr  bool
    }{
        {
            name:     "time.Time类型",
            input:    time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC),
            expected: Date(time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)),
            wantErr:  false,
        },
        {
            name:     "nil值",
            input:    nil,
            expected: Date(time.Time{}),
            wantErr:  false,
        },
        {
            name:     "字符串类型（应该出错）",
            input:    "2023-12-25",
            expected: Date{},
            wantErr:  true,
        },
        {
            name:     "整数类型（应该出错）",
            input:    20231225,
            expected: Date{},
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var d Date
            err := d.Scan(tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !tt.wantErr {
                expectedTime := time.Time(tt.expected).UTC()
                actualTime := time.Time(d).UTC()
                
                if !expectedTime.Equal(actualTime) {
                    t.Errorf("Scan() got = %v, want %v", 
                        actualTime.Format("2006-01-02"),
                        expectedTime.Format("2006-01-02"))
                }
            }
        })
    }
}

// 测试 JSON 完整的序列化和反序列化
func TestDate_JSONRoundTrip(t *testing.T) {
    tests := []struct {
        name string
        date Date
    }{
        {
            name: "普通日期",
            date: Date(time.Date(2023, 12, 25, 0, 0, 0, 0, time.Local)),
        },
        {
            name: "零值日期",
            date: Date(time.Date(1, 1, 1, 0, 0, 0, 0, time.Local)),
        },
        {
            name: "闰年日期",
            date: Date(time.Date(2024, 2, 29, 0, 0, 0, 0, time.Local)),
        },
    }

    for _, tt := range tests { 
		// 注意此检查方法在 Marshal 的时候一定会丢失时间精度，因为时区丢失了
		// 而在 Unmarshal 时，我们均以当前的时区来识别时间
        t.Run(tt.name, func(t *testing.T) {
            // 序列化
            jsonBytes, err := json.Marshal(tt.date)
            if err != nil {
                t.Fatalf("Marshal failed: %v", err)
            }
            
            // 反序列化
            var newDate Date
            err = json.Unmarshal(jsonBytes, &newDate)
            if err != nil {
                t.Fatalf("Unmarshal failed: %v", err)
            }
            
            // 比较
            originalTime := time.Time(tt.date)
            newTime := time.Time(newDate)
            
            if !originalTime.Equal(newTime) {
                t.Errorf("RoundTrip mismatch: original = %v, after roundtrip = %v",
                    originalTime,
                    newTime)
            }
        })
    }
}

// 测试数据库值的完整往返
func TestDate_DatabaseRoundTrip(t *testing.T) {
    originalDate := Date(time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC))
    
    // 获取数据库值
    dbValue, err := originalDate.Value()
    if err != nil {
        t.Fatalf("Value() failed: %v", err)
    }
    
    // 从数据库值扫描回来
    var scannedDate Date
    err = scannedDate.Scan(dbValue)
    if err != nil {
        t.Fatalf("Scan() failed: %v", err)
    }
    
    // 比较
    originalTime := time.Time(originalDate).UTC()
    scannedTime := time.Time(scannedDate).UTC()
    
    if !originalTime.Equal(scannedTime) {
        t.Errorf("Database roundtrip mismatch: original = %v, scanned = %v",
            originalTime.Format("2006-01-02"),
            scannedTime.Format("2006-01-02"))
    }
}
