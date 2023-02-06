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

package cni

import (
	"context"
	"net"
	"testing"

	cnilibrary "github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/types"
	types020 "github.com/containernetworking/cni/pkg/types/020"
	types040 "github.com/containernetworking/cni/pkg/types/040"
	types100 "github.com/containernetworking/cni/pkg/types/100"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestLibCNIType020 tests the cni version 0.2.0 plugin
// config and parses the result into structured data
func TestLibCNIType020(t *testing.T) {
	// Get the default CNI config
	l := defaultCNIConfig()

	// Create a fake cni config directory and file
	cniDir, confDir := makeFakeCNIConfig(t)
	defer tearDownCNIConfig(t, cniDir)
	l.pluginConfDir = confDir
	// Set the minimum network count as 2 for this test
	l.networkCount = 2
	err := l.Load(WithAllConf)
	assert.NoError(t, err)

	err = l.Status()
	assert.NoError(t, err)

	mockCNI := &MockCNI{}
	l.networks[0].cni = mockCNI
	expectedRT := &cnilibrary.RuntimeConf{ContainerID: "container-id1",
		NetNS:          "/proc/12345/ns/net",
		IfName:         "eth0",
		Args:           [][2]string(nil),
		CapabilityArgs: map[string]interface{}{},
	}
	mockCNI.On("AddNetworkList", l.networks[0].config, expectedRT).Return(&types020.Result{
		CNIVersion: "0.2.0",
		IP4: &types020.IPConfig{
			IP: net.IPNet{
				IP: []byte{10, 0, 0, 1},
			},
		},
	}, nil)
	mockCNI.On("DelNetworkList", l.networks[0].config, expectedRT).Return(nil)

	l.networks[1].cni = mockCNI
	expectedRT = &cnilibrary.RuntimeConf{ContainerID: "container-id1",
		NetNS:          "/proc/12345/ns/net",
		IfName:         "eth1",
		Args:           [][2]string(nil),
		CapabilityArgs: map[string]interface{}{},
	}
	mockCNI.On("AddNetworkList", l.networks[1].config, expectedRT).Return(&types020.Result{
		CNIVersion: "0.2.0",
		IP4: &types020.IPConfig{
			IP: net.IPNet{
				IP: []byte{10, 0, 0, 2},
			},
		},
	}, nil)
	mockCNI.On("DelNetworkList", l.networks[1].config, expectedRT).Return(nil)

	ctx := context.Background()

	r, err := l.Setup(ctx, "container-id1", "/proc/12345/ns/net")
	assert.NoError(t, err)
	assert.Contains(t, r.Interfaces, "eth0")
	assert.NotNil(t, r.Interfaces["eth0"].IPConfigs)
	assert.Equal(t, r.Interfaces["eth0"].IPConfigs[0].IP.String(), "10.0.0.1")
	assert.Equal(t, r.Interfaces["eth0"].IPConfigs[1].IP.String(), "10.0.0.2")
	err = l.Remove(ctx, "container-id1", "/proc/12345/ns/net")
	assert.NoError(t, err)

	c := l.GetConfig()
	assert.NotNil(t, c)
	assert.NotNil(t, c.Prefix)
	assert.Equal(t, "eth", c.Prefix)
	assert.NotNil(t, c.PluginDirs)
	assert.Equal(t, DefaultCNIDir, c.PluginDirs[0])
	assert.NotNil(t, c.PluginConfDir)
	assert.Equal(t, confDir, c.PluginConfDir)
	assert.NotNil(t, c.Networks)
	assert.Equal(t, "plugin1", c.Networks[0].Config.Name)
	assert.Equal(t, "eth0", c.Networks[0].IFName)
}

// TestLibCNIType040 tests the cni version 0.4.0 plugin
// config and parses the result into structured data
func TestLibCNIType040(t *testing.T) {
	// Get the default CNI config
	l := defaultCNIConfig()
	// Create a fake cni config directory and file
	cniDir, confDir := makeFakeCNIConfig(t)
	defer tearDownCNIConfig(t, cniDir)
	l.pluginConfDir = confDir
	// Set the minimum network count as 2 for this test
	l.networkCount = 2
	err := l.Load(WithAllConf)
	assert.NoError(t, err)

	err = l.Status()
	assert.NoError(t, err)

	mockCNI := &MockCNI{}
	l.networks[0].cni = mockCNI
	l.networks[1].cni = mockCNI
	ipv4, err := types.ParseCIDR("10.0.0.1/24")
	assert.NoError(t, err)
	expectedRT := &cnilibrary.RuntimeConf{
		ContainerID:    "container-id1",
		NetNS:          "/proc/12345/ns/net",
		IfName:         "eth0",
		Args:           [][2]string(nil),
		CapabilityArgs: map[string]interface{}{},
	}
	mockCNI.On("AddNetworkList", l.networks[0].config, expectedRT).Return(&types040.Result{
		CNIVersion: "0.3.1", // covered by types040
		Interfaces: []*types040.Interface{
			{
				Name: "eth0",
			},
		},
		IPs: []*types040.IPConfig{
			{
				Version:   "4",
				Interface: types100.Int(0),
				Address:   *ipv4,
				Gateway:   net.ParseIP("10.0.0.255"),
			},
		},
	}, nil)
	mockCNI.On("DelNetworkList", l.networks[0].config, expectedRT).Return(nil)

	ipv4, err = types.ParseCIDR("10.0.0.2/24")
	assert.NoError(t, err)
	l.networks[1].cni = mockCNI
	expectedRT = &cnilibrary.RuntimeConf{
		ContainerID:    "container-id1",
		NetNS:          "/proc/12345/ns/net",
		IfName:         "eth1",
		Args:           [][2]string(nil),
		CapabilityArgs: map[string]interface{}{},
	}
	mockCNI.On("AddNetworkList", l.networks[1].config, expectedRT).Return(&types040.Result{
		CNIVersion: "0.3.1", // covered by types040
		Interfaces: []*types040.Interface{
			{
				Name: "eth1",
			},
		},
		IPs: []*types040.IPConfig{
			{
				Version:   "4",
				Interface: types100.Int(0),
				Address:   *ipv4,
				Gateway:   net.ParseIP("10.0.0.2"),
			},
		},
	}, nil)
	mockCNI.On("DelNetworkList", l.networks[1].config, expectedRT).Return(nil)

	ctx := context.Background()
	r, err := l.Setup(ctx, "container-id1", "/proc/12345/ns/net")
	assert.NoError(t, err)
	assert.Contains(t, r.Interfaces, "eth0")
	assert.NotNil(t, r.Interfaces["eth0"].IPConfigs)
	assert.Equal(t, r.Interfaces["eth0"].IPConfigs[0].IP.String(), "10.0.0.1")
	assert.Contains(t, r.Interfaces, "eth1")
	assert.NotNil(t, r.Interfaces["eth1"].IPConfigs)
	assert.Equal(t, r.Interfaces["eth1"].IPConfigs[0].IP.String(), "10.0.0.2")
	err = l.Remove(ctx, "container-id1", "/proc/12345/ns/net")
	assert.NoError(t, err)
}

// TestLibCNIType100 tests the cni version 1.0.0 plugin
// config and parses the result into structured data
func TestLibCNIType100(t *testing.T) {
	// Get the default CNI config
	l := defaultCNIConfig()
	// Create a fake cni config directory and file
	cniDir, confDir := makeFakeCNIConfig(t)
	defer tearDownCNIConfig(t, cniDir)
	l.pluginConfDir = confDir
	// Set the minimum network count as 2 for this test
	l.networkCount = 2
	err := l.Load(WithAllConf)
	assert.NoError(t, err)

	err = l.Status()
	assert.NoError(t, err)

	mockCNI := &MockCNI{}
	l.networks[0].cni = mockCNI
	l.networks[1].cni = mockCNI
	ipv4, err := types.ParseCIDR("10.0.0.1/24")
	assert.NoError(t, err)
	expectedRT := &cnilibrary.RuntimeConf{
		ContainerID:    "container-id1",
		NetNS:          "/proc/12345/ns/net",
		IfName:         "eth0",
		Args:           [][2]string(nil),
		CapabilityArgs: map[string]interface{}{},
	}
	mockCNI.On("AddNetworkList", l.networks[0].config, expectedRT).Return(&types100.Result{
		CNIVersion: "1.0.0",
		Interfaces: []*types100.Interface{
			{
				Name: "eth0",
			},
		},
		IPs: []*types100.IPConfig{
			{
				Interface: types100.Int(0),
				Address:   *ipv4,
				Gateway:   net.ParseIP("10.0.0.255"),
			},
		},
	}, nil)
	mockCNI.On("DelNetworkList", l.networks[0].config, expectedRT).Return(nil)
	mockCNI.On("CheckNetworkList", l.networks[0].config, expectedRT).Return(nil)
	ipv4, err = types.ParseCIDR("10.0.0.2/24")
	assert.NoError(t, err)
	l.networks[1].cni = mockCNI
	expectedRT = &cnilibrary.RuntimeConf{
		ContainerID:    "container-id1",
		NetNS:          "/proc/12345/ns/net",
		IfName:         "eth1",
		Args:           [][2]string(nil),
		CapabilityArgs: map[string]interface{}{},
	}
	mockCNI.On("AddNetworkList", l.networks[1].config, expectedRT).Return(&types100.Result{
		CNIVersion: "1.0.0",
		Interfaces: []*types100.Interface{
			{
				Name: "eth1",
			},
		},
		IPs: []*types100.IPConfig{
			{
				Interface: types100.Int(0),
				Address:   *ipv4,
				Gateway:   net.ParseIP("10.0.0.2"),
			},
		},
	}, nil)
	mockCNI.On("DelNetworkList", l.networks[1].config, expectedRT).Return(nil)
	mockCNI.On("CheckNetworkList", l.networks[1].config, expectedRT).Return(nil)
	ctx := context.Background()
	r, err := l.Setup(ctx, "container-id1", "/proc/12345/ns/net")
	assert.NoError(t, err)
	assert.Contains(t, r.Interfaces, "eth0")
	assert.NotNil(t, r.Interfaces["eth0"].IPConfigs)
	assert.Equal(t, r.Interfaces["eth0"].IPConfigs[0].IP.String(), "10.0.0.1")
	assert.Contains(t, r.Interfaces, "eth1")
	assert.NotNil(t, r.Interfaces["eth1"].IPConfigs)
	assert.Equal(t, r.Interfaces["eth1"].IPConfigs[0].IP.String(), "10.0.0.2")

	err = l.Check(ctx, "container-id1", "/proc/12345/ns/net")
	assert.NoError(t, err)

	err = l.Remove(ctx, "container-id1", "/proc/12345/ns/net")
	assert.NoError(t, err)
}

type MockCNI struct {
	mock.Mock
}

func (m *MockCNI) AddNetwork(ctx context.Context, net *cnilibrary.NetworkConfig, rt *cnilibrary.RuntimeConf) (types.Result, error) {
	args := m.Called(net, rt)
	return args.Get(0).(types.Result), args.Error(1)
}

func (m *MockCNI) DelNetwork(ctx context.Context, net *cnilibrary.NetworkConfig, rt *cnilibrary.RuntimeConf) error {
	args := m.Called(net, rt)
	return args.Error(0)
}

func (m *MockCNI) DelNetworkList(ctx context.Context, net *cnilibrary.NetworkConfigList, rt *cnilibrary.RuntimeConf) error {
	args := m.Called(net, rt)
	return args.Error(0)
}

func (m *MockCNI) AddNetworkList(ctx context.Context, net *cnilibrary.NetworkConfigList, rt *cnilibrary.RuntimeConf) (types.Result, error) {
	args := m.Called(net, rt)
	return args.Get(0).(types.Result), args.Error(1)
}

func (m *MockCNI) CheckNetworkList(ctx context.Context, net *cnilibrary.NetworkConfigList, rt *cnilibrary.RuntimeConf) error {
	args := m.Called(net, rt)
	return args.Error(0)
}

func (m *MockCNI) CheckNetwork(ctx context.Context, net *cnilibrary.NetworkConfig, rt *cnilibrary.RuntimeConf) error {
	args := m.Called(net, rt)
	return args.Error(0)
}

func (m *MockCNI) GetNetworkCachedConfig(net *cnilibrary.NetworkConfig, rt *cnilibrary.RuntimeConf) ([]byte, *cnilibrary.RuntimeConf, error) {
	args := m.Called(net, rt)
	return args.Get(0).([]byte), args.Get(1).(*cnilibrary.RuntimeConf), args.Error(1)
}

func (m *MockCNI) GetNetworkCachedResult(net *cnilibrary.NetworkConfig, rt *cnilibrary.RuntimeConf) (types.Result, error) {
	args := m.Called(net, rt)
	return args.Get(0).(types.Result), args.Error(1)
}

func (m *MockCNI) ValidateNetworkList(ctx context.Context, net *cnilibrary.NetworkConfigList) ([]string, error) {
	args := m.Called(net)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCNI) ValidateNetwork(ctx context.Context, net *cnilibrary.NetworkConfig) ([]string, error) {
	args := m.Called(net)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCNI) GetNetworkListCachedConfig(net *cnilibrary.NetworkConfigList, rt *cnilibrary.RuntimeConf) ([]byte, *cnilibrary.RuntimeConf, error) {
	args := m.Called(net, rt)
	return args.Get(0).([]byte), args.Get(1).(*cnilibrary.RuntimeConf), args.Error(1)
}

func (m *MockCNI) GetNetworkListCachedResult(net *cnilibrary.NetworkConfigList, rt *cnilibrary.RuntimeConf) (types.Result, error) {
	args := m.Called(net, rt)
	return args.Get(0).(types.Result), args.Error(1)
}
