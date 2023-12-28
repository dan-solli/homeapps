package main_test

import (
	"testing"

	tu "github.com/dan-solli/homeapps/common/testutil"
	c "github.com/dan-solli/homeapps/microservice/checklist"
)

func TestAddTodo(t *testing.T) {

	str := tu.RandomString(10)

	type args struct {
		todo string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: str,
			args: args{str},
			want: str,
		},
	}
	var c c.CheckList
	c.New("testlist")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.AddTodo(tt.args.todo); got != tt.want {
				t.Errorf("AddTodo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountTodo(t *testing.T) {
	str := tu.RandomString(10)

	type args struct {
		todo string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: str,
			args: args{str},
			want: 1,
		},
	}
	var c c.CheckList
	c.New("testlist")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.AddTodo(tt.args.todo)
			if got := c.CountTodo(); got != tt.want {
				t.Errorf("CountTodo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckListItem(t *testing.T) {

	cli := c.NewCheckListItem().
		WithTitle("Test").
		WithGeolocation(nil).
		WithStorage(nil).
		WithPriority(c.LOW).
		WithRaci(nil).
		WithRecurrence("Every monday").
		WithBlock(c.BLOCKED, "11223344-5566-7788-9900-aabbccddeeff")
	if cli.GetTitle() != "Test" {
		t.Error("Expected Test, got ", cli.GetTitle())
	}
	_, err := cli.Store()
	if err == nil {
		t.Error("Expected error, got nil")
	}
	_, err = cli.Fetch("11223344-5566-7788-9900-aabbccddeeff")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
