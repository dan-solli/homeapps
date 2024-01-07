package testutil

import "testing"

func TestGetRandomPort(t *testing.T) {
	tests := []struct {
		name    string
		GTPort  int
		wantErr bool
	}{
		// TODO: Add test cases.
		{"TestGetRandomPort", 1024, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPort, err := GetRandomPort()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRandomPort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPort <= tt.GTPort {
				t.Errorf("GetRandomPort() = %v, want at least %v", gotPort, tt.GTPort)
			}
		})
	}
}
