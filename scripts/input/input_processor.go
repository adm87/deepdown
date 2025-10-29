package input

var bindings = make(map[InputAction]InputBinding, 10)

func Register(newBindings ...InputBinding) {
	for _, b := range newBindings {
		bindings[b.Action()] = b
	}
}

func Update(dt float64) {
	for _, b := range bindings {
		b.Update(dt)
	}
}

func IsActive(action InputAction) bool {
	if b, ok := bindings[action]; ok {
		return b.IsActive()
	}
	return false
}

func GetBinding[T InputBinding](action InputAction) T {
	if b, ok := bindings[action]; ok {
		if binding, ok := b.(T); ok {
			return binding
		}
	}
	var zero T
	return zero
}
