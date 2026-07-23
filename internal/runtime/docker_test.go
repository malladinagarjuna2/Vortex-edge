// package runtime
// import "testing"

// // func TestNewDockerRuntime_nonNil(t *testing.T){
// //     runtime, err :=  TestNewDockerRuntime()
// //     if err != nil {
// //          t.Fatalf("failed to create Docker runtime %v", err)
// //     } 
// //     if runtime == nil {
// //          t.Fatal("expected runtime, got nil")
// //     }

// //     if runtime.client== nil {
// //         t.Fatal("expected Docker client to be initialised")
// //     }


// // }
// /*What this tests

// It verifies that:

// ✅ NewDockerRuntime() doesn't return an error.
// ✅ It returns a non-nil DockerRuntime.
// ✅ The Docker client inside it is initialized.*/