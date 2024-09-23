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

func resourceTektonTriggerBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceTektonTriggerBindingCreate,
		Read:   resourceTektonTriggerBindingRead,
		Update: resourceTektonTriggerBindingUpdate,
		Delete: resourceTektonTriggerBindingDelete,

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
			"bindings": {
				Type:     schema.TypeList,
				Required: true,
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

func resourceTektonTriggerBindingCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(struct {
		TektonClient         *tektonclient.Clientset
		TektonTriggersClient *triggersclient.Clientset
	})
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)

	bindings := getTriggerBindingParams(d.Get("bindings").([]interface{}))

	triggerBinding := &tektonv1alpha1.TriggerBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: tektonv1alpha1.TriggerBindingSpec{
			Params: bindings,
		},
	}

	_, err := clients.TektonTriggersClient.TriggersV1alpha1().TriggerBindings(namespace).Create(context.Background(), triggerBinding, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Tekton TriggerBinding: %v", err)
	}

	d.SetId(name)
	return resourceTektonTriggerBindingRead(d, m)
}

func resourceTektonTriggerBindingRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(struct {
		TektonClient         *tektonclient.Clientset
		TektonTriggersClient *triggersclient.Clientset
	})
	name := d.Id()
	namespace := d.Get("namespace").(string)

	_, err := clients.TektonTriggersClient.TriggersV1alpha1().TriggerBindings(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceTektonTriggerBindingUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceTektonTriggerBindingRead(d, m)
}

func resourceTektonTriggerBindingDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(struct {
		TektonClient         *tektonclient.Clientset
		TektonTriggersClient *triggersclient.Clientset
	})
	name := d.Id()
	namespace := d.Get("namespace").(string)

	err := clients.TektonTriggersClient.TriggersV1alpha1().TriggerBindings(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete Tekton TriggerBinding: %v", err)
	}

	d.SetId("")
	return nil
}

func getTriggerBindingParams(tfBindings []interface{}) []tektonv1alpha1.Param {
	var bindings []tektonv1alpha1.Param
	for _, tfBinding := range tfBindings {
		bindingData := tfBinding.(map[string]interface{})
		binding := tektonv1alpha1.Param{
			Name:  bindingData["name"].(string),
			Value: tektonv1alpha1.ArrayOrString{StringVal: bindingData["value"].(string)},
		}
		bindings = append(bindings, binding)
	}
	return bindings
}
