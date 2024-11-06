package url_shortener

type Application struct {
	store Storer
}

func NewApplication() *Application {
	store := NewStore()
	store.Save("https://medium.com/equify-tech/the-three-fundamental-stages-of-an-engineering-career-54dac732fc74",
		"https://localhost/hardcoded")
	return &Application{store: NewInMemorySqlite()}
}
