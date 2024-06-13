package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"image"
	"slices"
)

var (
	nilNode   = Node{x: -1, y: -1}
	pathStart = nilNode
	pathEnd   = nilNode
	obstacles []Node
)

type Game struct {
	nodeSize int
	gap      int
	grid     *Grid
	width    int
	height   int
}

func NewGame(grid *Grid, nodeSize, gap int) *Game {
	var width = (nodeSize * grid.Width()) + (gap * (grid.Width() + 1))
	var height = (nodeSize * grid.Height()) + (gap * (grid.Height() + 1))

	return &Game{
		nodeSize: nodeSize,
		gap:      gap,
		grid:     grid,
		width:    width,
		height:   height,
	}
}

func (g *Game) getNodeInScreenSpace(node Node) image.Rectangle {
	xGap := (node.X() + 1) * g.gap
	yGap := (node.Y() + 1) * g.gap
	x0 := node.X()*g.nodeSize + xGap
	y0 := node.Y()*g.nodeSize + yGap

	return image.Rect(x0, y0, x0+g.nodeSize, y0+g.nodeSize)
}

func (g *Game) getNodeCentreInScreenSpace(node Node) image.Point {
	xGap := (node.X() + 1) * g.gap
	yGap := (node.Y() + 1) * g.gap
	x := ((node.X() * g.nodeSize) + xGap) + g.nodeSize/2
	y := ((node.Y() * g.nodeSize) + yGap) + g.nodeSize/2

	return image.Point{X: x, Y: y}
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	p := image.Point{X: x, Y: y}

	for _, node := range g.grid.Nodes() {
		rect := g.getNodeInScreenSpace(node)

		if p.In(rect) {
			// add obstacles
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {

				// get the index of the current node in obstacles
				index := slices.IndexFunc(obstacles, func(n Node) bool {
					return n.Eq(node)
				})

				// not there add it or if there delete it
				if index == -1 {
					obstacles = append(obstacles, node)
				} else {
					obstacles = slices.Delete(obstacles, index, index+1)
				}
			}

			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if slices.ContainsFunc(obstacles, func(n Node) bool {
					return n.Eq(node)
				}) {
					continue
				} else if pathStart.Eq(nilNode) { // set start node
					pathStart = node
				} else if pathEnd.Eq(node) { // unset end node
					pathEnd = nilNode
				} else if pathStart.Eq(node) && pathEnd.Eq(nilNode) { // unset start node
					pathStart = nilNode
				} else if !pathStart.Eq(nilNode) && pathEnd.Eq(nilNode) { // set end node
					pathEnd = node
				}
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, node := range g.grid.Nodes() {
		for _, neighbour := range node.Neighbours() {
			p1 := g.getNodeCentreInScreenSpace(node)
			p2 := g.getNodeCentreInScreenSpace(neighbour)

			vector.StrokeLine(
				screen,
				float32(p1.X),
				float32(p1.Y),
				float32(p2.X),
				float32(p2.Y),
				4,
				colornames.Gray,
				false)
		}

	}

	for _, node := range g.grid.Nodes() {
		r := g.getNodeInScreenSpace(node)

		colour := colornames.Gray

		if slices.ContainsFunc(obstacles, func(n Node) bool {
			return n.Eq(node)
		}) {
			colour = colornames.Royalblue
		}

		if node.Eq(pathStart) {
			colour = colornames.Green
		}

		if node.Eq(pathEnd) {
			colour = colornames.Red
		}

		vector.DrawFilledRect(screen, float32(r.Min.X), float32(r.Min.Y), float32(r.Dx()), float32(r.Dy()), colour, false)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	screenWidth = g.width
	screenHeight = g.height

	return screenWidth, screenHeight
}

func (g *Game) Width() int {
	return g.width
}

func (g *Game) Height() int {
	return g.height
}
