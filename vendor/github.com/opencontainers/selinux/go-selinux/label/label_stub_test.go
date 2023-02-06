//go:build !linux
// +build !linux

package label

import "testing"

const testLabel = "system_u:object_r:container_file_t:s0:c1,c2"

func TestInit(t *testing.T) {
	var testNull []string
	_, _, err := InitLabels(testNull)
	if err != nil {
		t.Log("InitLabels Failed")
		t.Fatal(err)
	}
	testDisabled := []string{"disable"}
	roMountLabel := ROMountLabel()
	if roMountLabel != "" {
		t.Errorf("ROMountLabel Failed")
	}
	plabel, mlabel, err := InitLabels(testDisabled)
	if err != nil {
		t.Log("InitLabels Disabled Failed")
		t.Fatal(err)
	}
	if plabel != "" {
		t.Fatal("InitLabels Disabled Failed")
	}
	if mlabel != "" {
		t.Fatal("InitLabels Disabled mlabel Failed")
	}
	testUser := []string{"user:user_u", "role:user_r", "type:user_t", "level:s0:c1,c15"}
	_, _, err = InitLabels(testUser)
	if err != nil {
		t.Log("InitLabels User Failed")
		t.Fatal(err)
	}
}

func TestRelabel(t *testing.T) {
	if err := Relabel("/etc", testLabel, false); err != nil {
		t.Fatalf("Relabel /etc succeeded")
	}
}

func TestSocketLabel(t *testing.T) {
	label := testLabel
	if err := SetSocketLabel(label); err != nil {
		t.Fatal(err)
	}
	if _, err := SocketLabel(); err != nil {
		t.Fatal(err)
	}
}

func TestKeyLabel(t *testing.T) {
	label := testLabel
	if err := SetKeyLabel(label); err != nil {
		t.Fatal(err)
	}
	if _, err := KeyLabel(); err != nil {
		t.Fatal(err)
	}
}

func TestProcessLabel(t *testing.T) {
	label := testLabel
	if err := SetProcessLabel(label); err != nil {
		t.Fatal(err)
	}
	if _, err := ProcessLabel(); err != nil {
		t.Fatal(err)
	}
}

func TestCheckLabelCompile(t *testing.T) {
	if _, _, err := GenLabels(""); err != nil {
		t.Fatal(err)
	}

	tmpDir := t.TempDir()
	if _, err := FileLabel(tmpDir); err != nil {
		t.Fatal(err)
	}

	if err := SetFileLabel(tmpDir, "foobar"); err != nil {
		t.Fatal(err)
	}

	if err := SetFileCreateLabel("foobar"); err != nil {
		t.Fatal(err)
	}

	if _, err := PidLabel(0); err != nil {
		t.Fatal(err)
	}

	ClearLabels()

	if err := ReserveLabel("foobar"); err != nil {
		t.Fatal(err)
	}

	if err := ReleaseLabel("foobar"); err != nil {
		t.Fatal(err)
	}

	_, _ = DupSecOpt("foobar")
	DisableSecOpt()

	if err := Validate("foobar"); err != nil {
		t.Fatal(err)
	}
	if relabel := RelabelNeeded("foobar"); relabel {
		t.Fatal("Relabel failed")
	}
	if shared := IsShared("foobar"); shared {
		t.Fatal("isshared failed")
	}
}
