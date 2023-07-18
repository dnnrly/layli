package layli

import (
	"math"
)

/*
A rough algorithm
1. Create 2d array of possible vertices
x			var verticies [width][height]bool
x			For each x,y: verticies[x][y] = l.InsideAny(x, y) || l.AnyPort(x, y)
x				graph.AddMappedVertex(toMap(x, y))
2. Create "outside" horizontal vertex arcs
x	For each row
x		for src := 1 ; src < width ; src++{
x			for dst := 1 ; dst < width ; dst++{
x				valid := true
x				for p := src ; p <= dst ; p++ {
x					valid = valid && verticies[row][p] {
x				}
x
x				if valid && !l.AnyPort(row, dst) {
x					graph.AddArc(row*width+src, row*width+dst, 1) // Cost of 1 for everything on the same line
x				}
x			}
x		}
3. Create "outside" vertical vertex arcs
x	For each column
x		for src := 1 ; src < height ; src++{
x			for dst := 1 ; dst < width ; dst++{
x				valid := true
x				for p := left ; p <= dst ; p++ {
x					valid = valid && verticies[p][column] {
x				}
x
x				if valid && !l.AnyPort(dst, column) {
x					graph.AddArc(src*width+column, dst*width+column, 1) // Cost of 1 for everything on the same line
x				}
x			}
x		}
4. Add src node arcs - only 1 way (out from the centre of the node)
	Get node centre
	For node port
		Calculate distance between points [d=√((x2 – x1)² + (y2 – y1)²).]
		graph.AddArc(-1, port.X*width+port.Y, distance*100) // *100 to account for float calculations
5. Add dst node arcs- only 1 way (toward the centre of the node)
	Get node centre
	For node port
		Calculate distance between points [d=√((x2 – x1)² + (y2 – y1)²).]
		graph.AddArc(port.X*width+port.Y, -2, distance*100) // *100 to account for float calculations
6. graph.Shortest() -> convert to LayoutPath
*/

type Graph interface {
	AddMappedArc(Source, Destination string, Distance int64) error
	AddMappedVertex(ID string) int
	GetMapped(a int) (string, error)
}

func BuildVertexMap(l *Layout) VertexMap {
	vm := NewVertexMap(l.LayoutWidth(), l.LayoutHeight())
	vm.MapUnset(l.InsideAny)
	vm.MapOr(l.IsAnyPort)

	return vm
}

func (l *Layout) AddPath(from, to string) {
	// g := dijkstra.NewGraph()

	nFrom := l.Nodes.ByID(from)
	nTo := l.Nodes.ByID(to)

	vm := BuildVertexMap(l)
	arcs := vm.GetArcs()

	// vm.GetVertexPoints().AddToGraph(g)

	{
		// Add "from" paths
		centre := nFrom.GetCentre()
		ports := nFrom.GetPorts()
		for _, to := range ports {
			arcs.Add(
				centre,
				to,
				int(math.Sqrt(
					(centre.X-to.X)*(centre.X-to.X)+
						(centre.Y-to.Y)*(centre.Y-to.Y),
				)*100),
			)
		}
		ports = append(ports, centre)
		// ports.AddToGraph(g)
	}

	{
		// Add "to" paths
		centre := nFrom.GetCentre()
		ports := nFrom.GetPorts()
		for _, from := range ports {
			arcs.Add(
				from,
				centre,
				int(math.Sqrt(
					(centre.X-from.X)*(centre.X-from.X)+
						(centre.Y-from.Y)*(centre.Y-from.Y),
				)*100),
			)
		}
		ports = append(ports, centre)
		// ports.AddToGraph(g)
	}

	// g.Shortest()

	l.Paths = append(
		l.Paths,
		LayoutPath{
			points: Points{
				nFrom.GetCentre(),
				nTo.GetCentre(),
			},
		},
	)

	// graph := dijkstra.NewGraph()
}
