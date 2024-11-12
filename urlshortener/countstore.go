package urlshortener

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

type CountStorer interface {
	Increment(url string) error
	Get(url string) (int, error)
}

type PGCountStore struct {
	db *gorm.DB
}

func (pcs *PGCountStore) Increment(url string) error {
	return pcs.db.Transaction(func(tx *gorm.DB) error {
		hits, err := (&PGCountStore{db: tx}).Get(url)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Create(&CountStoreRow{URL: url, Hits: 1}).Error
		}
		return tx.Model(&CountStoreRow{URL: url}).Update("hits", hits+1).Error
	})
}

func (pcs *PGCountStore) Get(url string) (int, error) {
	var row CountStoreRow
	tx := pcs.db.First(&row, "url = ?", url)
	return row.Hits, tx.Error
}

type CountStoreRow struct {
	URL  string `gorm:"primaryKey"`
	Hits int
}

func NewInMemoryCountStore() *PGCountStore {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&CountStoreRow{})
	if err != nil {
		panic("failed to migrate to schema")
	}
	return &PGCountStore{db: db}
}

func NewPGCountStore() *PGCountStore {
	dsn := fmt.Sprintf("host=%s dbname=%s port=5432 user=%s password=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	_ = db.AutoMigrate(&CountStoreRow{})
	return &PGCountStore{db: db}
}
