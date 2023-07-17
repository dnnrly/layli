package layli

/*
A rough algorithm
1. Create 2d array of possible vertices
x			var verticies [width][height]bool
x			For each x,y: verticies[x][y] = l.InsideAny(x, y) || l.AnyPort(x, y)
x				graph.AddMappedVertex(toMap(x, y))
2. Create "outside" horizontal vertex arcs
	For each row
		for src := 1 ; src < width ; src++{
			for dst := 1 ; dst < width ; dst++{
				valid := true
				for p := src ; p <= dst ; p++ {
					valid = valid && verticies[row][p] {
				}

				if valid && !l.AnyPort(row, dst) {
					graph.AddArc(row*width+src, row*width+dst, 1) // Cost of 1 for everything on the same line
				}
			}
		}
3. Create "outside" vertical vertex arcs
	For each column
		for src := 1 ; src < height ; src++{
			for dst := 1 ; dst < width ; dst++{
				valid := true
				for p := left ; p <= dst ; p++ {
					valid = valid && verticies[p][column] {
				}

				if valid && !l.AnyPort(dst, column) {
					graph.AddArc(src*width+column, dst*width+column, 1) // Cost of 1 for everything on the same line
				}
			}
		}
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
}

func BuildVertexMap(l *Layout) VertexMap {
	vm := NewVertexMap(l.LayoutWidth(), l.LayoutHeight())
	vm.MapUnset(l.InsideAny)
	vm.MapOr(l.IsAnyPort)

	return vm
}

func (l *Layout) AddPath(from, to string) {
	nFrom := l.Nodes.ByID(from)
	nTo := l.Nodes.ByID(to)

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
