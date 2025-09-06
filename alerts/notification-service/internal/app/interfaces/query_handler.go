package interfaces

type QueryHandler interface {
	Execute(query any) (any, error)
}
