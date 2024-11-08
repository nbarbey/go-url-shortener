package urlshortener

import "gorm.io/gorm"
import "gorm.io/driver/sqlite"

type Storer interface {
	Get(shortened string) (string, error)
	Save(url, shortened string) error
}

type InMemoryStore struct {
	data map[string]string
}

func NewStore() *InMemoryStore {
	return &InMemoryStore{data: make(map[string]string)}
}

func (s *InMemoryStore) Get(shortened string) (string, error) {
	u, ok := s.data[shortened]
	if !ok {
		return "", ErrNotFound
	}
	return u, nil
}

func (s *InMemoryStore) Save(url, shortened string) error {
	s.data[shortened] = url
	return nil
}

type PGStore struct {
	db *gorm.DB
}

type URLAssociation struct {
	URL       string
	Shortened string `gorm:"primaryKey"`
}

func (p PGStore) Get(shortened string) (string, error) {
	var association = URLAssociation{}
	tx := p.db.First(&association, "shortened = ?", shortened)
	return association.URL, tx.Error
}

func (p PGStore) Save(url, shortened string) error {
	tx := p.db.Create(&URLAssociation{URL: url, Shortened: shortened})
	return tx.Error
}

func NewInMemorySqlite() *PGStore {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	_ = db.AutoMigrate(&URLAssociation{})
	return &PGStore{db: db}
}
