package helpers

import "gorm.io/gorm"

func CommonFilter(query *gorm.DB, filter map[string]interface{}) *gorm.DB {
	query.Where("deleted_at IS NULL")
	if id, ok := filter["id"].(string); ok {
		query = query.Where("id = ?", id)
	}
	return query
}
