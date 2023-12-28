package testutil

import "testing"

func TestRandomString(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "Positive length", args: args{length: 10}, want: 10},
		{name: "Negative length", args: args{length: -5}, want: 0},
		{name: "Zero length", args: args{length: 0}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomString(tt.args.length); len(got) != tt.want {
				t.Errorf("RandomString() = %v, want %v", got, tt.want)
			}
		})
	}
}
