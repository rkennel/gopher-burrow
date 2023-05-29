package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type DBContext struct {
	DB *gorm.DB
}

type BasicEntity interface {
	GetId() uuid.UUID
}

type BasicFields struct {
	ID        uuid.UUID `gorm:"type:uuid;primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BasicRepository[T BasicEntity] interface {
	FindByID(id uuid.UUID) (*T, error)
	FindAll() ([]T, error)
	Save(t *T) (*T, error)
	Delete(id uuid.UUID) error
	Exists(id uuid.UUID) (bool, error)
	Count() (int64, error)
}

func NewBasicRepository[T BasicEntity](serverContext *DBContext) BasicRepository[T] {
	return basicRepositoryImpl[T]{ServerContext: serverContext}
}

type basicRepositoryImpl[T BasicEntity] struct {
	ServerContext *DBContext
}

func (repo basicRepositoryImpl[T]) FindByID(id uuid.UUID) (*T, error) {
	var basic T
	result := repo.ServerContext.DB.First(&basic, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &basic, nil
}

func (repo basicRepositoryImpl[T]) FindAll() ([]T, error) {
	var basics []T
	result := repo.ServerContext.DB.Find(&basics)
	if result.Error != nil {
		return nil, result.Error
	}
	return basics, nil
}

func (repo basicRepositoryImpl[T]) Save(basic *T) (*T, error) {
	result := repo.ServerContext.DB.Save(basic)
	if result.Error != nil {
		return nil, result.Error
	}
	return basic, nil
}

func (repo basicRepositoryImpl[T]) Delete(id uuid.UUID) error {
	emptyObject := new(T)
	tx := repo.ServerContext.DB.Delete(emptyObject, id)
	return tx.Error
}

func (repo basicRepositoryImpl[T]) Exists(id uuid.UUID) (bool, error) {
	var basic T
	result := repo.ServerContext.DB.First(&basic, id)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (repo basicRepositoryImpl[T]) Count() (int64, error) {
	emptyObject := new(T)
	var count int64
	tx := repo.ServerContext.DB.Model(&emptyObject).Count(&count)
	return count, tx.Error
}

func (basic *BasicFields) BeforeCreate(tx *gorm.DB) (err error) {
	ts := time.Now()
	basic.ID = uuid.New()
	basic.CreatedAt = ts
	basic.UpdatedAt = ts

	return nil
}

func (basic *BasicFields) BeforeUpdate(tx *gorm.DB) (err error) {
	ts := time.Now()
	basic.UpdatedAt = ts

	return nil
}
