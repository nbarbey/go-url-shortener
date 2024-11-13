package urlshortener

type Application struct {
	server *HTTPServer
	*CountingUsecase
}

func (a *Application) Start() error {
	return a.server.Start()
}

func NewInMemoryApplication() *Application {
	return NewApplicationFromInfrastructure(NewInMemoryInfrastructure())
}

func NewPGpplication() *Application {
	return NewApplicationFromInfrastructure(NewPGInfrastructure())
}

func NewApplicationFromInfrastructure(i *InfraStructure) *Application {
	useCases := NewCountingUsecase(i.store, i.countStore)
	return &Application{
		CountingUsecase: useCases,
		server:          NewHTTPServer(useCases, i.countStore),
	}
}
