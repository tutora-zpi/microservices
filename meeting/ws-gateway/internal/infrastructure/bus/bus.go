package bus

type HandlerRegistry[T any] struct {
	handlers map[string][]T
}

func NewHandlerRegistry[T any]() *HandlerRegistry[T] {
	return &HandlerRegistry[T]{handlers: make(map[string][]T)}
}

func (r *HandlerRegistry[T]) Register(pattern string, handler T) {
	r.handlers[pattern] = append(r.handlers[pattern], handler)
}

func (r *HandlerRegistry[T]) GetHandlers(pattern string) []T {
	return r.handlers[pattern]
}

func (r *HandlerRegistry[T]) Patterns() []string {
	patterns := make([]string, 0, len(r.handlers))
	for p := range r.handlers {
		patterns = append(patterns, p)
	}
	return patterns
}
