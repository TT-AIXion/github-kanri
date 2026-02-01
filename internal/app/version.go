package app

func (a App) runVersion() int {
	a.Out.Raw(a.Version)
	return 0
}
