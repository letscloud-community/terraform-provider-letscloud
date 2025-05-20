// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/letscloud-community/letscloud-go/domains"
)

// letsCloudClientMock é um mock que implementa a interface LetsCloudClient.
type letsCloudClientMock struct {
	instances map[string]*domains.Instance
}

func (m *letsCloudClientMock) CreateInstance(req *domains.CreateInstanceRequest) error {
	if m.instances == nil {
		m.instances = make(map[string]*domains.Instance)
	}
	inst := &domains.Instance{
		Identifier: req.Label,
		Label:      req.Label,
		Hostname:   req.Hostname,
		Built:      true,
		Booted:     true,
		Location:   domains.Location{Slug: req.LocationSlug},
		IPAddresses: []domains.IPAddress{
			{Address: "192.168.1.1"},
			{Address: "2001:db8::1"},
		},
	}
	m.instances[req.Label] = inst
	return nil
}

func (m *letsCloudClientMock) Instance(id string) (*domains.Instance, error) {
	if m.instances == nil {
		m.instances = make(map[string]*domains.Instance)
	}
	inst, ok := m.instances[id]
	if !ok {
		return nil, fmt.Errorf("Instance not found")
	}
	return inst, nil
}

func (m *letsCloudClientMock) DeleteInstance(id string) error {
	if m.instances == nil {
		m.instances = make(map[string]*domains.Instance)
	}
	delete(m.instances, id)
	return nil
}

func (m *letsCloudClientMock) ResetPasswordInstance(id string, password string) error {
	if m.instances == nil {
		m.instances = make(map[string]*domains.Instance)
	}
	_, ok := m.instances[id]
	if !ok {
		return fmt.Errorf("Instance not found")
	}
	// Just a mock, so we don't store password, but could add a field if needed
	return nil
}

func (m *letsCloudClientMock) LocationPlans(location string) ([]domains.Plan, error) {
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

func (m *letsCloudClientMock) UpdateInstance(id, newLabel, newHostname string) error {
	if m.instances == nil {
		m.instances = make(map[string]*domains.Instance)
	}
	inst, ok := m.instances[id]
	if !ok {
		return fmt.Errorf("Instance not found")
	}
	inst.Label = newLabel
	inst.Hostname = newHostname
	return nil
}
