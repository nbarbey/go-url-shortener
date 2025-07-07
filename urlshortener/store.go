package urlshortener

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
import "gorm.io/driver/sqlite"

type Storer interface {
	Get(shortened string) (URL, error)
	Save(url, shortened string, expiration *time.Time) error
}

type InMemoryStore struct {
	data map[string]string
}

func NewStore() *InMemoryStore {
	return &InMemoryStore{data: make(map[string]string)}
}

func (s *InMemoryStore) Get(shortened string) (URL, error) {
	u, ok := s.data[shortened]
	if !ok {
		return URL{}, ErrNotFound
	}
	return NewURL(u, nil)
}

func (s *InMemoryStore) Save(url, shortened string) error {
	s.data[shortened] = url
	return nil
}

type PGStore struct {
	db *gorm.DB
}

type URLAssociation struct {
	URL        string
	Shortened  string `gorm:"primaryKey"`
	Expiration sql.NullTime
}

func (p PGStore) Get(shortened string) (URL, error) {
	var association = URLAssociation{}
	tx := p.db.First(&association, "shortened = ?", shortened)
	if tx.Error != nil {
		return URL{}, tx.Error
	}
	if association.Expiration.Valid {
		location, err := time.LoadLocation("Local")
		if err != nil {
			return URL{}, err
		}
		utc := association.Expiration.Time.In(location)
		return NewURL(association.URL, &utc)
	}
	return NewURL(association.URL, nil)
}

func (p PGStore) Save(url, shortened string, expiration *time.Time) error {
	var t sql.NullTime
	if expiration == nil {
		t.Valid = false
	} else {
		t.Valid = true
		t.Time = *expiration
	}
	tx := p.db.Create(&URLAssociation{URL: url, Shortened: shortened, Expiration: t})
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

func NewPG() *PGStore {
	dsn := fmt.Sprintf("host=%s dbname=%s port=5432 user=%s password=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	_ = db.AutoMigrate(&URLAssociation{})
	return &PGStore{db: db}
}
