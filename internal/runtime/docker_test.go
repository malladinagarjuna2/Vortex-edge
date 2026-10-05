
package runtime

import( "testing"
"vortex-edge/internal/models"
)
func TestDockerConnection(t*testing.T){
	runtime, err:= NewDockerRuntime()
	if err!= nil {
		t.Fatalf("failed to create Docker runtime: %v", err)
	}
	 version, err:= runtime.ServerVersion()
	 if err!= nil {
	t.Fatalf("failed to connect to Docker: %v", err)
	}

	t.Logf("Connected to Docker Engine version: %s", version)
}

func TestDockerStop(t *testing.T) {
	runtime, err := NewDockerRuntime()
	if err != nil {
		t.Fatalf("failed to create docker runtime: %v", err)
	}

	service := &models.Service{
		ID:     "stop-test",
		Name:   "vortex-stop-test",
		Image:  "hello-world",
		Status: models.Deploying,
	}

	err = runtime.Run(service)
	if err != nil {
		t.Fatalf("failed to run container: %v", err)
	}

	t.Logf("Started container: %s", service.ContainerID)

	err = runtime.Stop(service)
	if err != nil {
		t.Fatalf("failed to stop container: %v", err)
	}

	if service.Status != models.Stopped {
		t.Fatalf(
			"expected Stopped, got %s",
			service.Status,
		)
	}

	t.Logf(
		"Container stopped successfully: %s",
		service.ContainerID,
	)
}