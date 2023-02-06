//go:build !windows
// +build !windows

package mountinfo

import "testing"

func TestGetMounts(t *testing.T) {
	mounts, err := GetMounts(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(mounts) < 2 {
		t.Fatalf("should have at least two mounts, got %d: %+v", len(mounts), mounts)
	}
}
