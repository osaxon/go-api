package config

import (
	"go-api/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func GenDal() {
	// Initialize the generator with configuration
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./tmp/dal", // output directory, default value is ./query
		Mode:          gen.WithQueryInterface | gen.WithDefaultQuery,
		FieldNullable: true,
	})

	// Initialize a *gorm.DB instance
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Use the above `*gorm.DB` instance to initialize the generator,
	// which is required to generate structs from db when using `GenerateModel/GenerateModelAs`
	g.UseDB(db)

	// Generate default DAO interface for those specified structs
	g.ApplyBasic(models.User{}, models.Task{})

	// Execute the generator
	g.Execute()
}
