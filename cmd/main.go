//go:build windows
// +build windows

package main

import (
	"encoding/csv"
	"fmt"
	"go-bms-dev/pkg/database"
	"log"
	"os"
)

func main() {

	connString := ""
	dbManager := database.NewDBManager(connString)
	defer dbManager.Close()

	query := "SELECT Sequence, ValueType, SampleValue FROM BAULNE__COGIR.dbo.tblTrendlog_2010100_0000000052"
	rows, err := dbManager.ExecuteQuery(query)
	if err != nil {
		log.Fatalf("Erreur lors de l'exécution de la requête : %v", err)
	}
	defer rows.Close()

	file, err := os.Create("resultats.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Créer un writer CSV à partir du fichier.
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// En-tête du fichier CSV.
	header := []string{"Sequence", "ValueType", "SampleValue"}
	if err := writer.Write(header); err != nil {
		log.Fatal(err)
	}

	// Parcourir les données et les écrire dans le fichier CSV.
	for rows.Next() {
		var Sequence int32
		var ValueType int16
		var SampleValue float32
		if err := rows.Scan(&Sequence, &ValueType, &SampleValue); err != nil {
			log.Fatal(err)
		}
		record := []string{fmt.Sprintf("%d", Sequence), fmt.Sprintf("%d", ValueType), fmt.Sprintf("%.2f", SampleValue)}
		if err := writer.Write(record); err != nil {
			log.Fatal(err)
		}
	}

	// Vérifier s'il y a des erreurs lors de l'itération sur les lignes.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}
