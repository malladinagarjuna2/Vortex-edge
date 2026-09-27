package main

import (
	"log"
	"time"

	"vortex-edge/internal/agent"
	"vortex-edge/internal/models"
	"vortex-edge/internal/runtime"
)

func main() {

	node := &models.Node{
		ID:      "node-01",
		Name:    "worker-01",
		Address: "localhost:9000",
		Status:  models.NodeReady,
		CPU:     8,
		Memory:  16 * 1024 * 1024 * 1024,
	}

	dockerRuntime, err := runtime.NewDockerRuntime()
	if err != nil {
		log.Fatal(err)
	}

	nodeAgent := agent.NewAgent(node, dockerRuntime)

	client, err := agent.NewClient("localhost:9001")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	accepted, err := client.Register(node)
	if err != nil {
		log.Fatal(err)
	}

	if !accepted {
		log.Fatal("node registration rejected")
	}

	log.Println("Node registered successfully")

	go nodeAgent.StartHeartbeat(client, 5*time.Second)

	select {}
}