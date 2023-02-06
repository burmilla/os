package mountinfo

import "testing"

func TestPrefixFilter(t *testing.T) {
	tests := []struct {
		prefix     string
		mountPoint string
		shouldSkip bool
	}{
		{prefix: "/a", mountPoint: "/a", shouldSkip: false},
		{prefix: "/a", mountPoint: "/a/b", shouldSkip: false},
		{prefix: "/a", mountPoint: "/aa", shouldSkip: true},
		{prefix: "/a", mountPoint: "/aa/b", shouldSkip: true},

		// invalid prefix: prefix path must be cleaned and have no trailing slash
		{prefix: "/a/", mountPoint: "/a", shouldSkip: true},
		{prefix: "/a/", mountPoint: "/a/b", shouldSkip: true},
	}
	for _, tc := range tests {
		filter := PrefixFilter(tc.prefix)
		skip, _ := filter(&Info{Mountpoint: tc.mountPoint})
		if skip != tc.shouldSkip {
			if tc.shouldSkip {
				t.Errorf("prefix %q: expected %q to be skipped", tc.prefix, tc.mountPoint)
			} else {
				t.Errorf("prefix %q: expected %q not to be skipped", tc.prefix, tc.mountPoint)
			}
		}
	}
}
