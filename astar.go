package main

import (
	"cmp"
	"math"
	"slices"
)

// Algorithm taken from https://github.com/OneLoneCoder/Javidx9/blob/master/ConsoleGameEngine/SmallerProjects/OneLoneCoder_PathFinding_AStar.cpp
func SolveAStar(grid *Grid) bool {
	// Reset nodes
	for _, node := range grid.Nodes() {
		node.IsVisited = false
		node.GlobalGoal = float32(math.Inf(1))
		node.LocalGoal = float32(math.Inf(1))
		node.Parent = nil
	}

	// Setup starting conditions
	nodeCurrent := grid.PathStart
	grid.PathStart.LocalGoal = 0
	grid.PathStart.GlobalGoal = heuristic(grid.PathStart, grid.PathEnd)

	// Add start node to not tested list - this will ensure it gets tested.
	// As the algorithm progresses, newly discovered nodes get added to this
	// list, and will themselves be tested later
	listNotTestedNodes := make([]*Node, 0)
	listNotTestedNodes = append(listNotTestedNodes, grid.PathStart)

	// if the not tested list contains nodes, there may be better paths
	// which have not yet been explored. However, we will also stop
	// searching when we reach the target - there may well be better
	// paths but this one will do - it won't be the longest.
	for len(listNotTestedNodes) > 0 && !nodeCurrent.Eq(*grid.PathEnd) { // Find absolutely shortest path // && nodeCurrent != nodeEnd)
		// Sort Untested nodes by global goal, so lowest is first
		slices.SortFunc(listNotTestedNodes, func(a, b *Node) int {
			return cmp.Compare(a.GlobalGoal, b.GlobalGoal)
		})

		// Front of listNotTestedNodes is potentially the lowest distance node. Our
		// list may also contain nodes that have been visited, so ditch these...
		for len(listNotTestedNodes) > 0 && listNotTestedNodes[0].IsVisited {
			listNotTestedNodes = listNotTestedNodes[1:]
		}

		// ...or abort because there are no valid nodes left to test
		if len(listNotTestedNodes) == 0 {
			break
		}

		nodeCurrent = listNotTestedNodes[0]
		nodeCurrent.IsVisited = true // We only explore a node once

		// Check each of this node's neighbours...
		for _, nodeNeighbour := range nodeCurrent.neighbours {
			// ... and only if the neighbour is not visited and is
			// not an obstacle, add it to NotTested List
			if !nodeNeighbour.IsVisited && !nodeNeighbour.IsObstacle {
				listNotTestedNodes = append(listNotTestedNodes, nodeNeighbour)
			}

			// Calculate the neighbours potential lowest parent distance
			possiblyLowerGoal := nodeCurrent.LocalGoal + distance(nodeCurrent, nodeNeighbour)

			// If choosing to path through this node is a lower distance than what
			// the neighbour currently has set, update the neighbour to use this node
			// as the path source, and set its distance scores as necessary
			if possiblyLowerGoal < nodeNeighbour.LocalGoal {
				nodeNeighbour.Parent = nodeCurrent
				nodeNeighbour.LocalGoal = possiblyLowerGoal

				// The best path length to the neighbour being tested has changed, so
				// update the neighbour's score. The heuristic is used to globally bias
				// the path algorithm, so it knows if its getting better or worse. At some
				// point the algo will realise this path is worse and abandon it, and then go
				// and search along the next best path.
				nodeNeighbour.GlobalGoal = nodeNeighbour.LocalGoal + heuristic(nodeNeighbour, grid.PathEnd)
			}
		}
	}
	return true
}

func distance(a, b *Node) float32 {
	return float32(math.Sqrt(float64((a.X()-b.X())*(a.X()-b.X()) + (a.Y()-b.Y())*(a.Y()-b.Y()))))
}

func heuristic(a, b *Node) float32 {
	return distance(a, b)
}
