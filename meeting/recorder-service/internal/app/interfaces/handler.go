package interfaces

type UseCaseHandler interface {
	//Executes logic on body
	Exec(body []byte) error
}
