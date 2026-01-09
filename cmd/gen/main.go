package gen

import (
	"core-consumer/config"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	dsn := cfg.DBString
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/app/gen/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(db)

	g.GenerateAllTable()

	g.Execute()
}
