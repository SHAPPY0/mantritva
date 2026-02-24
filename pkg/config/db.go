package config

import (
	"log"
	"strings"

	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase() *gorm.DB {
	var dsnBuilder strings.Builder
	dsnBuilder.WriteString("host=")
	dsnBuilder.WriteString(GetConfig().GetString("DB_HOST"))
	dsnBuilder.WriteString(" user=")
	dsnBuilder.WriteString(GetConfig().GetString("DB_USER"))
	dsnBuilder.WriteString(" dbname=")
	dsnBuilder.WriteString(GetConfig().GetString("DB_NAME"))
	dsnBuilder.WriteString(" sslmode=disable password=")
	dsnBuilder.WriteString(GetConfig().GetString("DB_PASSWORD"))
	var dns = dsnBuilder.String()
	log.Printf("Postgres DSN: %s", dns)
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	//Auto Migrate
	db.AutoMigrate(&model.Country{})
	db.AutoMigrate(&model.State{})
	db.AutoMigrate(&model.District{})
	db.AutoMigrate(&model.Loksabha{})
	db.AutoMigrate(&model.Vidhansabha{})
	db.AutoMigrate(&model.Party{})
	db.AutoMigrate(&model.Alliance{})
	db.AutoMigrate(&model.AllianceParty{})
	db.AutoMigrate(&model.Ruler{})
	db.AutoMigrate(&model.Member{})
	db.AutoMigrate(&model.CriminalRecord{})
	db.AutoMigrate(&model.Career{})
	db.AutoMigrate(&model.PartyMembers{})
	db.AutoMigrate(&model.Leadership{})
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Entity{})
	db.AutoMigrate(&model.KSA{})
	db.AutoMigrate(&model.Questionnaire{})
	db.AutoMigrate(&model.Enquiry{})

	return db
}
