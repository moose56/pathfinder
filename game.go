package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"image"
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
		rect := g.getNodeInScreenSpace(*node)

		if p.In(rect) {
			// add obstacles
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
				node.IsObstacle = !node.IsObstacle
			}

			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if node.IsObstacle {
					continue
				} else if g.grid.PathStart == nil { // set start node
					g.grid.PathStart = node
				} else if g.grid.PathEnd == nil { // set end node
					g.grid.PathEnd = node
				} else if g.grid.PathEnd.Eq(*node) { // unset end node
					g.grid.PathEnd = nil
				} else if g.grid.PathEnd == nil && g.grid.PathStart.Eq(*node) { // unset start node
					g.grid.PathStart = nil
				}
			}
		}
	}

	if g.grid.PathStart != nil && g.grid.PathEnd != nil {
		g.grid.IsSolved = SolveAStar(g.grid)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, node := range g.grid.Nodes() {
		for _, neighbour := range node.Neighbours() {
			p1 := g.getNodeCentreInScreenSpace(*node)
			p2 := g.getNodeCentreInScreenSpace(*neighbour)

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
		r := g.getNodeInScreenSpace(*node)

		colour := colornames.Gray

		if node.IsObstacle {
			colour = colornames.Royalblue
		}

		if g.grid.PathStart != nil && node.Eq(*g.grid.PathStart) {
			colour = colornames.Green
		}

		if g.grid.PathEnd != nil && node.Eq(*g.grid.PathEnd) {
			colour = colornames.Red
		}

		vector.DrawFilledRect(screen, float32(r.Min.X), float32(r.Min.Y), float32(r.Dx()), float32(r.Dy()), colour, false)
	}

	// Draw Path by starting ath the end, and following the parent node trail
	// back to the start - the start node will not have a parent path to follow
	if g.grid.IsSolved {
		p := g.grid.PathEnd

		for p != nil {

			if p.Parent != nil {
				p1 := g.getNodeCentreInScreenSpace(*p)
				p2 := g.getNodeCentreInScreenSpace(*p.Parent)

				vector.StrokeLine(
					screen,
					float32(p1.X),
					float32(p1.Y),
					float32(p2.X),
					float32(p2.Y),
					4,
					colornames.Yellow,
					false)
			}

			r := g.getNodeInScreenSpace(*p)

			vector.DrawFilledRect(
				screen,
				float32(r.Min.X),
				float32(r.Min.Y),
				float32(r.Dx()),
				float32(r.Dy()),
				colornames.Yellow,
				false)

			p = p.Parent
		}
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
