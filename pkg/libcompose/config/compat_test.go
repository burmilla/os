package config

import (
	"testing"

	composetypes "github.com/compose-spec/compose-go/types"
)

func TestServiceConfigV1ToServiceConfig_Nil(t *testing.T) {
	if sc := ServiceConfigV1ToServiceConfig(nil); sc != nil {
		t.Errorf("expected nil, got %v", sc)
	}
}

func TestServiceConfigV1ToServiceConfig_Basic(t *testing.T) {
	v1 := &ServiceConfigV1{
		Image:         "busybox:latest",
		ContainerName: "test",
		Privileged:    true,
		Net:           "host",
		Pid:           "host",
		Restart:       "always",
	}
	sc := ServiceConfigV1ToServiceConfig(v1)
	if sc == nil {
		t.Fatal("expected non-nil ServiceConfig")
	}
	if sc.Image != "busybox:latest" {
		t.Errorf("Image = %q, want %q", sc.Image, "busybox:latest")
	}
	if sc.ContainerName != "test" {
		t.Errorf("ContainerName = %q, want %q", sc.ContainerName, "test")
	}
	if !sc.Privileged {
		t.Error("Privileged = false, want true")
	}
	if sc.NetworkMode != "host" {
		t.Errorf("NetworkMode = %q, want %q", sc.NetworkMode, "host")
	}
	if sc.Pid != "host" {
		t.Errorf("Pid = %q, want %q", sc.Pid, "host")
	}
	if sc.Restart != "always" {
		t.Errorf("Restart = %q, want %q", sc.Restart, "always")
	}
}

func TestServiceConfigV1ToServiceConfig_Logging(t *testing.T) {
	v1 := &ServiceConfigV1{
		Image:     "busybox",
		LogDriver: "syslog",
		LogOpt:    map[string]string{"tag": "myapp"},
	}
	sc := ServiceConfigV1ToServiceConfig(v1)
	if sc.Logging == nil {
		t.Fatal("Logging should not be nil when LogDriver is set")
	}
	if sc.Logging.Driver != "syslog" {
		t.Errorf("Logging.Driver = %q, want %q", sc.Logging.Driver, "syslog")
	}
	if sc.Logging.Options["tag"] != "myapp" {
		t.Errorf("Logging.Options[tag] = %q, want %q", sc.Logging.Options["tag"], "myapp")
	}
}

func TestServiceConfigV1ToServiceConfig_NoLogging(t *testing.T) {
	v1 := &ServiceConfigV1{Image: "busybox"}
	sc := ServiceConfigV1ToServiceConfig(v1)
	if sc.Logging != nil {
		t.Errorf("Logging should be nil when LogDriver is empty, got %v", sc.Logging)
	}
}

func TestServiceConfigV1ToServiceConfig_Build(t *testing.T) {
	v1 := &ServiceConfigV1{
		Build:      "./app",
		Dockerfile: "Dockerfile.dev",
	}
	sc := ServiceConfigV1ToServiceConfig(v1)
	if sc.Build == nil {
		t.Fatal("Build should not be nil")
	}
	if sc.Build.Context != "./app" {
		t.Errorf("Build.Context = %q, want %q", sc.Build.Context, "./app")
	}
	if sc.Build.Dockerfile != "Dockerfile.dev" {
		t.Errorf("Build.Dockerfile = %q, want %q", sc.Build.Dockerfile, "Dockerfile.dev")
	}
}

func TestServiceConfigV1ToServiceConfig_Ports(t *testing.T) {
	v1 := &ServiceConfigV1{
		Image: "nginx",
		Ports: []string{"80", "8080:80", "127.0.0.1:443:443/tcp", "9090:9090/udp"},
	}
	sc := ServiceConfigV1ToServiceConfig(v1)

	if len(sc.Ports) != 4 {
		t.Fatalf("len(Ports) = %d, want 4", len(sc.Ports))
	}

	// "80"
	if sc.Ports[0].Target != 80 || sc.Ports[0].Published != "" {
		t.Errorf("port 0: Target=%d Published=%q", sc.Ports[0].Target, sc.Ports[0].Published)
	}
	// "8080:80"
	if sc.Ports[1].Target != 80 || sc.Ports[1].Published != "8080" {
		t.Errorf("port 1: Target=%d Published=%q", sc.Ports[1].Target, sc.Ports[1].Published)
	}
	// "127.0.0.1:443:443/tcp"
	if sc.Ports[2].Target != 443 || sc.Ports[2].HostIP != "127.0.0.1" || sc.Ports[2].Protocol != "tcp" {
		t.Errorf("port 2: %+v", sc.Ports[2])
	}
	// "9090:9090/udp"
	if sc.Ports[3].Protocol != "udp" {
		t.Errorf("port 3: Protocol=%q, want udp", sc.Ports[3].Protocol)
	}
}

func TestServiceConfigV1ToServiceConfig_Volumes(t *testing.T) {
	v1 := &ServiceConfigV1{
		Image:   "app",
		Volumes: []string{"/data", "/host:/container", "/host:/container:ro", "named-vol"},
	}
	sc := ServiceConfigV1ToServiceConfig(v1)

	if len(sc.Volumes) != 4 {
		t.Fatalf("len(Volumes) = %d, want 4", len(sc.Volumes))
	}
	// "/data"
	if sc.Volumes[0].Target != "/data" || sc.Volumes[0].Type != composetypes.VolumeTypeBind {
		t.Errorf("vol 0: %+v", sc.Volumes[0])
	}
	// "/host:/container"
	if sc.Volumes[1].Source != "/host" || sc.Volumes[1].Target != "/container" {
		t.Errorf("vol 1: %+v", sc.Volumes[1])
	}
	// "/host:/container:ro"
	if !sc.Volumes[2].ReadOnly {
		t.Error("vol 2: expected ReadOnly=true")
	}
	// "named-vol"
	if sc.Volumes[3].Type != composetypes.VolumeTypeVolume {
		t.Errorf("vol 3: Type=%q, want %q", sc.Volumes[3].Type, composetypes.VolumeTypeVolume)
	}
}

func TestServiceConfigV1ToServiceConfig_Environment(t *testing.T) {
	v1 := &ServiceConfigV1{
		Image:       "app",
		Environment: []string{"FOO=bar", "EMPTY=", "BARE"},
	}
	sc := ServiceConfigV1ToServiceConfig(v1)
	if sc.Environment == nil {
		t.Fatal("Environment should not be nil")
	}
	if v, ok := sc.Environment["FOO"]; !ok || v == nil || *v != "bar" {
		t.Errorf("FOO = %v", v)
	}
}

func TestConvertV1Services_SetsName(t *testing.T) {
	v1s := map[string]*ServiceConfigV1{
		"web": {Image: "nginx"},
		"db":  {Image: "postgres"},
	}
	result := ConvertV1Services(v1s)
	if len(result) != 2 {
		t.Fatalf("len = %d, want 2", len(result))
	}
	if result["web"].Name != "web" {
		t.Errorf("web.Name = %q", result["web"].Name)
	}
	if result["db"].Name != "db" {
		t.Errorf("db.Name = %q", result["db"].Name)
	}
}

func TestConvertV1Services_SkipsNil(t *testing.T) {
	v1s := map[string]*ServiceConfigV1{
		"good": {Image: "nginx"},
		"bad":  nil,
	}
	result := ConvertV1Services(v1s)
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1 (nil should be skipped)", len(result))
	}
	if _, ok := result["good"]; !ok {
		t.Error("expected 'good' in result")
	}
}

func TestVolumesToStringSlice(t *testing.T) {
	vols := []composetypes.ServiceVolumeConfig{
		{Source: "/host", Target: "/container"},
		{Source: "/host", Target: "/container", ReadOnly: true},
	}
	result := VolumesToStringSlice(vols)
	if len(result) != 2 {
		t.Fatalf("len = %d, want 2", len(result))
	}
	if result[0] != "/host:/container" {
		t.Errorf("[0] = %q", result[0])
	}
	if result[1] != "/host:/container:ro" {
		t.Errorf("[1] = %q", result[1])
	}
}

func TestPortsToStringSlice(t *testing.T) {
	ports := []composetypes.ServicePortConfig{
		{Target: 80, Protocol: "tcp"},
		{Target: 443, Published: "443", HostIP: "0.0.0.0", Protocol: "tcp"},
		{Target: 53, Published: "53", Protocol: "udp"},
	}
	result := PortsToStringSlice(ports)
	if len(result) != 3 {
		t.Fatalf("len = %d, want 3", len(result))
	}
	if result[0] != "80" {
		t.Errorf("[0] = %q, want %q", result[0], "80")
	}
	if result[1] != "0.0.0.0:443:443" {
		t.Errorf("[1] = %q, want %q", result[1], "0.0.0.0:443:443")
	}
	if result[2] != "53:53/udp" {
		t.Errorf("[2] = %q, want %q", result[2], "53:53/udp")
	}
}

func TestEnvironmentToSlice(t *testing.T) {
	val := "bar"
	env := composetypes.MappingWithEquals{
		"FOO": &val,
		"BAR": nil,
	}
	result := EnvironmentToSlice(env)
	found := map[string]bool{}
	for _, r := range result {
		found[r] = true
	}
	if !found["FOO=bar"] {
		t.Error("expected FOO=bar")
	}
	if !found["BAR"] {
		t.Error("expected bare BAR")
	}
}

func TestLoggingDriver_Nil(t *testing.T) {
	sc := &composetypes.ServiceConfig{LogDriver: "json-file"}
	if d := LoggingDriver(sc); d != "json-file" {
		t.Errorf("LoggingDriver = %q, want %q", d, "json-file")
	}
}

func TestLoggingDriver_NonNil(t *testing.T) {
	sc := &composetypes.ServiceConfig{
		Logging: &composetypes.LoggingConfig{Driver: "syslog"},
	}
	if d := LoggingDriver(sc); d != "syslog" {
		t.Errorf("LoggingDriver = %q, want %q", d, "syslog")
	}
}

func TestLoggingOptions_Nil(t *testing.T) {
	opts := map[string]string{"tag": "test"}
	sc := &composetypes.ServiceConfig{LogOpt: opts}
	if o := LoggingOptions(sc); o["tag"] != "test" {
		t.Errorf("LoggingOptions = %v", o)
	}
}

func TestLoggingOptions_NonNil(t *testing.T) {
	opts := map[string]string{"tag": "test"}
	sc := &composetypes.ServiceConfig{
		Logging: &composetypes.LoggingConfig{Options: opts},
	}
	if o := LoggingOptions(sc); o["tag"] != "test" {
		t.Errorf("LoggingOptions = %v", o)
	}
}

func TestValidateV1ServiceConfig_Valid(t *testing.T) {
	v1 := &ServiceConfigV1{Image: "busybox", Ports: []string{"80", "8080:80/tcp"}}
	if err := ValidateV1ServiceConfig("test", v1); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateV1ServiceConfig_Nil(t *testing.T) {
	if err := ValidateV1ServiceConfig("test", nil); err == nil {
		t.Error("expected error for nil config")
	}
}

func TestValidateV1ServiceConfig_EmptyVolume(t *testing.T) {
	v1 := &ServiceConfigV1{Image: "busybox", Volumes: []string{"/data", ""}}
	if err := ValidateV1ServiceConfig("test", v1); err == nil {
		t.Error("expected error for empty volume")
	}
}

func TestValidateV1ServiceConfig_BadProtocol(t *testing.T) {
	v1 := &ServiceConfigV1{Image: "busybox", Ports: []string{"80/sctp"}}
	if err := ValidateV1ServiceConfig("test", v1); err == nil {
		t.Error("expected error for unsupported protocol")
	}
}
