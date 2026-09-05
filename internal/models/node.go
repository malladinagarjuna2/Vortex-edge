package models

import "time"

type NodeStatus string 

const (
	NodeReady NodeStatus = "Ready"
	NodeUnhealthy NodeStatus ="Unhealthy"
	NodeOffline NodeStatus = "Offline"
)
 
type Node struct {
	 ID  string 
	 Name string 
	 Address string 
	 Status NodeStatus 
	 CPU int
	 Memory int64
	 LastHeartbeat  time.Time
	  CreatedAt time.Time
	UpdatedAt time.Time
}
