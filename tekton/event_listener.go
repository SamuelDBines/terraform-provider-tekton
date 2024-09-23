package tekton

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	tektonv1alpha1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func resourceTektonEventListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceTektonEventListenerCreate,
		Read:   resourceTektonEventListenerRead,
		Update: resourceTektonEventListenerUpdate,
		Delete: resourceTektonEventListenerDelete,

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
			"triggers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger_template_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"trigger_binding_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceTektonEventListenerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tektonclient.Clientset)
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)

	triggers := getEventListenerTriggers(d.Get("triggers").([]interface{}))

	eventListener := &tektonv1alpha1.EventListener{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: tektonv1alpha1.EventListenerSpec{
			Triggers: triggers,
		},
	}

	_, err := client.TektonV1alpha1().EventListeners(namespace).Create(context.Background(), eventListener, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Tekton EventListener: %v", err)
	}

	d.SetId(name)
	return resourceTektonEventListenerRead(d, m)
}

func resourceTektonEventListenerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*tektonclient.Clientset)
	name := d.Id()
	namespace := d.Get("namespace").(string)

	_, err := client.TektonV1alpha1().EventListeners(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceTektonEventListenerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceTektonEventListenerRead(d, m)
}

func resourceTektonEventListenerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*tektonclient.Clientset)
	name := d.Id()
	namespace := d.Get("namespace").(string)

	err := client.TektonV1alpha1().EventListeners(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete Tekton EventListener: %v", err)
	}

	d.SetId("")
	return nil
}

func getEventListenerTriggers(tfTriggers []interface{}) []tektonv1alpha1.EventListenerTrigger {
	var triggers []tektonv1alpha1.EventListenerTrigger
	for _, tfTrigger := range tfTriggers {
		triggerData := tfTrigger.(map[string]interface{})
		trigger := tektonv1alpha1.EventListenerTrigger{
			Template: tektonv1alpha1.EventListenerTemplate{
				Name: triggerData["trigger_template_name"].(string),
			},
			Bindings: []*tektonv1alpha1.EventListenerBinding{
				{
					Name: triggerData["trigger_binding_name"].(string),
				},
			},
		}
		triggers = append(triggers, trigger)
	}
	return triggers
}
