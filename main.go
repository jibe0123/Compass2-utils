//go:build windows
// +build windows

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	mdbDir := "./Trendlog_2013000_0000000020/"
	rootExportDir := "./data-exported"

	if err := os.MkdirAll(rootExportDir, os.ModePerm); err != nil {
		fmt.Println("Erreur lors de la création du dossier racine d'exportation:", err)
		return
	}

	files, err := ioutil.ReadDir(mdbDir)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du dossier:", err)
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".mdb" {
			dbPath := filepath.Join(mdbDir, file.Name())
			fmt.Println("Traitement du fichier:", dbPath)

			specificExportDir := filepath.Join(rootExportDir, strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
			if err := os.MkdirAll(specificExportDir, os.ModePerm); err != nil {
				fmt.Println("Erreur lors de la création du dossier d'exportation spécifique:", err)
				continue
			}

			processMDB(dbPath, specificExportDir)
		}
	}
}

func processMDB(dbPath, outputDir string) {
	tables, err := getTables(dbPath)
	if err != nil {
		fmt.Println("Erreur lors de la récupération des tables:", err)
		return
	}

	for _, table := range tables {
		fmt.Println("Exportation de la table:", table)
		if err := exportTableToCsv(dbPath, table, outputDir); err != nil {
			fmt.Println("Erreur lors de l'exportation de la table:", table, err)
			continue
		}
	}
}

func getTables(dbPath string) ([]string, error) {
	cmd := exec.Command("mdb-tables", "-1", dbPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	tables := strings.Split(strings.TrimSpace(out.String()), "\n")
	return tables, nil
}

func exportTableToCsv(dbPath, tableName, outputDir string) error {
	outputFilePath := filepath.Join(outputDir, fmt.Sprintf("%s.csv", tableName))
	outFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	cmd := exec.Command("mdb-export", dbPath, tableName)
	cmd.Stdout = outFile
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println("Données exportées avec succès dans:", outputFilePath)
	return nil
}
