package database

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

// DBManager est une structure pour gérer les connexions à la base de données.
type DBManager struct {
	db *sql.DB
}

// NewDBManager crée une nouvelle instance de DBManager.
func NewDBManager(connString string) *DBManager {
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de la connexion : %v", err)
	}

	// Vérifiez la connexion
	err = db.Ping()
	if err != nil {
		log.Fatalf("Erreur lors de la vérification de la connexion : %v", err)
	}

	fmt.Println("Connexion à la base de données réussie !")
	return &DBManager{db: db}
}

// Close ferme la connexion à la base de données.
func (manager *DBManager) Close() {
	if err := manager.db.Close(); err != nil {
		log.Fatalf("Erreur lors de la fermeture de la connexion à la base de données : %v", err)
	}
}

// ExecuteQuery exécute une requête SQL avec des arguments optionnels et retourne un *sql.Rows.
func (manager *DBManager) ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := manager.db.Query(query, args...)
	if err != nil {
		log.Printf("Erreur lors de l'exécution de la requête : %v", err)
		return nil, err
	}
	return rows, nil
}
