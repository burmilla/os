package control

import (
	"reflect"
	"strings"
	"testing"

	"github.com/burmilla/os/config"
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

func TestQemuScripts(t *testing.T) {
	cfg := &config.CloudConfig{}
	scripts := qemuScripts(cfg)
	if len(scripts) != len(qemuDefaultTools) {
		t.Fatalf("default tools: got %d scripts, want %d", len(scripts), len(qemuDefaultTools))
	}
	virsh, ok := scripts["/usr/bin/virsh"]
	if !ok || !strings.Contains(virsh, "system-docker exec -i $tty qemu virsh") {
		t.Errorf("unexpected virsh wrapper: %q", virsh)
	}

	cfg.Rancher.Qemu.Tools = []string{"virt-top", "../evil", ""}
	scripts = qemuScripts(cfg)
	if _, ok := scripts["/usr/bin/virt-top"]; !ok || len(scripts) != 2 {
		t.Errorf("configured tools: %v", scripts)
	}
	if _, ok := scripts["/usr/bin/evil"]; !ok {
		t.Errorf("path is reduced to its base name: %v", scripts)
	}
}
