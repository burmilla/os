package compose

import (
	"testing"

	"github.com/burmilla/os/config"
	composeConfig "github.com/burmilla/os/pkg/libcompose/config"
)

func TestResolveEngine(t *testing.T) {
	available := []string{"docker-25.0.14", "docker-29.10.0", "docker-29.6.1", "docker-29.9.2"}

	tests := []struct {
		name      string
		engine    string
		available []string
		want      string
		wantExact bool
	}{
		{"available", "docker-29.6.1", available, "docker-29.6.1", true},
		{"missing falls back to latest of its type", "docker-29.6.2", available, "docker-29.10.0", false},
		{"index not readable", "docker-29.6.2", nil, "docker-29.6.2", true},
		{"type resolves to latest available", "docker", available, "docker-29.10.0", false},
		{"type without index", "docker", nil, "docker", true},
		{"type without engines of that type", "docker", []string{"other-1.0"}, "docker", true},
		{"missing type falls back to latest docker", "other-1.0", available, "docker-29.10.0", false},
		{"missing type without docker falls back to latest of any type", "other-1.0", []string{"foo-1.0", "foo-2.0"}, "foo-2.0", false},
		{"local file", "/var/lib/rancher/conf/docker.yml", available, "/var/lib/rancher/conf/docker.yml", true},
		{"url", "https://example.com/docker.yml", available, "https://example.com/docker.yml", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, exact := resolveEngine(tt.engine, tt.available)
			if got != tt.want || exact != tt.wantExact {
				t.Errorf("resolveEngine(%q, %v) = %q, %v; want %q, %v", tt.engine, tt.available, got, exact, tt.want, tt.wantExact)
			}
		})
	}
}

func TestEngineServiceName(t *testing.T) {
	name, err := engineServiceName("docker-29.9.2", map[interface{}]interface{}{"docker": map[interface{}]interface{}{}})
	if err != nil || name != "docker" {
		t.Fatalf("engineServiceName = %q, %v; want docker", name, err)
	}
	if _, err := engineServiceName("broken", map[interface{}]interface{}{"a": nil, "b": nil}); err == nil {
		t.Fatal("engineServiceName accepted two services")
	}
	if _, err := engineServiceName("broken", map[interface{}]interface{}{}); err == nil {
		t.Fatal("engineServiceName accepted no services")
	}
}

func TestEngineServices(t *testing.T) {
	system := map[string]string{config.ScopeLabel: config.System}
	services := map[string]*composeConfig.ServiceConfigV1{
		"docker":  {Labels: system},
		"console": {Labels: system},
		"user":    {},
	}

	cfg := &config.CloudConfig{}
	got := engineServices(cfg, services)
	if len(got) != 3 {
		t.Fatalf("docker engine keeps all services, got %v", got)
	}

	cfg.Rancher.Engine = "other-1.0"
	got = engineServices(cfg, services)
	if _, ok := got["docker"]; ok {
		t.Error("docker service kept with another engine")
	}
	if _, ok := got["user"]; ok {
		t.Error("user Docker service kept with another engine")
	}
	if _, ok := got["console"]; !ok {
		t.Error("console service dropped")
	}
}

func TestUserService(t *testing.T) {
	system := map[interface{}]interface{}{"svc": map[interface{}]interface{}{"labels": map[interface{}]interface{}{config.ScopeLabel: config.System}}}
	user := map[interface{}]interface{}{"svc": map[interface{}]interface{}{"image": "x"}}
	if userService(system) {
		t.Error("system service reported as user service")
	}
	if !userService(user) {
		t.Error("user service not detected")
	}
}
