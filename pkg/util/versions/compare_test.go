package versions

import (
	"testing"
)

func assertVersion(t *testing.T, a, b string, result int) {
	if r := compare(a, b); r != result {
		t.Fatalf("Unexpected version comparison result. Found %d, expected %d", r, result)
	}
}

func TestCompareVersion(t *testing.T) {
	assertVersion(t, "1.12", "1.12", 0)
	assertVersion(t, "1.0.0", "1", 0)
	assertVersion(t, "1", "1.0.0", 0)
	assertVersion(t, "1.05.00.0156", "1.0.221.9289", 1)
	assertVersion(t, "1", "1.0.1", -1)
	assertVersion(t, "1.0.1", "1", 1)
	assertVersion(t, "1.0.1", "1.0.2", -1)
	assertVersion(t, "1.0.2", "1.0.3", -1)
	assertVersion(t, "1.0.3", "1.1", -1)
	assertVersion(t, "1.1", "1.1.1", -1)
	assertVersion(t, "1.1.1", "1.1.2", -1)
	assertVersion(t, "1.1.2", "1.2", -1)
}

func TestSortEngines(t *testing.T) {
	engines := []string{"k3s-v1.34.1-k3s2", "docker-29.9.2", "qemu-10.1.0", "docker-29.10.0", "k3s-v1.34.1-k3s1", "docker-25.0.14", "k3s-v1.33.5-k3s1", "docker-29.6.1", "qemu-9.2.4"}
	SortEngines(engines)
	want := []string{"docker-25.0.14", "docker-29.6.1", "docker-29.9.2", "docker-29.10.0", "k3s-v1.33.5-k3s1", "k3s-v1.34.1-k3s1", "k3s-v1.34.1-k3s2", "qemu-9.2.4", "qemu-10.1.0"}
	for i := range want {
		if engines[i] != want[i] {
			t.Fatalf("SortEngines = %v, want %v", engines, want)
		}
	}
}
