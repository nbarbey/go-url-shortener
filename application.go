package url_shortener

type Application struct {
	store  Storer
	server *HTTPServer
	*Usecase
}

func (a *Application) Start() error {
	return a.server.Start()
}

func NewApplication() *Application {
	store := NewInMemorySqlite()
	usecase := NewUsecase(store)
	return &Application{
		store:   store,
		Usecase: usecase,
		server:  NewHTTPServer(usecase)}
}
