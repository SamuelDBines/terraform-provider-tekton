package tekton

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// resourceTektonPipelineRun defines a Tekton PipelineRun resource.
func resourceTektonPipelineRun() *schema.Resource {
	return &schema.Resource{
		Create: resourceTektonPipelineRunCreate,
		Read:   resourceTektonPipelineRunRead,
		Update: resourceTektonPipelineRunUpdate,
		Delete: resourceTektonPipelineRunDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"pipeline_ref_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Tekton Pipeline to reference in this PipelineRun.",
			},
			"service_account_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"params": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

// resourceTektonPipelineRunCreate creates a Tekton PipelineRun.
func resourceTektonPipelineRunCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tektonclient.Clientset)
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)
	pipelineRefName := d.Get("pipeline_ref_name").(string)
	serviceAccountName := d.Get("service_account_name").(string)

	params := getPipelineRunParams(d.Get("params").([]interface{}))

	pipelineRun := &tektonv1beta1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: tektonv1beta1.PipelineRunSpec{
			PipelineRef: &tektonv1beta1.PipelineRef{
				Name: pipelineRefName,
			},
			ServiceAccountName: serviceAccountName,
			Params:             params,
		},
	}

	_, err := client.TektonV1beta1().PipelineRuns(namespace).Create(context.Background(), pipelineRun, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Tekton PipelineRun: %v", err)
	}

	d.SetId(name)
	return resourceTektonPipelineRunRead(d, m)
}

// resourceTektonPipelineRunRead reads the state of a Tekton PipelineRun.
func resourceTektonPipelineRunRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*tektonclient.Clientset)
	name := d.Id()
	namespace := d.Get("namespace").(string)

	_, err := client.TektonV1beta1().PipelineRuns(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		// If the pipeline run is not found, remove it from the state
		d.SetId("")
		return nil
	}

	return nil
}

// resourceTektonPipelineRunUpdate updates a Tekton PipelineRun (if necessary).
func resourceTektonPipelineRunUpdate(d *schema.ResourceData, m interface{}) error {
	// PipelineRuns are typically immutable once created. Handle as necessary.
	return resourceTektonPipelineRunRead(d, m)
}

// resourceTektonPipelineRunDelete deletes a Tekton PipelineRun.
func resourceTektonPipelineRunDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*tektonclient.Clientset)
	name := d.Id()
	namespace := d.Get("namespace").(string)

	err := client.TektonV1beta1().PipelineRuns(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete Tekton PipelineRun: %v", err)
	}

	d.SetId("")
	return nil
}

// Helper function to convert Terraform params into Tekton params
func getPipelineRunParams(tfParams []interface{}) []tektonv1beta1.Param {
	var params []tektonv1beta1.Param

	for _, tfParam := range tfParams {
		paramData := tfParam.(map[string]interface{})
		param := tektonv1beta1.Param{
			Name: paramData["name"].(string),
			Value: tektonv1beta1.ArrayOrString{
				Type:      tektonv1beta1.ParamTypeString,
				StringVal: paramData["value"].(string),
			},
		}
		params = append(params, param)
	}

	return params
}
