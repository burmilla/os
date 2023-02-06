/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

// Copyright 2018 CNI authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
//
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integration

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"sync"
	"syscall"
	"testing"

	"github.com/containerd/continuity/fs"
	"github.com/containerd/continuity/testutil"
	"github.com/containerd/go-cni"
	"github.com/stretchr/testify/assert"
)

var (
	baseNetNSDir = "/var/run/netns/"

	defaultCNIPluginDir = "/opt/cni/bin/"

	cniBridgePluginCfg = `
{
  "cniVersion": "1.0.0",
  "name": "gocni-test",
  "plugins": [
    {
      "type":"bridge",
      "bridge":"gocni-test0",
      "isGateway":true,
      "ipMasq":true,
      "promiscMode":true,
      "ipam":{
        "type":"host-local",
        "ranges":[
          [{
            "subnet":"10.88.0.0/16"
          }],
          [{
            "subnet":"2001:4860:4860::/64"
          }]
        ],
        "routes":[
          {"dst":"0.0.0.0/0"},
          {"dst":"::/0"}
        ]
       }
    },
    {
      "type":"portmap",
      "capabilities":{
        "portMappings":true
      }
    }
  ]
}
`

	cniBridgePluginCfgWithoutVersion = `
{
  "name": "gocni-test",
  "plugins": [
    {
      "type":"bridge",
      "bridge":"gocni-test0",
      "isGateway":true,
      "ipMasq":true,
      "promiscMode":true,
      "ipam":{
        "type":"host-local",
        "ranges":[
          [{
            "subnet":"10.88.0.0/16"
          }],
          [{
            "subnet":"2001:4860:4860::/64"
          }]
        ],
        "routes":[
          {"dst":"0.0.0.0/0"},
          {"dst":"::/0"}
        ]
       }
    },
    {
      "type":"portmap",
      "capabilities":{
        "portMappings":true
      }
    }
  ]
}
`
)

// TestBasicSetupAndRemove tests the cni.Setup/Remove with real bridge and
// loopback CNI plugins.
//
// NOTE:
//
// 1. It required that the both bridge and loopback CNI plugins are installed
// in /opt/cni/bin.
//
// 2. Since #76 enables parallel mode, we should enable -race option for this.
func TestBasicSetupAndRemove(t *testing.T) {
	testutil.RequiresRoot(t)

	// setup config dir
	tmpPluginConfDir, err := os.MkdirTemp("", t.Name()+"-conf")
	assert.NoError(t, err, "create temp dir for plugin conf dir")
	defer os.RemoveAll(tmpPluginConfDir)

	assert.NoError(t,
		ioutil.WriteFile(
			path.Join(tmpPluginConfDir, "10-gocni-test-net.conflist"),
			[]byte(cniBridgePluginCfg),
			0600,
		),
		"init cni config",
	)

	// copy plugins from /opt/cni/bin
	tmpPluginDir, err := os.MkdirTemp("", t.Name()+"-bin")
	assert.NoError(t, err, "create temp dir for plugin bin dir")
	defer os.RemoveAll(tmpPluginDir)

	assert.NoError(t,
		fs.CopyDir(tmpPluginDir, defaultCNIPluginDir),
		"copy %v into %v", defaultCNIPluginDir, tmpPluginDir)

	nsPath, done, err := createNetNS()
	assert.NoError(t, err, "create temp netns")
	defer func() {
		assert.NoError(t, done(), "cleanup temp netns")
	}()

	defaultIfName := "eth0"
	ctx := context.Background()
	id := t.Name()

	for idx, opts := range [][]cni.Opt{
		// Use default plugin dir
		{
			cni.WithMinNetworkCount(2),
			cni.WithPluginConfDir(tmpPluginConfDir),
		},
		// Use customize plugin dir
		{
			cni.WithMinNetworkCount(2),
			cni.WithPluginConfDir(tmpPluginConfDir),
			cni.WithPluginDir([]string{
				tmpPluginDir,
			}),
		},
	} {
		l, err := cni.New(opts...)
		assert.NoError(t, err, "[%v] initialize cni library", idx)

		assert.NoError(t,
			l.Load(cni.WithLoNetwork, cni.WithDefaultConf),
			"[%v] load cni configuration", idx,
		)

		// Setup network
		result, err := l.Setup(ctx, id, nsPath)
		assert.NoError(t, err, "[%v] setup network interfaces for namespace in parallel %v", idx, nsPath)

		ip := result.Interfaces[defaultIfName].IPConfigs[0].IP.String()
		t.Logf("[%v] ip is %v", idx, ip)

		assert.NoError(t,
			l.Remove(ctx, id, nsPath),
			"[%v] teardown network interfaces for namespace %v", idx, nsPath,
		)

		// Setup network serially
		result, err = l.SetupSerially(ctx, id, nsPath)
		assert.NoError(t, err, "[%v] setup network interfaces for namespace serially%v", idx, nsPath)

		ip = result.Interfaces[defaultIfName].IPConfigs[0].IP.String()
		t.Logf("[%v] ip is %v", idx, ip)

		assert.NoError(t,
			l.Remove(ctx, id, nsPath),
			"[%v] teardown network interfaces for namespace %v", idx, nsPath,
		)
	}
}

func TestBasicSetupAndRemovePluginWithoutVersion(t *testing.T) {
	testutil.RequiresRoot(t)

	// setup config dir
	tmpPluginConfDir, err := os.MkdirTemp("", t.Name()+"-conf")
	assert.NoError(t, err, "create temp dir for plugin conf dir")
	defer os.RemoveAll(tmpPluginConfDir)

	assert.NoError(t,
		ioutil.WriteFile(
			path.Join(tmpPluginConfDir, "10-gocni-test-net.conflist"),
			[]byte(cniBridgePluginCfgWithoutVersion),
			0600,
		),
		"init cni config",
	)

	// copy plugins from /opt/cni/bin
	tmpPluginDir, err := os.MkdirTemp("", t.Name()+"-bin")
	assert.NoError(t, err, "create temp dir for plugin bin dir")
	defer os.RemoveAll(tmpPluginDir)

	assert.NoError(t,
		fs.CopyDir(tmpPluginDir, defaultCNIPluginDir),
		"copy %v into %v", defaultCNIPluginDir, tmpPluginDir)

	nsPath, done, err := createNetNS()
	assert.NoError(t, err, "create temp netns")
	defer func() {
		assert.NoError(t, done(), "cleanup temp netns")
	}()

	defaultIfName := "eth0"
	ctx := context.Background()
	id := t.Name()

	for idx, opts := range [][]cni.Opt{
		// Use default plugin dir
		{
			cni.WithMinNetworkCount(2),
			cni.WithPluginConfDir(tmpPluginConfDir),
		},
		// Use customize plugin dir
		{
			cni.WithMinNetworkCount(2),
			cni.WithPluginConfDir(tmpPluginConfDir),
			cni.WithPluginDir([]string{
				tmpPluginDir,
			}),
		},
	} {
		l, err := cni.New(opts...)
		assert.NoError(t, err, "[%v] initialize cni library", idx)

		assert.NoError(t,
			l.Load(cni.WithLoNetwork, cni.WithDefaultConf),
			"[%v] load cni configuration", idx,
		)

		// Setup network
		result, err := l.Setup(ctx, id, nsPath)
		assert.NoError(t, err, "[%v] setup network interfaces for namespace in parallel %v", idx, nsPath)

		ip := result.Interfaces[defaultIfName].IPConfigs[0].IP.String()
		t.Logf("[%v] ip is %v", idx, ip)

		assert.NoError(t,
			l.Remove(ctx, id, nsPath),
			"[%v] teardown network interfaces for namespace %v", idx, nsPath,
		)

		// Setup network serially
		result, err = l.SetupSerially(ctx, id, nsPath)
		assert.NoError(t, err, "[%v] setup network interfaces for namespace serially%v", idx, nsPath)

		ip = result.Interfaces[defaultIfName].IPConfigs[0].IP.String()
		t.Logf("[%v] ip is %v", idx, ip)

		assert.NoError(t,
			l.Remove(ctx, id, nsPath),
			"[%v] teardown network interfaces for namespace %v", idx, nsPath,
		)
	}
}

// createNetNS returns temp netns path.
//
// NOTE: It is based on https://github.com/containernetworking/plugins/blob/v1.0.1/pkg/testutils/netns_linux.go.
// That can prevent from introducing unnessary dependencies in go.mod.
func createNetNS() (_ string, _ func() error, retErr error) {
	b := make([]byte, 16)
	if _, err := rand.Reader.Read(b); err != nil {
		return "", nil, fmt.Errorf("failed to generate random netns name: %w", err)
	}

	// Create the directory for mounting network namespaces
	// This needs to be a shared mountpoint in case it is mounted in to
	// other namespaces (containers)
	if err := os.MkdirAll(baseNetNSDir, 0755); err != nil {
		return "", nil, fmt.Errorf("failed to init base netns dir %s: %v", baseNetNSDir, err)
	}

	// create an empty file at the mount point
	nsName := fmt.Sprintf("gocni-test-%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	nsPath := path.Join(baseNetNSDir, nsName)
	mountPointFd, err := os.Create(nsPath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp nspath %s: %v", nsPath, err)
	}
	mountPointFd.Close()

	defer func() {
		if retErr != nil {
			_ = os.RemoveAll(nsPath)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	// do namespace work in a dedicated goroutine, so that we can safely
	// Lock/Unlock OSThread without upsetting the lock/unlock state of
	// the caller of this function
	go (func() {
		defer wg.Done()

		// Don't unlock. By not unlocking, golang will kill the OS thread
		// when the goroutine is done (>= go1.10). Since <= go1.10 has
		// been deprecated, we don't need to get current net ns and
		// reset.
		runtime.LockOSThread()

		// create a new netns on the current thread
		if err = syscall.Unshare(syscall.CLONE_NEWNET); err != nil {
			return
		}

		// bind mount the netns from the current thread (from /proc) onto the
		// mount point. This causes the namespace to persist, even when there
		// are no threads in the ns.
		err = syscall.Mount(getCurrentThreadNetNSPath(), nsPath, "none", syscall.MS_BIND, "")
		if err != nil {
			err = fmt.Errorf("failed to bind mount ns at %s: %w", nsPath, err)
		}
	})()
	wg.Wait()

	if err != nil {
		return "", nil, fmt.Errorf("failed to create net namespace: %w", err)
	}

	return nsPath, func() error {
		if err := syscall.Unmount(nsPath, 0); err != nil {
			return fmt.Errorf("failed to unmount netns: at %s: %v", nsPath, err)
		}

		if err := os.Remove(nsPath); err != nil {
			return fmt.Errorf("failed to remove nspath %s: %v", nsPath, err)
		}
		return nil
	}, nil
}

// getCurrentThreadNetNSPath copied from pkg/ns
//
// NOTE: It is from https://github.com/containernetworking/plugins/blob/v1.0.1/pkg/testutils/netns_linux.go.
func getCurrentThreadNetNSPath() string {
	// /proc/self/ns/net returns the namespace of the main thread, not
	// of whatever thread this goroutine is running on.  Make sure we
	// use the thread's net namespace since the thread is switching around
	return fmt.Sprintf("/proc/%d/task/%d/ns/net", os.Getpid(), syscall.Gettid())
}
