package tekton

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// resourceTektonPipeline defines a Tekton Pipeline resource.
func resourceTektonPipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourceTektonPipelineCreate,
		Read:   resourceTektonPipelineRead,
		Update: resourceTektonPipelineUpdate,
		Delete: resourceTektonPipelineDelete,

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
			"tasks": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"task_ref_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the Tekton Task to reference in this Pipeline",
						},
						"run_after": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Tasks that should run after this task.",
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
									"workspace_ref": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
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
					},
				},
			},
		},
	}
}

func getPipelineWorkspaces(tfWorkspaces []interface{}) []tektonv1beta1.PipelineWorkspaceDeclaration {
	var workspaces []tektonv1beta1.PipelineWorkspaceDeclaration

	for _, tfWorkspace := range tfWorkspaces {
		workspaceData := tfWorkspace.(map[string]interface{})
		workspace := tektonv1beta1.PipelineWorkspaceDeclaration{
			Name: workspaceData["name"].(string),
		}

		workspaces = append(workspaces, workspace)
	}

	return workspaces
}

func getPipelineTaskWorkspaces(tfWorkspaces []interface{}) []tektonv1beta1.WorkspacePipelineTaskBinding {
	var workspaces []tektonv1beta1.WorkspacePipelineTaskBinding

	for _, tfWorkspace := range tfWorkspaces {
		workspaceData := tfWorkspace.(map[string]interface{})
		workspace := tektonv1beta1.WorkspacePipelineTaskBinding{
			Name:      workspaceData["name"].(string),
			Workspace: workspaceData["workspace_ref"].(string),
		}

		workspaces = append(workspaces, workspace)
	}

	return workspaces
}

// resourceTektonPipelineCreate creates a Tekton Pipeline.
func resourceTektonPipelineCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tektonclient.Clientset)
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)

	tasks := getPipelineTasks(d.Get("tasks").([]interface{}))
	workspaces := getPipelineWorkspaces(d.Get("workspaces").([]interface{}))

	pipeline := &tektonv1beta1.Pipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: tektonv1beta1.PipelineSpec{
			Tasks:      tasks,
			Workspaces: workspaces,
		},
	}

	_, err := client.TektonV1beta1().Pipelines(namespace).Create(context.Background(), pipeline, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Tekton Pipeline: %v", err)
	}

	d.SetId(name)
	return resourceTektonPipelineRead(d, m)
}

// resourceTektonPipelineRead reads the state of a Tekton Pipeline.
func resourceTektonPipelineRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*tektonclient.Clientset)
	name := d.Id()
	namespace := d.Get("namespace").(string)

	_, err := client.TektonV1beta1().Pipelines(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		// If the pipeline is not found, remove it from the state
		d.SetId("")
		return nil
	}

	return nil
}

// resourceTektonPipelineUpdate updates a Tekton Pipeline.
func resourceTektonPipelineUpdate(d *schema.ResourceData, m interface{}) error {
	// Tekton Pipelines are typically immutable once created. You may handle update logic based on your needs.
	return resourceTektonPipelineRead(d, m)
}

// resourceTektonPipelineDelete deletes a Tekton Pipeline.
func resourceTektonPipelineDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*tektonclient.Clientset)
	name := d.Id()
	namespace := d.Get("namespace").(string)

	err := client.TektonV1beta1().Pipelines(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete Tekton Pipeline: %v", err)
	}

	d.SetId("")
	return nil
}

// Helper function to convert Terraform tasks into Tekton pipeline tasks
func getPipelineTasks(tfTasks []interface{}) []tektonv1beta1.PipelineTask {
	var tasks []tektonv1beta1.PipelineTask

	for _, tfTask := range tfTasks {
		taskData := tfTask.(map[string]interface{})
		task := tektonv1beta1.PipelineTask{
			Name: taskData["name"].(string),
			TaskRef: &tektonv1beta1.TaskRef{
				Name: taskData["task_ref_name"].(string),
			},
		}

		// Add run_after tasks if specified
		if v, ok := taskData["run_after"]; ok {
			task.RunAfter = toStringSlice(v.([]interface{}))
		}

		tasks = append(tasks, task)
	}

	return tasks
}

// Helper function to convert a Terraform list to a string slice
func toStringSlice(tfList []interface{}) []string {
	var result []string
	for _, v := range tfList {
		result = append(result, v.(string))
	}
	return result
}
