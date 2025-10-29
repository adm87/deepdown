package input

type InputAction uint8

type InputBinding interface {
	Action() InputAction
	Update(dt float64)
	IsActive() bool
}

type Binding struct {
	action InputAction
	active bool
}

func NewBinding(action InputAction) Binding {
	return Binding{action: action}
}

func (b *Binding) Action() InputAction {
	return b.action
}

func (b *Binding) IsActive() bool {
	return b.active
}

func (b *Binding) State(action InputAction) int {
	return 0
}
