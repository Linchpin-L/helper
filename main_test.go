package helper

import "testing"

func TestIsIDCard(t *testing.T) {
	type args struct {
		idcard string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name:"test1", args:args{""}, want:false},
		{name:"test1", args:args{"1"}, want:false},
		{name:"test1", args:args{"440102198001021231"}, want:false},
		{name:"test1", args:args{"440102198001021230"}, want:true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIDCard(tt.args.idcard); got != tt.want {
				t.Errorf("IsIDCard() = %v, want %v", got, tt.want)
			}
		})
	}
}
