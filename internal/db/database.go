package db

// логика работы с базой данных
import (
	"context"
	"fmt"

	"mycli/internal/models"

	"github.com/jackc/pgx/v5"
)

func DbInsert(configs []models.Config, dbConfig string) error {

	ctx := context.Background()
	// postgres://USERNAME:PASSWORD@HOST:PORT/DBNAME
	conn, err := pgx.Connect(ctx, dbConfig)
	if err != nil {
		fmt.Println("база не поднялась")
		return err
	}
	defer conn.Close(ctx)

	query := `
        INSERT INTO configs (name, description, version, author, tags)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (name)
        DO UPDATE SET 
			description = EXCLUDED.description,
			version = EXCLUDED.version,
			author = EXCLUDED.author,
			tags = EXCLUDED.tags
    `
	for _, config := range configs {
		_, err = conn.Exec(ctx, query,
			config.Name,
			config.Description,
			config.Version,
			config.Metadata.Author,
			config.Metadata.Tags,
		)
		if err != nil {
			return fmt.Errorf("Error insert, %s", err)
		}
	}
	return nil
}
