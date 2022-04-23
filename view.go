package main

type viewModel struct {
	height int
	width int
}

type view interface {
	setWidth(int)
	setHeight(int)
	render() string
	name() string
}

func (v *viewModel) setHeight(height int) {
	v.height = height
}

func (v *viewModel) setWidth(width int) {
	v.width = width
}