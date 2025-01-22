package helper

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestIsIDCard(t *testing.T) {
	type args struct {
		idcard string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "1", args: args{""}, want: false},
		{name: "2", args: args{"1"}, want: false},
		{name: "3", args: args{"440102198001021231"}, want: false},
		{name: "4", args: args{"440102198001021230"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIDCard(tt.args.idcard); got != tt.want {
				t.Errorf("IsIDCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeroClock(t *testing.T) {
	a := time.Date(2023, 1, 11, 12, 12, 12, 12, time.Local)
	b := time.Date(2024, 2, 29, 23, 59, 59, 59, time.Local)
	c := time.Date(2024, 2, 29, 0, 0, 0, 12, time.Local)
	e := time.Now()

	type args struct {
		t *time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{name: "1", args: args{&a}, want: time.Date(2023, 1, 11, 0, 0, 0, 0, time.Local)},
		{name: "2", args: args{&b}, want: time.Date(2024, 2, 29, 0, 0, 0, 0, time.Local)},
		{name: "3", args: args{&c}, want: time.Date(2024, 2, 29, 0, 0, 0, 0, time.Local)},
		{name: "4", args: args{nil}, want: time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, time.Local)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ZeroClock(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZeroClock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumToAlphabet(t *testing.T) {
	type args struct {
		num int64
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"1", args{0}, []byte("0")},
		{"2", args{-0}, []byte("0")},
		{"3", args{1}, []byte("1")},
		{"4", args{-1}, []byte("-1")},
		{"6", args{10}, []byte("a")},
		{"7", args{11}, []byte("b")},
		{"8", args{61}, []byte("1p")},
		{"9", args{62}, []byte("1q")},
		{"10", args{-653543}, []byte("-e09z")},
		{"5", args{9223372036854775807}, []byte("1y2p0ij32e8e7")},
		{"11", args{-9223372036854775808}, []byte("-1y2p0ij32e8e8")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumToAlphabet(tt.args.num); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NumToAlphabet3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeDir(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{"test/a"}, false},
		{"2", args{"test/a/b/c"}, false},
		{"3", args{"test/a/b/c/d.txt"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MakeDir(tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("MakeDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMakeDirTrimFileName(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	type args struct {
		dir string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantPath string
	}{
		// WARNING: 测试后会删除新建的目录，测试用例请小心不要误删文件
		{"0", args{"/testasjldjfiwsersmx"}, false, "/"},
		{"1", args{root + "/testasjldjfiwsersmx/testasjldjfiwsersmx.txt"}, false, root + "/testasjldjfiwsersmx"},
		{"2", args{"test.txt"}, false, ""},
		{"3", args{"test/a"}, false, "test/a"},
		{"4", args{"test/b/b.txt"}, false, "test/b"},
		{"5", args{"test/c.txt/d.txt"}, false, "test/c.txt"},
		{"6", args{"test/d/a/l/s/dlfiejrse.m.s.pdf"}, false, "test/d/a/l/s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alreadyExist := IsFileExist(tt.wantPath) // 已经存在的目录不删除，避免因测试用例错误导致真实文件被删除
			// if alreadyExist {
			// 	fmt.Println("已经存在的目录 不删除")
			// }

			if err := MakeDirTrimFileName(tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("MakeDirTrimFileName() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.wantPath == "" {
					return
				}
				// 检查目录是否存在
				if !IsFileExist(tt.wantPath) {
					t.Errorf("MakeDirTrimFileName() wantPath %s, but it is not here", tt.wantPath)
				}
				// 删除测试目录。如果目录下有其他文件，该函数将会报错，可以用来检查是否错误的生成了多余的文件夹
				if !alreadyExist {
					err := os.Remove(tt.wantPath)
					if err != nil {
						t.Error(err)
					}
				}
			}
		})
	}
	err = os.RemoveAll("test")
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveElementIgnoreOrder(t *testing.T) {
	type args struct {
		slice []any
		i     int
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		{"1", args{[]any{}, 1}, []any{}},
		{"2", args{nil, 1}, nil},
		{"3", args{[]any{1, 2, 3}, 1}, []any{1, 3}},
		{"4", args{[]any{1, 2, 3}, 2}, []any{1, 2}},
		{"5", args{[]any{1, 2, 3}, 3}, []any{1, 2, 3}},
		{"6", args{[]any{1, 2, 3}, -1}, []any{1, 2, 3}},
		{"7", args{[]any{"a", "b", "c"}, 1}, []any{"a", "c"}},
		{"8", args{[]any{true, false, true}, 0}, []any{true, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveElementIgnoreOrder(tt.args.slice, tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveElementWithoutOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGetFirstDayOfWeek(t *testing.T) {
	tests := []struct {
		name string
		today time.Time
		want time.Time
	}{
		{
			name: "Monday",
			today: time.Date(2023, 10, 2, 12, 0, 0, 0, time.Local), // Monday
			want: time.Date(2023, 10, 2, 0, 0, 0, 0, time.Local),
		},
		{
			name: "Wednesday",
			today: time.Date(2023, 10, 4, 12, 0, 0, 0, time.Local), // Wednesday
			want: time.Date(2023, 10, 2, 0, 0, 0, 0, time.Local),
		},
		{
			name: "Sunday",
			today: time.Date(2023, 10, 8, 12, 23, 59, 59, time.Local), // Sunday
			want: time.Date(2023, 10, 2, 0, 0, 0, 0, time.Local),
		},
		{
			name: "Saturday",
			today: time.Date(2023, 10, 7, 12, 1, 1, 1, time.Local), // Saturday
			want: time.Date(2023, 10, 2, 0, 0, 0, 0, time.Local),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFirstDayOfWeek(tt.today); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFirstDayOfWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

