// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"fmt"

	"github.com/letscloud-community/letscloud-go/domains"
	"github.com/letscloud-community/terraform-provider-letscloud/internal/provider/client"
)

// MockLetsCloudClient is a mock implementation of LetsCloudClient for testing.
type MockLetsCloudClient struct {
	sshKeys   map[string]*domains.SSHKey
	instances map[string]*domains.Instance
}

// NewMockLetsCloudClient creates a new mock client.
func NewMockLetsCloudClient() client.LetsCloudClient {
	return &MockLetsCloudClient{
		sshKeys:   make(map[string]*domains.SSHKey),
		instances: make(map[string]*domains.Instance),
	}
}

// Close closes the mock client.
func (m *MockLetsCloudClient) Close() {
	// Nothing to do for mock client
}

// SSH Key methods.
func (m *MockLetsCloudClient) SSHKey(id string) (*domains.SSHKey, error) {
	if key, exists := m.sshKeys[id]; exists {
		return key, nil
	}
	return nil, fmt.Errorf("SSH key not found: %s", id)
}

func (m *MockLetsCloudClient) SSHKeys() ([]domains.SSHKey, error) {
	keys := make([]domains.SSHKey, 0, len(m.sshKeys))
	for _, key := range m.sshKeys {
		keys = append(keys, *key)
	}
	return keys, nil
}

func (m *MockLetsCloudClient) CreateSSHKey(req *domains.SSHKeyCreateRequest) (*domains.SSHKey, error) {
	// Check if label already exists
	for _, key := range m.sshKeys {
		if key.Title == req.Title {
			return nil, fmt.Errorf("SSH key with label '%s' already exists", req.Title)
		}
	}

	// Create new SSH key
	id := fmt.Sprintf("mock-ssh-key-%d", len(m.sshKeys)+1)
	key := &domains.SSHKey{
		Slug:      id,
		Title:     req.Title,
		PublicKey: req.Key,
	}
	m.sshKeys[id] = key
	return key, nil
}

func (m *MockLetsCloudClient) DeleteSSHKey(id string) error {
	if _, exists := m.sshKeys[id]; !exists {
		return fmt.Errorf("SSH key not found: %s", id)
	}
	delete(m.sshKeys, id)
	return nil
}

// Instance methods.
func (m *MockLetsCloudClient) Instance(id string) (*domains.Instance, error) {
	if instance, exists := m.instances[id]; exists {
		// Simulate instance building process
		if !instance.Built {
			instance.Built = true
		} else if !instance.Booted {
			instance.Booted = true
		}
		return instance, nil
	}
	return nil, fmt.Errorf("Instance not found: %s", id)
}

func (m *MockLetsCloudClient) Instances() ([]domains.Instance, error) {
	instances := make([]domains.Instance, 0, len(m.instances))
	for _, instance := range m.instances {
		// Simulate instance building process
		if !instance.Built {
			instance.Built = true
		} else if !instance.Booted {
			instance.Booted = true
		}
		instances = append(instances, *instance)
	}
	return instances, nil
}

func (m *MockLetsCloudClient) CreateInstance(req *domains.CreateInstanceRequest) error {
	id := fmt.Sprintf("mock-instance-%d", len(m.instances)+1)
	instance := &domains.Instance{
		Identifier: id,
		Label:      req.Label,
		Hostname:   req.Hostname,
		Built:      false,
		Booted:     false,
		Location:   domains.Location{Slug: req.LocationSlug},
		IPAddresses: []domains.IPAddress{
			{Address: "192.168.1.1"},
			{Address: "2001:db8::1"},
		},
	}
	m.instances[id] = instance
	return nil
}

func (m *MockLetsCloudClient) DeleteInstance(id string) error {
	if _, exists := m.instances[id]; !exists {
		return fmt.Errorf("Instance not found: %s", id)
	}
	delete(m.instances, id)
	return nil
}

func (m *MockLetsCloudClient) ResetPasswordInstance(id string, password string) error {
	if _, exists := m.instances[id]; !exists {
		return fmt.Errorf("Instance not found: %s", id)
	}
	return nil
}

func (m *MockLetsCloudClient) LocationPlans(location string) ([]domains.Plan, error) {
	return []domains.Plan{
		{
			Slug:         "plan-1",
			Shortcode:    "Basic Plan",
			Core:         1,
			Memory:       1024,
			Disk:         10,
			Bandwidth:    1000,
			MonthlyValue: "10.00",
			CurrencyCode: "USD",
		},
	}, nil
}
