package datavolume

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/utils"
	kubevirtapiv1 "kubevirt.io/api/core/v1"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

func dataVolumeStatusFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"phase": {
			Type:        schema.TypeString,
			Description: "DataVolumePhase is the current phase of the DataVolume.",
			Optional:    true,
			Computed:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"",
				"Pending",
				"PVCBound",
				"ImportScheduled",
				"ImportInProgress",
				"CloneScheduled",
				"CloneInProgress",
				"SnapshotForSmartCloneInProgress",
				"SmartClonePVCInProgress",
				"UploadScheduled",
				"UploadReady",
				"Succeeded",
				"Failed",
				"Unknown",
			}, false),
		},
		"progress": {
			Type:         schema.TypeString,
			Description:  "DataVolumePhase is the current phase of the DataVolume.",
			Optional:     true,
			Computed:     true,
			ValidateFunc: utils.StringIsIntInRange(0, 100),
		},
	}
}

func dataVolumeStatusSchema() *schema.Schema {
	fields := dataVolumeStatusFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: fmt.Sprintf("DataVolumeStatus provides the parameters to store the phase of the Data Volume"),
		Optional:    true,
		MaxItems:    1,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandDataVolumeStatus(dataVolumeStatus []interface{}) *kubevirtapiv1.DataVolumeTemplateDummyStatus {
	result := &kubevirtapiv1.DataVolumeTemplateDummyStatus{}

	if len(dataVolumeStatus) == 0 || dataVolumeStatus[0] == nil {
		return result
	}

	// in := dataVolumeStatus[0].(map[string]interface{})

	// if v, ok := in["phase"].(string); ok {
	// 	result = cdiv1.DataVolumePhase(v)
	// }
	// if v, ok := in["progress"].(string); ok {
	// 	result.Progress = cdiv1.DataVolumeProgress(v)
	// }

	return result
}

func flattenDataVolumeStatus(in cdiv1.DataVolumeStatus) []interface{} {
	att := map[string]interface{}{
		"phase":    string(in.Phase),
		"progress": string(in.Progress),
	}
	return []interface{}{att}
}
