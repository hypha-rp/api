package db

import (
	"embed"
	"fmt"
	"hypha/api/internal/config"
	"hypha/api/internal/db/tables"
	"hypha/api/internal/utils/logging"

	"github.com/go-orm/gorm"
	_ "github.com/go-orm/gorm/dialects/postgres"
)

var log = logging.Logger

//go:embed views/test_results_view.sql
var test_results_view embed.FS

var (
	// DBConn is the global database connection.
	DBConn *gorm.DB
)

// DBConnWrapper wraps a gorm.DB connection.
type DBConnWrapper struct {
	DB *gorm.DB
}

// Connect establishes a connection to the database using the provided configuration.
// It returns the database connection and any error encountered.
//
// Parameters:
//   - cfg: A pointer to the configuration object containing database connection details.
//
// Returns:
//   - *gorm.DB: A pointer to the established database connection.
//   - error: An error object if the connection fails, otherwise nil.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	log.Info().
		Str("host", cfg.Database.Host).
		Int("port", cfg.Database.Port).
		Str("dbname", cfg.Database.Dbname).
		Msg("Attempting to connect to the database")

	db, err := gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Dbname,
		cfg.Database.Password,
		cfg.Database.Sslmode,
	))
	if err != nil {
		log.Error().
			Err(err).
			Str("host", cfg.Database.Host).
			Int("port", cfg.Database.Port).
			Str("dbname", cfg.Database.Dbname).
			Msg("Failed to connect to the database")
		return nil, err
	}

	log.Info().
		Str("host", cfg.Database.Host).
		Int("port", cfg.Database.Port).
		Str("dbname", cfg.Database.Dbname).
		Msg("Successfully connected to the database")

	if cfg.Database.Debug {
		DBConn = db.Debug()
	} else {
		DBConn = db
	}
	return DBConn, nil
}

var tables_slice = []interface{}{
	&tables.Product{},
	&tables.Relationship{},
	&tables.Result{},
	&tables.TestSuite{},
	&tables.TestCase{},
	&tables.Property{},
	&tables.ResultsRule{},
}

// AutoMigrate performs database migration for all the tables defined in tables_slice.
// It returns any error encountered during the migration process.
//
// Parameters:
//   - db: A pointer to the gorm.DB connection.
//
// Returns:
//   - error: An error object if the migration fails, otherwise nil.
func AutoMigrate(db *gorm.DB) error {
	log.Info().Msg("Starting database migration")
	for _, table := range tables_slice {
		if err := db.AutoMigrate(table).Error; err != nil {
			log.Error().Err(err).Msg("Database migration failed")
			return err
		}
	}
	createViews(db)
	log.Info().Msg("Database migration completed successfully")
	return nil
}

// CreateViews reads SQL files from the assets package and executes them to create views in the database.
// It returns any error encountered during the execution.
//
// Parameters:
//   - db: A pointer to the gorm.DB connection.
//
// Returns:
//   - error: An error object if the execution fails, otherwise nil.
func createViews(db *gorm.DB) error {

	log.Info().Msg("Creating test_results_view")

	// Read the embedded SQL file
	content, err := test_results_view.ReadFile("views/test_results_view.sql")
	if err != nil {
		log.Error().Err(err).Msg("Failed to read embedded SQL file: test_results_view.sql")
		return err
	}

	// Execute the SQL content
	if err := db.Exec(string(content)).Error; err != nil {
		log.Error().Err(err).Msg("Failed to execute SQL file: test_results_view.sql")
		return err
	}

	log.Info().Msg("Successfully created test_results_view")
	return nil
}

// Connection returns the wrapped gorm.DB connection.
//
// Returns:
//   - *gorm.DB: A pointer to the wrapped gorm.DB connection.
func (wrapper *DBConnWrapper) Connection() *gorm.DB {
	return wrapper.DB
}

// Create inserts a new record into the database.
// It returns any error encountered during the insertion.
//
// Parameters:
//   - record: The record to be inserted into the database.
//
// Returns:
//   - error: An error object if the insertion fails, otherwise nil.
func (wrapper *DBConnWrapper) Create(record interface{}) error {
	return wrapper.DB.Create(record).Error
}

// First retrieves the first record that matches the given conditions.
// It returns any error encountered during the retrieval.
//
// Parameters:
//   - result: A pointer to the object where the result will be stored.
//   - conditions: Optional conditions to filter the query.
//
// Returns:
//   - error: An error object if the retrieval fails, otherwise nil.
func (wrapper *DBConnWrapper) First(result interface{}, conditions ...interface{}) error {
	return wrapper.DB.First(result, conditions...).Error
}
