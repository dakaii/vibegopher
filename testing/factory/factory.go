package factory

import (
	"github.com/dakaii/vibegopher/internal/database"
)

// Save saves any entity to the database using GORM
func Save[T any](entity T) (T, error) {
	db := database.GetDatabase(true)
	tx := db.Begin()
	if tx.Error != nil {
		return entity, tx.Error
	}

	if err := tx.Create(&entity).Error; err != nil {
		tx.Rollback()
		return entity, err
	}

	if err := tx.Commit().Error; err != nil {
		return entity, err
	}

	return entity, nil
}

// SaveMany saves multiple entities to the database using GORM
func SaveMany[T any](entities []T) ([]T, error) {
	db := database.GetDatabase(true)
	tx := db.Begin()
	if tx.Error != nil {
		return entities, tx.Error
	}

	for i := range entities {
		if err := tx.Create(&entities[i]).Error; err != nil {
			tx.Rollback()
			return entities, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return entities, err
	}

	return entities, nil
}
