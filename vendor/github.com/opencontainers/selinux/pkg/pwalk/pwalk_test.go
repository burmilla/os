package pwalk

import (
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestWalk(t *testing.T) {
	var count uint32
	concurrency := runtime.NumCPU() * 2

	dir, total, err := prepareTestSet(3, 2, 1)
	if err != nil {
		t.Fatalf("dataset creation failed: %v", err)
	}
	defer os.RemoveAll(dir)

	err = WalkN(dir,
		func(_ string, _ os.FileInfo, _ error) error {
			atomic.AddUint32(&count, 1)
			return nil
		},
		concurrency)

	if err != nil {
		t.Errorf("Walk failed: %v", err)
	}
	if count != uint32(total) {
		t.Errorf("File count mismatch: found %d, expected %d", count, total)
	}

	t.Logf("concurrency: %d, files found: %d\n", concurrency, count)
}

func TestWalkManyErrors(t *testing.T) {
	var count uint32

	dir, total, err := prepareTestSet(3, 3, 2)
	if err != nil {
		t.Fatalf("dataset creation failed: %v", err)
	}
	defer os.RemoveAll(dir)

	max := uint32(total / 2)
	e42 := errors.New("42")
	err = Walk(dir,
		func(p string, i os.FileInfo, _ error) error {
			if atomic.AddUint32(&count, 1) > max {
				return e42
			}
			return nil
		})
	t.Logf("found %d of %d files", count, total)

	if err == nil {
		t.Errorf("Walk succeeded, but error is expected")
		if count != uint32(total) {
			t.Errorf("File count mismatch: found %d, expected %d", count, total)
		}
	}
}

func makeManyDirs(prefix string, levels, dirs, files int) (count int, err error) {
	for d := 0; d < dirs; d++ {
		var dir string
		dir, err = os.MkdirTemp(prefix, "d-")
		if err != nil {
			return
		}
		count++
		for f := 0; f < files; f++ {
			var fi *os.File
			fi, err = os.CreateTemp(dir, "f-")
			if err != nil {
				return count, err
			}
			_ = fi.Close()
			count++
		}
		if levels == 0 {
			continue
		}
		var c int
		if c, err = makeManyDirs(dir, levels-1, dirs, files); err != nil {
			return
		}
		count += c
	}

	return
}

// prepareTestSet() creates a directory tree of shallow files,
// to be used for testing or benchmarking.
//
// Total dirs: dirs^levels + dirs^(levels-1) + ... + dirs^1
// Total files: total_dirs * files
func prepareTestSet(levels, dirs, files int) (dir string, total int, err error) {
	dir, err = os.MkdirTemp(".", "pwalk-test-")
	if err != nil {
		return
	}
	total, err = makeManyDirs(dir, levels, dirs, files)
	if err != nil && total > 0 {
		_ = os.RemoveAll(dir)
		dir = ""
		total = 0
		return
	}
	total++ // this dir

	return
}

type walkerFunc func(root string, walkFn WalkFunc) error

func genWalkN(n int) walkerFunc {
	return func(root string, walkFn WalkFunc) error {
		return WalkN(root, walkFn, n)
	}
}

func BenchmarkWalk(b *testing.B) {
	const (
		levels = 5 // how deep
		dirs   = 3 // dirs on each levels
		files  = 8 // files on each levels
	)

	benchmarks := []struct {
		walk filepath.WalkFunc
		name string
	}{
		{name: "Empty", walk: cbEmpty},
		{name: "ReadFile", walk: cbReadFile},
		{name: "ChownChmod", walk: cbChownChmod},
		{name: "RandomSleep", walk: cbRandomSleep},
	}

	walkers := []struct {
		walker walkerFunc
		name   string
	}{
		{name: "filepath.Walk", walker: filepath.Walk},
		{name: "pwalk.Walk", walker: Walk},
		// test WalkN with various values of N
		{name: "pwalk.Walk1", walker: genWalkN(1)},
		{name: "pwalk.Walk2", walker: genWalkN(2)},
		{name: "pwalk.Walk4", walker: genWalkN(4)},
		{name: "pwalk.Walk8", walker: genWalkN(8)},
		{name: "pwalk.Walk16", walker: genWalkN(16)},
		{name: "pwalk.Walk32", walker: genWalkN(32)},
		{name: "pwalk.Walk64", walker: genWalkN(64)},
		{name: "pwalk.Walk128", walker: genWalkN(128)},
		{name: "pwalk.Walk256", walker: genWalkN(256)},
	}

	dir, total, err := prepareTestSet(levels, dirs, files)
	if err != nil {
		b.Fatalf("dataset creation failed: %v", err)
	}
	defer os.RemoveAll(dir)
	b.Logf("dataset: %d levels x %d dirs x %d files, total entries: %d", levels, dirs, files, total)

	for _, bm := range benchmarks {
		for _, w := range walkers {
			walker := w.walker
			walkFn := bm.walk
			// preheat
			if err := w.walker(dir, bm.walk); err != nil {
				b.Errorf("walk failed: %v", err)
			}
			// benchmark
			b.Run(bm.name+"/"+w.name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					if err := walker(dir, walkFn); err != nil {
						b.Errorf("walk failed: %v", err)
					}
				}
			})
		}
	}
}

func cbEmpty(_ string, _ os.FileInfo, _ error) error {
	return nil
}

func cbChownChmod(path string, info os.FileInfo, _ error) error {
	_ = os.Chown(path, 0, 0)
	mode := os.FileMode(0o644)
	if info.Mode().IsDir() {
		mode = os.FileMode(0o755)
	}
	_ = os.Chmod(path, mode)

	return nil
}

func cbReadFile(path string, info os.FileInfo, _ error) error {
	var err error
	if info.Mode().IsRegular() {
		_, err = os.ReadFile(path)
	}
	return err
}

func cbRandomSleep(_ string, _ os.FileInfo, _ error) error {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Microsecond) //nolint:gosec // ignore G404: Use of weak random number generator
	return nil
}
