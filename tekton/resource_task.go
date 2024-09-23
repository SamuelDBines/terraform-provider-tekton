package tekton

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	triggersclient "github.com/tektoncd/triggers/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// resourceTektonTask defines a Tekton Task.
func resourceTektonTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTektonTaskCreate,
		Read:   resourceTektonTaskRead,
		Update: resourceTektonTaskUpdate,
		Delete: resourceTektonTaskDelete,

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
			"steps": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"image": {
							Type:     schema.TypeString,
							Required: true,
						},
						"command": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"workspaces": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func getTaskWorkspaces(tfWorkspaces []interface{}) []tektonv1beta1.WorkspaceDeclaration {
	var workspaces []tektonv1beta1.WorkspaceDeclaration

	for _, tfWorkspace := range tfWorkspaces {
		workspaceData := tfWorkspace.(map[string]interface{})
		workspace := tektonv1beta1.WorkspaceDeclaration{
			Name: workspaceData["name"].(string),
		}

		if v, ok := workspaceData["description"]; ok {
			workspace.Description = v.(string)
		}

		workspaces = append(workspaces, workspace)
	}

	return workspaces
}

// resourceTektonTaskCreate creates a Tekton Task.
func resourceTektonTaskCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(struct {
		TektonClient         *tektonclient.Clientset
		TektonTriggersClient *triggersclient.Clientset
	})
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)

	steps := getTaskSteps(d.Get("steps").([]interface{}))
	workspaces := getTaskWorkspaces(d.Get("workspaces").([]interface{}))

	task := &tektonv1beta1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: tektonv1beta1.TaskSpec{
			Steps:      steps,
			Workspaces: workspaces,
		},
	}

	_, err := clients.TektonClient.TektonV1beta1().Tasks(namespace).Create(context.Background(), task, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Tekton Task: %v", err)
	}

	d.SetId(name)
	return resourceTektonTaskRead(d, m)
}

func resourceTektonTaskRead(d *schema.ResourceData, m interface{}) error {
	// Implement read logic
	return nil
}

func resourceTektonTaskUpdate(d *schema.ResourceData, m interface{}) error {
	// Implement update logic
	return nil
}

func resourceTektonTaskDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(struct {
		TektonClient         *tektonclient.Clientset
		TektonTriggersClient *triggersclient.Clientset
	})
	name := d.Id()
	namespace := d.Get("namespace").(string)

	err := clients.TektonClient.TektonV1beta1().Tasks(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete Tekton Task: %v", err)
	}

	d.SetId("")
	return nil
}

// Helper function to convert Terraform steps to Tekton steps
func getTaskSteps(tfSteps []interface{}) []tektonv1beta1.Step {
	var steps []tektonv1beta1.Step

	for _, tfStep := range tfSteps {
		stepData := tfStep.(map[string]interface{})
		step := tektonv1beta1.Step{
			Name:    stepData["name"].(string),
			Image:   stepData["image"].(string),
			Command: toStringSlice(stepData["command"].([]interface{})),
		}
		steps = append(steps, step)
	}

	return steps
}

func toStringSlice(tfList []interface{}) []string {
	var result []string
	for _, v := range tfList {
		result = append(result, v.(string))
	}
	return result
}
