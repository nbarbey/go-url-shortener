package urlshortener

type Application struct {
	store  Storer
	server *HTTPServer
	*CountingUsecase
}

func (a *Application) Start() error {
	return a.server.Start()
}

func NewInMemoryApplication() *Application {
	store := NewInMemorySqlite()
	countStore := NewInMemoryCountStore()
	usecase := &CountingUsecase{Usecase: NewUsecase(store), countStore: countStore}
	return &Application{
		store:           store,
		CountingUsecase: usecase,
		server:          NewHTTPServer(usecase, countStore)}
}

func NewPGpplication() *Application {
	store := NewPG()
	countStore := NewPGCountStore()
	usecase := &CountingUsecase{Usecase: NewUsecase(store), countStore: countStore}
	return &Application{
		store:           store,
		CountingUsecase: usecase,
		server:          NewHTTPServer(usecase, countStore)}
}
