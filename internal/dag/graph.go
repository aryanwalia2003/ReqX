package dag

import "fmt"

// Node represents a single request inside the scenario graph.
// Index maps 1-to-1 with the position of the request in ExecutionPlan.Requests.
type Node struct {
	Index int
	Name  string
}

// ScenarioGraph is a directed acyclic graph where each node is a request and
// each directed edge means "target depends on source completing first".
type ScenarioGraph struct {
	Nodes []Node
	// Edges[i] holds the indices of all nodes that node i must wait for.
	Edges map[int][]int
	// nameToIndex is used during build to resolve "depends_on" name strings.
	nameToIndex map[string]int
}

// Build constructs a ScenarioGraph from a slice of request names and their
// dependency lists. names[i] is the name of request i; deps[i] is the slice of
// request names that request i depends on.
//
// Returns nil, nil when no request declares any dependency — callers use this
// signal to fall back to the existing linear runner.
func Build(names []string, deps [][]string) (*ScenarioGraph, error) {
	hasDeps := false
	for _, d := range deps {
		if len(d) > 0 {
			hasDeps = true
			break
		}
	}
	if !hasDeps {
		return nil, nil
	}

	g := &ScenarioGraph{
		Nodes:       make([]Node, len(names)),
		Edges:       make(map[int][]int, len(names)),
		nameToIndex: make(map[string]int, len(names)),
	}

	for i, name := range names {
		g.Nodes[i] = Node{Index: i, Name: name}
		g.nameToIndex[name] = i
	}

	for i, depNames := range deps {
		for _, depName := range depNames {
			depIdx, ok := g.nameToIndex[depName]
			if !ok {
				return nil, fmt.Errorf(
					"dag: request %q declares depends_on %q but no such request exists in the collection",
					names[i], depName,
				)
			}
			g.Edges[i] = append(g.Edges[i], depIdx)
		}
	}

	return g, nil
}