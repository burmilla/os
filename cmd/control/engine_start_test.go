package control

import (
	"reflect"
	"testing"
)

func TestEngineLogCommand(t *testing.T) {
	got := engineLogCommand("/var/log/k3s.log", []string{"/usr/bin/k3s", "server", "--node-name", "it's"})
	want := []string{"bash", "-c", `exec '/usr/bin/k3s' 'server' '--node-name' 'it'\''s' >> /var/log/k3s.log 2>&1`}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("engineLogCommand = %q, want %q", got, want)
	}
}

func TestEngineStarterFor(t *testing.T) {
	if _, err := engineStarterFor(nil); err == nil {
		t.Error("no engine type accepted")
	}
	if _, err := engineStarterFor([]string{"unknown"}); err == nil {
		t.Error("unknown engine type accepted")
	}
	if _, err := engineStarterFor([]string{"k3s"}); err != nil {
		t.Errorf("k3s rejected: %v", err)
	}
}
