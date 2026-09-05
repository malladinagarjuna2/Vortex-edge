

package agent 
import (
	"vortex-edge/internal/models"
	"vortex-edge/internal/runtime"
)

type Agent struct {
	 node *models.Node
	 runtime runtime.Runtime
}

func NewAgent(node *models.Node, runtime runtime.Runtime)*Agent{
	  return &Agent{
		node: node, 
		runtime: runtime,
	  }
}
func(a*Agent)Node() *models.Node{
	 return a.node
}
