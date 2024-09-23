terraform provider for tekton

go get github.com/hashicorp/terraform-plugin-sdk/v2
go get github.com/tektoncd/pipeline/pkg/client/clientset/versioned

## Provider
```
terraform {
  required_providers {
    tekton = {
      source  = "SamuelDBines/tekton"
      version = "1.0.16"
    }
    ...
  }
}

provider "tekton" {
  kubeconfig = "~/.kube/config"
}
```

## Tasks
```
resource "tekton_task" "hello_task" {
  name = "hello-task"
  namespace = "default"
  steps {
    name    = "echo"
    image   = "alpine"
    command = ["echo", "Hello, World!"]
  }
}
```

## Tasks Runs
```
resource "tekton_taskrun" "hello_taskrun" {
  name                  = "example-taskrun"
  namespace             = "default"
  task_ref_name         = tekton_task.hello_task.name
  service_account_name  = "default"
  
  params {
    name  = "example-param"
    value = "Hello, World!"
  }
}
```

## Pipeline

```
resource "tekton_pipeline" "example_pipeline" {
  name      = "example-pipeline"
  namespace = "default"

  tasks {
    name          = "task-1"
    task_ref_name = "example-task"  # This should reference an existing Tekton Task
  }

  tasks {
    name          = "task-2"
    task_ref_name = "example-task-2"
    run_after     = ["task-1"]
  }
}

resource "tekton_pipelinerun" "example_pipelinerun" {
  name                  = "example-pipelinerun"
  namespace             = "default"
  pipeline_ref_name     = tekton_pipeline.example_pipeline.name  # Reference the pipeline created above
  service_account_name  = "default"

  params {
    name  = "example-param"
    value = "Hello, World!"
  }
}
```

## Triggers

```

resource "tekton_triggertemplate" "my_template" {
  name      = "my-template"
  namespace = "default"

  params {
    name        = "param1"
    description = "A parameter for the pipeline"
  }

  resourcetemplates {
    api_version = "tekton.dev/v1beta1"
    kind        = "PipelineRun"
    metadata = {
      name = "example-pipelinerun"
    }
    spec = {
      pipelineRef = {
        name = "example-pipeline"
      }
      params = [{
        name = "param1"
        value = "$(params.param1)"
      }]
    }
  }
}

resource "tekton_triggerbinding" "my_binding" {
  name      = "my-binding"
  namespace = "default"

  bindings {
    name  = "param1"
    value = "$(body.param1)"
  }
}

resource "tekton_eventlistener" "my_listener" {
  name      = "my-listener"
  namespace = "default"

  triggers {
    trigger_template_name = tekton_triggertemplate.my_template.name
    trigger_binding_name  = tekton_triggerbinding.my_binding.name
  }
}
```