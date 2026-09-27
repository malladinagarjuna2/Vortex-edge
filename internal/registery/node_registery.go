package registery 

import (
	"sync"
	"vortex-edge/internal/models"
)

type NodeRegistry struct {
	 nodes map[string]*models.Node
	 mu sync.RWMutex
}

 func NewNodeRegistry()*NodeRegistry{
	 return &NodeRegistry{
		 nodes: make(map[string]*models.Node),
	 }
 }

 func(r*NodeRegistry)Add(node *models.Node)error{
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes[node.ID] = node
	return nil 
 }



func( r*NodeRegistry)Get(id string)( *models.Node, error){
	 r.mu.RLock()
	defer r.mu.RUnlock()
	node, exists := r.nodes[id]
	if !exists{
		 return nil, nil
	}
	 return node, nil
}

 func (r*NodeRegistry)Update(node *models.Node)error{
	 r.mu.Lock()
	 defer r.mu.Unlock()
     r.nodes[node.ID]= node
	 return nil

 }

 func (r *NodeRegistry) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.nodes, id)

	return nil
}

func (r *NodeRegistry) List() []*models.Node {
	r.mu.RLock()
	defer r.mu.RUnlock()

	nodes := make([]*models.Node, 0, len(r.nodes))

	for _, node := range r.nodes {
		nodes = append(nodes, node)
	}

	return nodes
}

    //       New Machine
    //           │
    //           │ Node Agent starts
    //           ▼
    //     ┌─────────────┐
    //     │ Node Agent  │
    //     └──────┬──────┘
    //            │
    //            │ Register()
    //            ▼
    //     ┌─────────────┐
    //     │   Control   │
    //     │    Plane    │
    //     └──────┬──────┘
    //            │
    //            ▼
    //      Node Registry
    //            │
    //    ┌───────┴───────┐
    //    ▼               ▼
    // node-01         node-02