package label

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/opencontainers/selinux/go-selinux"
)

func needSELinux(t *testing.T) {
	t.Helper()
	if !selinux.GetEnabled() {
		t.Skip("SELinux not enabled, skipping.")
	}
}

func TestInit(t *testing.T) {
	needSELinux(t)

	var testNull []string
	_, _, err := InitLabels(testNull)
	if err != nil {
		t.Fatalf("InitLabels failed: %v:", err)
	}
	testDisabled := []string{"disable"}
	roMountLabel := ROMountLabel()
	if roMountLabel == "" {
		t.Fatal("ROMountLabel: empty")
	}
	plabel, mlabel, err := InitLabels(testDisabled)
	if err != nil {
		t.Fatalf("InitLabels(disabled) failed: %v", err)
	}
	if plabel != "" {
		t.Fatalf("InitLabels(disabled): %q not empty", plabel)
	}
	if mlabel != "system_u:object_r:container_file_t:s0:c1022,c1023" {
		t.Fatalf("InitLabels Disabled mlabel Failed, %s", mlabel)
	}

	testUser := []string{"user:user_u", "role:user_r", "type:user_t", "level:s0:c1,c15"}
	plabel, mlabel, err = InitLabels(testUser)
	if err != nil {
		t.Fatalf("InitLabels(user) failed: %v", err)
	}
	if plabel != "user_u:user_r:user_t:s0:c1,c15" || (mlabel != "user_u:object_r:container_file_t:s0:c1,c15" && mlabel != "user_u:object_r:svirt_sandbox_file_t:s0:c1,c15") {
		t.Fatalf("InitLabels(user) failed (plabel=%q, mlabel=%q)", plabel, mlabel)
	}

	testBadData := []string{"user", "role:user_r", "type:user_t", "level:s0:c1,c15"}
	if _, _, err = InitLabels(testBadData); err == nil {
		t.Fatal("InitLabels(bad): expected error, got nil")
	}
}

func TestDuplicateLabel(t *testing.T) {
	secopt, err := DupSecOpt("system_u:system_r:container_t:s0:c1,c2")
	if err != nil {
		t.Fatalf("DupSecOpt: %v", err)
	}
	for _, opt := range secopt {
		con := strings.SplitN(opt, ":", 2)
		if con[0] == "user" {
			if con[1] != "system_u" {
				t.Errorf("DupSecOpt Failed user incorrect")
			}
			continue
		}
		if con[0] == "role" {
			if con[1] != "system_r" {
				t.Errorf("DupSecOpt Failed role incorrect")
			}
			continue
		}
		if con[0] == "type" {
			if con[1] != "container_t" {
				t.Errorf("DupSecOpt Failed type incorrect")
			}
			continue
		}
		if con[0] == "level" {
			if con[1] != "s0:c1,c2" {
				t.Errorf("DupSecOpt Failed level incorrect")
			}
			continue
		}
		t.Errorf("DupSecOpt failed: invalid field %q", con[0])
	}
	secopt = DisableSecOpt()
	if secopt[0] != "disable" {
		t.Errorf("DisableSecOpt failed: expected \"disable\", got %q", secopt[0])
	}
}

func TestRelabel(t *testing.T) {
	needSELinux(t)

	testdir := t.TempDir()
	label := "system_u:object_r:container_file_t:s0:c1,c2"
	if err := Relabel(testdir, "", true); err != nil {
		t.Fatalf("Relabel with no label failed: %v", err)
	}
	if err := Relabel(testdir, label, true); err != nil {
		t.Fatalf("Relabel shared failed: %v", err)
	}
	if err := Relabel(testdir, label, false); err != nil {
		t.Fatalf("Relabel unshared failed: %v", err)
	}
	if err := Relabel("/etc", label, false); err == nil {
		t.Fatalf("Relabel /etc succeeded")
	}
	if err := Relabel("/", label, false); err == nil {
		t.Fatalf("Relabel / succeeded")
	}
	if err := Relabel("/usr", label, false); err == nil {
		t.Fatalf("Relabel /usr succeeded")
	}
	if err := Relabel("/usr/", label, false); err == nil {
		t.Fatalf("Relabel /usr/ succeeded")
	}
	if err := Relabel("/etc/passwd", label, false); err == nil {
		t.Fatalf("Relabel /etc/passwd succeeded")
	}
	if home := os.Getenv("HOME"); home != "" {
		if err := Relabel(home, label, false); err == nil {
			t.Fatalf("Relabel %s succeeded", home)
		}
	}
}

func TestValidate(t *testing.T) {
	if err := Validate("zZ"); !errors.Is(err, ErrIncompatibleLabel) {
		t.Fatalf("Expected incompatible error, got %v", err)
	}
	if err := Validate("Z"); err != nil {
		t.Fatal(err)
	}
	if err := Validate("z"); err != nil {
		t.Fatal(err)
	}
	if err := Validate(""); err != nil {
		t.Fatal(err)
	}
}

func TestIsShared(t *testing.T) {
	if shared := IsShared("Z"); shared {
		t.Fatalf("Expected label `Z` to not be shared, got %v", shared)
	}
	if shared := IsShared("z"); !shared {
		t.Fatalf("Expected label `z` to be shared, got %v", shared)
	}
	if shared := IsShared("Zz"); !shared {
		t.Fatalf("Expected label `Zz` to be shared, got %v", shared)
	}
}

func TestSELinuxNoLevel(t *testing.T) {
	needSELinux(t)

	tlabel := "system_u:system_r:container_t"
	dup, err := DupSecOpt(tlabel)
	if err != nil {
		t.Fatal(err)
	}

	if len(dup) != 3 {
		t.Errorf("DupSecOpt failed on non mls label: expected 3, got %d", len(dup))
	}
	con, err := selinux.NewContext(tlabel)
	if err != nil {
		t.Fatal(err)
	}
	if con.Get() != tlabel {
		t.Errorf("NewContaxt and con.Get() failed on non mls label: expected %q, got %q", tlabel, con.Get())
	}
}

func TestSocketLabel(t *testing.T) {
	needSELinux(t)

	label := "system_u:object_r:container_t:s0:c1,c2"
	if err := selinux.SetSocketLabel(label); err != nil {
		t.Fatal(err)
	}
	nlabel, err := selinux.SocketLabel()
	if err != nil {
		t.Fatal(err)
	}
	if label != nlabel {
		t.Errorf("SocketLabel %s != %s", nlabel, label)
	}
}

func TestKeyLabel(t *testing.T) {
	needSELinux(t)

	label := "system_u:object_r:container_t:s0:c1,c2"
	if err := selinux.SetKeyLabel(label); err != nil {
		t.Fatal(err)
	}
	nlabel, err := selinux.KeyLabel()
	if err != nil {
		t.Fatal(err)
	}
	if label != nlabel {
		t.Errorf("KeyLabel %s != %s", nlabel, label)
	}
}

func TestFileLabel(t *testing.T) {
	needSELinux(t)

	testUser := []string{"filetype:test_file_t", "level:s0:c1,c15"}
	_, mlabel, err := InitLabels(testUser)
	if err != nil {
		t.Fatalf("InitLabels(user) failed: %v", err)
	}
	if mlabel != "system_u:object_r:test_file_t:s0:c1,c15" {
		t.Fatalf("InitLabels(filetype) failed: %v", err)
	}
}
