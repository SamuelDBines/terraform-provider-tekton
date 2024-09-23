package tekton

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	triggersclient "github.com/tektoncd/triggers/pkg/client/clientset/versioned"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Provider defines the provider schema and resources.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"kubeconfig": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KUBECONFIG", nil),
				Description: "Path to the Kubernetes configuration file.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"tekton_task":            resourceTektonTask(),
			"tekton_taskrun":         resourceTektonTaskRun(),
			"tekton_pipeline":        resourceTektonPipeline(),
			"tekton_pipelinerun":     resourceTektonPipelineRun(),
			"tekton_triggers":        resourceTektonTriggers(),
			"tekton_triggertemplate": resourceTektonTriggerTemplate(),
			"tekton_triggerbinding":  resourceTektonTriggerBinding(),
			"tekton_eventlistener":   resourceTektonEventListener(),
			// Define other resources like "tekton_pipeline" here
		},
		ConfigureFunc: providerConfigure,
	}
}

// providerConfigure sets up the Tekton client for interacting with Tekton resources.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	configPath := d.Get("kubeconfig").(string)

	kubeConfig, err := loadKubeConfig(configPath)
	if err != nil {
		return nil, err
	}

	tektonClient, err := tektonclient.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	tektonTriggersClient, err := triggersclient.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	return struct {
		TektonClient         *tektonclient.Clientset
		TektonTriggersClient *triggersclient.Clientset
	}{
		TektonClient:         tektonClient,
		TektonTriggersClient: tektonTriggersClient,
	}, nil
}

// loadKubeConfig loads the Kubernetes configuration from a file.
func loadKubeConfig(configPath string) (*rest.Config, error) {
	if configPath == "" {
		return nil, fmt.Errorf("KUBECONFIG environment variable or kubeconfig file must be provided")
	}

	// Expand "~" to the user's home directory
	if configPath[:2] == "~/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %v", err)
		}
		configPath = filepath.Join(homeDir, configPath[2:])
	}

	// Build the Kubernetes configuration from the file
	return clientcmd.BuildConfigFromFlags("", configPath)
}
