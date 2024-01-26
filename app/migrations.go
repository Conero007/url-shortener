package app

func getTableCreationQueries() map[string]string {
	return map[string]string{
		"urls": `CREATE TABLE IF NOT EXISTS urls (
			id INT PRIMARY KEY AUTO_INCREMENT,
			original_url VARCHAR(255) NOT NULL,
			short_key VARCHAR(20) NOT NULL UNIQUE,
			expire_time TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);`,
	}
}

func RunMigrations(dbName string) error {
	if err := createAndUseDatabase(dbName); err != nil {
		return err
	}

	for _, query := range getTableCreationQueries() {
		if _, err := App.DB.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func RollbackMigrations() error {
	for tableName := range getTableCreationQueries() {
		if _, err := App.DB.Exec("DROP TABLE " + tableName); err != nil {
			return err
		}
	}
	return nil
}

func createAndUseDatabase(dbName string) error {
	if _, err := App.DB.Exec("CREATE DATABASE IF NOT EXISTS " + dbName); err != nil {
		return err
	}
	if _, err := App.DB.Exec("USE " + dbName); err != nil {
		return err
	}
	return nil
}
