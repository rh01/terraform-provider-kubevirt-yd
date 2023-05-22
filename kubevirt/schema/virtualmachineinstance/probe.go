package virtualmachineinstance

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/schema/k8s"
	api "k8s.io/api/core/v1"
	utilValidation "k8s.io/apimachinery/pkg/util/validation"
	kubevirtapiv1 "kubevirt.io/client-go/api/v1"
)

func probeFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"http_get": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Specifies the http request to perform.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"host": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: `Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.`,
					},
					"path": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: `Path to access on the HTTP server.`,
					},
					"scheme": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     string(api.URISchemeHTTP),
						Description: `Scheme to use for connecting to the host.`,
						ValidateFunc: validation.StringInSlice([]string{
							string(api.URISchemeHTTP),
							string(api.URISchemeHTTPS),
						}, false),
					},
					"port": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validatePortNumOrName,
						Description:  `Name or number of the port to access on the container. Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME.`,
					},
					"http_header": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: `Scheme to use for connecting to the host.`,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "The header field name",
								},
								"value": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "The header field value",
								},
							},
						},
					},
				},
			},
		},
		"tcp_socket": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "TCPSocket specifies an action involving a TCP port. TCP hooks not yet supported",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"port": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validatePortNumOrName,
						Description:  "Number or name of the port to access on the container. Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME.",
					},
				},
			},
		},
		"initial_delay_seconds": {
			Type:        schema.TypeInt,
			Description: "Number of seconds after the VirtualMachineInstance has started before liveness probes are initiated.",
			Optional:    true,
		},
		"timeout_seconds": {
			Type:        schema.TypeInt,
			Description: "Number of seconds after which the probe times out.",
			Optional:    true,
		},
		"period_seconds": {
			Type:        schema.TypeInt,
			Description: "How often (in seconds) to perform the probe.",
			Optional:    true,
		},
		"success_threshold": {
			Type:        schema.TypeInt,
			Description: "Minimum consecutive successes for the probe to be considered successful after having failed.",
			Optional:    true,
		},
		"failure_threshold": {
			Type:        schema.TypeInt,
			Description: "Minimum consecutive failures for the probe to be considered failed after having succeeded.",
			Optional:    true,
		},
	}
}

func validatePortNum(value interface{}, key string) (ws []string, es []error) {
	errors := utilValidation.IsValidPortNum(value.(int))
	if len(errors) > 0 {
		for _, err := range errors {
			es = append(es, fmt.Errorf("%s %s", key, err))
		}
	}
	return
}

func validatePortName(value interface{}, key string) (ws []string, es []error) {
	errors := utilValidation.IsValidPortName(value.(string))
	if len(errors) > 0 {
		for _, err := range errors {
			es = append(es, fmt.Errorf("%s %s", key, err))
		}
	}
	return
}

func validatePortNumOrName(value interface{}, key string) (ws []string, es []error) {
	switch value.(type) {
	case string:
		intVal, err := strconv.Atoi(value.(string))
		if err != nil {
			return validatePortName(value, key)
		}
		return validatePortNum(intVal, key)
	case int:
		return validatePortNum(value, key)

	default:
		es = append(es, fmt.Errorf("%s must be defined of type string or int on the schema", key))
		return
	}
}

func probeSchema() *schema.Schema {
	fields := probeFields()

	return &schema.Schema{
		Type: schema.TypeList,

		Description: fmt.Sprintf("Specification of the desired behavior of the VirtualMachineInstance on the host."),
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandProbe(probe []interface{}) *kubevirtapiv1.Probe {
	if len(probe) == 0 || probe[0] == nil {
		return nil
	}

	result := &kubevirtapiv1.Probe{}

	in := probe[0].(map[string]interface{})

	if v, ok := in["http_get"].([]interface{}); ok && len(v) > 0 {
		result.HTTPGet = k8s.ExpandHTTPGet(v)
	}
	if v, ok := in["tcp_socket"].([]interface{}); ok && len(v) > 0 {
		result.TCPSocket = k8s.ExpandTCPSocket(v)
	}

	return result
}

func flattenProbe(in kubevirtapiv1.Probe) []interface{} {
	att := make(map[string]interface{})

	if in.HTTPGet != nil {
		att["http_get"] = k8s.FlattenHTTPGet(in.HTTPGet)
	}
	if in.TCPSocket != nil {
		att["tcp_socket"] = k8s.FlattenTCPSocket(in.TCPSocket)
	}

	return []interface{}{att}
}
