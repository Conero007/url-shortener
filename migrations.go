package main

func getTableCreationQueries() map[string]string {
	queries := make(map[string]string)

	// create urls table
	queries["urls"] = `CREATE TABLE IF NOT EXISTS urls (
		id INT PRIMARY KEY AUTO_INCREMENT,
		original_url VARCHAR(255) NOT NULL,
		short_key VARCHAR(20) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`

	return queries
}

func RunMigrations(app App) error {
	for _, query := range getTableCreationQueries() {
		if _, err := app.DB.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func RollbackMigrations(app App) error {
	for tableName := range getTableCreationQueries() {
		if _, err := app.DB.Exec("DELETE FROM " + tableName); err != nil {
			return err
		}
		if _, err := app.DB.Exec("ALTER TABLE " + tableName + " AUTO_INCREMENT = 1"); err != nil {
			return err
		}
	}
	return nil
}
