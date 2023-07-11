package layli

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
}
