package main

type Grid struct {
	data   []Node
	width  int
	height int
}

func NewGrid(width, height int) *Grid {
	g := &Grid{
		data:   make([]Node, width*height),
		width:  width,
		height: height,
	}

	// set x and y for each node in grid
	for x := 0; x < g.Width(); x++ {
		for y := 0; y < g.Height(); y++ {
			n := g.Get(x, y)
			n.x, n.y = x, y
			g.Put(x, y, n)
		}
	}

	// add connections for each node in grid
	for i, node := range g.data {
		if !g.IsTopEdge(node) {
			node.AddNeighbour(g.North(node))
		}

		if !g.IsTopEdge(node) && !g.IsRightEdge(node) {
			node.AddNeighbour(g.NorthEast(node))
		}

		if !g.IsRightEdge(node) {
			node.AddNeighbour(g.East(node))
		}

		if !g.IsBottomEdge(node) && !g.IsRightEdge(node) {
			node.AddNeighbour(g.SouthEast(node))
		}

		if !g.IsBottomEdge(node) {
			node.AddNeighbour(g.South(node))
		}

		if !g.IsBottomEdge(node) && !g.IsLeftEdge(node) {
			node.AddNeighbour(g.SouthWest(node))
		}

		if !g.IsLeftEdge(node) {
			node.AddNeighbour(g.West(node))
		}

		if !g.IsTopEdge(node) && !g.IsLeftEdge(node) {
			node.AddNeighbour(g.NorthWest(node))
		}

		// update node in the grid
		g.data[i] = node
	}
	return g
}

func (g *Grid) Width() int {
	return g.width
}
func (g *Grid) Height() int {
	return g.height
}
func (g *Grid) Nodes() []Node { return g.data }

func (g *Grid) Put(x, y int, value Node) {
	g.data[x*g.Height()+y] = value
}
func (g *Grid) Get(x, y int) Node {
	return g.data[x*g.Height()+y]
}

func (g *Grid) IsTopLeftCorner(node Node) bool {
	return node.X() == 0 && node.Y() == 0
}
func (g *Grid) IsTopRightCorner(node Node) bool {
	return node.X() == g.Width()-1 && node.Y() == 0
}
func (g *Grid) IsBottomLeftCorner(node Node) bool {
	return node.X() == 0 && node.Y() == g.Height()-1
}
func (g *Grid) IsBottomRightCorner(node Node) bool {
	return node.X() == g.Width()-1 && node.Y() == g.Height()-1
}

func (g *Grid) IsTopEdge(node Node) bool {
	return node.Y() == 0
}
func (g *Grid) IsBottomEdge(node Node) bool {
	return node.Y() == g.Height()-1
}
func (g *Grid) IsLeftEdge(node Node) bool {
	return node.X() == 0
}
func (g *Grid) IsRightEdge(node Node) bool {
	return node.X() == g.Width()-1
}

func (g *Grid) IsWithinGrid(x, y int) bool {
	return x >= 0 && x < g.width && y >= 0 && y < g.height
}

func (g *Grid) North(node Node) Node {
	return g.Get(node.X(), node.Y()-1)
}
func (g *Grid) NorthEast(node Node) Node {
	return g.Get(node.X()+1, node.Y()-1)
}
func (g *Grid) East(node Node) Node {
	return g.Get(node.X()+1, node.Y())
}
func (g *Grid) SouthEast(node Node) Node {
	return g.Get(node.X()+1, node.Y()+1)
}
func (g *Grid) South(node Node) Node {
	return g.Get(node.X(), node.Y()+1)
}
func (g *Grid) SouthWest(node Node) Node {
	return g.Get(node.X()-1, node.Y()+1)
}
func (g *Grid) West(node Node) Node {
	return g.Get(node.X()-1, node.Y())
}
func (g *Grid) NorthWest(node Node) Node {
	return g.Get(node.X()-1, node.Y()-1)
}

type Node struct {
	x          int
	y          int
	neighbours []Node
	Parent     *Node
}

func (n *Node) X() int {
	return n.x
}
func (n *Node) Y() int {
	return n.y
}
func (n *Node) Neighbours() []Node {
	return n.neighbours
}

func (n *Node) AddNeighbour(node Node) {
	n.neighbours = append(n.neighbours, node)
}
func (n *Node) Eq(node Node) bool {
	return n.x == node.x && n.y == node.y
}
