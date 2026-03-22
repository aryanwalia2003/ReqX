package dag

import (
	"fmt"
	"strings"
)

// TopoSort runs Kahn's algorithm on g and returns the nodes grouped into levels.
// All nodes in the same level have no dependency on each other and can be
// executed concurrently. Levels must be executed in order.
//
// Returns an error if the graph contains a cycle.
func TopoSort(g *ScenarioGraph) ([][]int, error) {
	n := len(g.Nodes)

	// inDegree[i] = number of nodes that i is waiting for.
	inDegree := make([]int, n)
	for i := 0; i < n; i++ {
		inDegree[i] = len(g.Edges[i])
	}

	// adjacency[i] = nodes that become unblocked when i finishes.
	adjacency := make([][]int, n)
	for i := 0; i < n; i++ {
		for _, dep := range g.Edges[i] {
			adjacency[dep] = append(adjacency[dep], i)
		}
	}

	var levels [][]int
	processed := 0

	// Seed with nodes that have no dependencies.
	current := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if inDegree[i] == 0 {
			current = append(current, i)
		}
	}

	for len(current) > 0 {
		levels = append(levels, current)
		processed += len(current)

		next := make([]int, 0)
		for _, nodeIdx := range current {
			for _, dependent := range adjacency[nodeIdx] {
				inDegree[dependent]--
				if inDegree[dependent] == 0 {
					next = append(next, dependent)
				}
			}
		}
		current = next
	}

	if processed != n {
		return nil, cycleError(g, inDegree)
	}

	return levels, nil
}

// cycleError builds a human-readable error listing the nodes involved in the cycle.
func cycleError(g *ScenarioGraph, inDegree []int) error {
	involved := make([]string, 0)
	for i, deg := range inDegree {
		if deg > 0 {
			involved = append(involved, fmt.Sprintf("%q", g.Nodes[i].Name))
		}
	}
	return fmt.Errorf(
		"dag: cycle detected among requests: [%s]",
		strings.Join(involved, ", "),
	)
}