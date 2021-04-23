package api

type API interface {
	ListenAndServeBlocking() error
}
