package tekton

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// resourceTektonTaskRun defines a Tekton TaskRun.
func resourceTektonTaskRun() *schema.Resource {
	return &schema.Resource{
		Create: resourceTektonTaskRunCreate,
		Read:   resourceTektonTaskRunRead,
		Update: resourceTektonTaskRunUpdate,
		Delete: resourceTektonTaskRunDelete,

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
			"task_ref_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Tekton Task to run",
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
			"service_account_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
		},
	}
}

// resourceTektonTaskRunCreate creates a Tekton TaskRun.
func resourceTektonTaskRunCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tektonclient.Clientset)
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)
	taskRefName := d.Get("task_ref_name").(string)
	serviceAccountName := d.Get("service_account_name").(string)

	params := getTaskRunParams(d.Get("params").([]interface{}))

	taskRun := &tektonv1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: tektonv1beta1.TaskRunSpec{
			TaskRef: &tektonv1beta1.TaskRef{
				Name: taskRefName,
			},
			ServiceAccountName: serviceAccountName,
			Params:             params,
		},
	}

	_, err := client.TektonV1beta1().TaskRuns(namespace).Create(context.Background(), taskRun, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Tekton TaskRun: %v", err)
	}

	d.SetId(name)
	return resourceTektonTaskRunRead(d, m)
}

// resourceTektonTaskRunRead reads the state of a Tekton TaskRun.
func resourceTektonTaskRunRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*tektonclient.Clientset)
	name := d.Id()
	namespace := d.Get("namespace").(string)

	_, err := client.TektonV1beta1().TaskRuns(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		// If the task run is not found, we should remove it from the state
		d.SetId("")
		return nil
	}

	return nil
}

// resourceTektonTaskRunUpdate updates a Tekton TaskRun (if necessary).
func resourceTektonTaskRunUpdate(d *schema.ResourceData, m interface{}) error {
	// TaskRuns are typically immutable after creation, so you may want to handle this accordingly
	return resourceTektonTaskRunRead(d, m)
}

// resourceTektonTaskRunDelete deletes a Tekton TaskRun.
func resourceTektonTaskRunDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*tektonclient.Clientset)
	name := d.Id()
	namespace := d.Get("namespace").(string)

	err := client.TektonV1beta1().TaskRuns(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete Tekton TaskRun: %v", err)
	}

	d.SetId("")
	return nil
}

// Helper function to convert Terraform params to Tekton params
func getTaskRunParams(tfParams []interface{}) []tektonv1beta1.Param {
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
