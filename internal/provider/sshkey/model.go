// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SSHKeyResourceModel describes the resource data model.
type SSHKeyResourceModel struct {
	Label types.String `tfsdk:"label"`
	Key   types.String `tfsdk:"key"`
	Id    types.String `tfsdk:"id"`
}
