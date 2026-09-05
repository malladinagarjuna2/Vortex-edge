
package runtime

import "testing"

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