package config

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedData(db *gorm.DB) error {
	log.Infof("Seeding initial data...")
	var err error
	err = SeedLeadershipPositions(db)
	err = SeedRoles(db)
	return err
}

func SeedLeadershipPositions(db *gorm.DB) error {
	log.Infof("Seeding leadership positions...")
	var positions []model.ConfigMap

	for k, v := range utils.PositionDisplayMap {
		positions = append(positions, model.ConfigMap{
			Category:    "LEADERSHIP_POSITION",
			Key:         string(k),
			Value:       v,
			AuditFields: model.AuditFields{CreatedBy: "SYSTEM"},
		})
	}

	for _, position := range positions {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}, {Name: "category"}},
			DoNothing: true,
		}).Create(&position).Error; err != nil {
			return err
		}
	}
	return nil
}

func SeedRoles(db *gorm.DB) error {
	log.Infof("Seeding roles...")
	var roles []model.Roles

	for _, roleName := range utils.RolesMap {
		roles = append(roles, model.Roles{
			Name:        roleName,
			AuditFields: model.AuditFields{CreatedBy: "SYSTEM"},
		})
	}

	for _, role := range roles {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).Create(&role).Error; err != nil {
			return err
		}
	}
	return nil
}
