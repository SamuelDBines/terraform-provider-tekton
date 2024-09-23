package tekton

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	tektonv1alpha1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	triggersclient "github.com/tektoncd/triggers/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// resourceTektonTriggerTemplate defines a Tekton TriggerTemplate.
func resourceTektonTriggerTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTektonTriggerTemplateCreate,
		Read:   resourceTektonTriggerTemplateRead,
		Update: resourceTektonTriggerTemplateUpdate,
		Delete: resourceTektonTriggerTemplateDelete,

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
			"params": {
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
			"resourcetemplates": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:     schema.TypeString,
							Required: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"metadata": {
							Type:     schema.TypeMap,
							Required: true,
						},
						"spec": {
							Type:     schema.TypeMap,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceTektonTriggerTemplateCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(struct {
		TektonClient         *tektonclient.Clientset
		TektonTriggersClient *triggersclient.Clientset
	})
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)

	params := getTriggerTemplateParams(d.Get("params").([]interface{}))
	resourceTemplates := getResourceTemplates(d.Get("resourcetemplates").([]interface{}))

	triggerTemplate := &tektonv1alpha1.TriggerTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: tektonv1alpha1.TriggerTemplateSpec{
			Params:            params,
			ResourceTemplates: resourceTemplates,
		},
	}

	_, err := clients.TektonTriggersClient.TriggersV1alpha1().TriggerTemplates(namespace).Create(context.Background(), triggerTemplate, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Tekton TriggerTemplate: %v", err)
	}

	d.SetId(name)
	return resourceTektonTriggerTemplateRead(d, m)
}

func resourceTektonTriggerTemplateRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(struct {
		TektonClient         *tektonclient.Clientset
		TektonTriggersClient *triggersclient.Clientset
	})
	name := d.Id()
	namespace := d.Get("namespace").(string)

	_, err := clients.TektonTriggersClient.TriggersV1alpha1().TriggerTemplates(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceTektonTriggerTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	// TriggerTemplates are typically immutable after creation, so you may not allow updates
	return resourceTektonTriggerTemplateRead(d, m)
}

func resourceTektonTriggerTemplateDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(struct {
		TektonClient         *tektonclient.Clientset
		TektonTriggersClient *triggersclient.Clientset
	})
	name := d.Id()
	namespace := d.Get("namespace").(string)

	err := clients.TektonTriggersClient.TriggersV1alpha1().TriggerTemplates(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete Tekton TriggerTemplate: %v", err)
	}

	d.SetId("")
	return nil
}

// Helper function to convert Terraform params into Tekton params
func getTriggerTemplateParams(tfParams []interface{}) []tektonv1alpha1.ParamSpec {
	var params []tektonv1alpha1.ParamSpec
	for _, tfParam := range tfParams {
		paramData := tfParam.(map[string]interface{})
		param := tektonv1alpha1.ParamSpec{
			Name: paramData["name"].(string),
		}
		if v, ok := paramData["description"]; ok {
			param.Description = v.(string)
		}
		params = append(params, param)
	}
	return params
}

// Helper function to convert resource templates for Tekton
func getResourceTemplates(tfResourceTemplates []interface{}) []tektonv1alpha1.TriggerResourceTemplate {
	var templates []tektonv1alpha1.TriggerResourceTemplate
	for _, tfTemplate := range tfResourceTemplates {
		templateData := tfTemplate.(map[string]interface{})
		template := tektonv1alpha1.TriggerResourceTemplate{
			APIVersion: templateData["api_version"].(string),
			Kind:       templateData["kind"].(string),
			Metadata:   templateData["metadata"].(map[string]interface{}),
			Spec:       templateData["spec"].(map[string]interface{}),
		}
		templates = append(templates, template)
	}
	return templates
}
