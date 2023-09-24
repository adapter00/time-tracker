package main

import (
	"fmt"
	"log"
    "os"

	"github.com/k0kubun/sqldef"
	"github.com/k0kubun/sqldef/database"
	"github.com/k0kubun/sqldef/database/postgres"
	"github.com/k0kubun/sqldef/parser"
	"github.com/k0kubun/sqldef/schema"
)

func main() {
    migrate("../../schema.sql")
}
func migrate( schemaFile string) error {
	sqlParser := database.NewParser(parser.ParserModePostgres)
	desiredDDLs, err := sqldef.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("Failed to read %s: %w", schemaFile, err)
	}
    db,err:=postgres.NewDatabase(database.Config{
        DbName:"slack",
        User:"user",
        Password:"password",
        Host: "0.0.0.0",
        Port:54321,
    })
        os.Setenv("PGSSLMODE", "disable")

    if err!=nil {
        log.Fatal(err)
    }
    options := &sqldef.Options{DesiredDDLs: desiredDDLs,Export:true}
	sqldef.Run(schema.GeneratorModePostgres, db, sqlParser, options)
    return nil
}
