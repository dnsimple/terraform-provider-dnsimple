package common

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Timeouts struct {
	Create types.String `tfsdk:"create"`
	Update types.String `tfsdk:"update"`
	Delete types.String `tfsdk:"delete"`
}

func (t Timeouts) CreateDuration() time.Duration {
	if !t.Create.IsUnknown() && !t.Create.IsNull() {
		d, _ := time.ParseDuration(t.Create.ValueString())
		return d
	}

	d, _ := time.ParseDuration("30s")
	return d
}

func (t Timeouts) UpdateDuration() time.Duration {
	if !t.Update.IsUnknown() && !t.Update.IsNull() {
		d, _ := time.ParseDuration(t.Update.ValueString())
		return d
	}

	d, _ := time.ParseDuration("30s")
	return d
}

func (t Timeouts) DeleteDuration() time.Duration {
	if !t.Delete.IsUnknown() && !t.Delete.IsNull() {
		d, _ := time.ParseDuration(t.Delete.ValueString())
		return d
	}

	d, _ := time.ParseDuration("30s")
	return d
}
