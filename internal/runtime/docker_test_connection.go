package  runtime
 import "testing"

 func TestNewDockerRuntime(t *testing.T){
	runtime, err := TestNewDockerRuntime()
	if err != nil {
		 t.Fatalf("failed to create Docker runtime %v", err)
	}
	if err := runtime.Ping(); err != nil {
		 t.Fatalf("failed to ping Docker daemon %v", err)
	}
}

/*This test verifies that:

Go
 ↓
Docker SDK
 ↓
Docker Daemon

is working.*/