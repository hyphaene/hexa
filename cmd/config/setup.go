package config

import (
	_ "embed"
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	ConfigCmd.AddCommand(SetupCmd)
}

//go:embed template.yml
var template []byte

var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup and initialize hexa CLI",
	Long:  `Setup and initialize hexa CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Setup logic here
		userHomeDir, _ := os.UserHomeDir()
		filePath := filepath.Join(userHomeDir, ".hexa.yml")
		_, err := os.Stat(filePath)
		if err == nil {
			println("User config file already exists at", filePath)
			return
		}
		if errors.Is(err, os.ErrNotExist) {
			println("Creating user config file at", filePath)
			if writeErr := os.WriteFile(filePath, template, 0644); writeErr != nil {
				println("Failed to create user config file:", writeErr.Error())
				return
			}
			return
		}
		println("Error checking user config file:", err.Error())

		// create the file from template
		// write to ./hexa.local.yml
		// check if .gitignore exists
		// - if yes, append .hexa.local.yml if not already present
		// - if no, create .gitignore with .hexa.local.yml inside

		// en mode interactif :
		// un truc dans le genre

		// voulez vous créer un fichier de config global ? (y/n)
		// check si existe
		// - si oui, dire qu'il existe deja
		// => écraser ou skip
		// - si non, le créer
		// voulez vous créer un fichier de config local ? (y/n)
		// ici on va : chercher le pwd
		// - creer le .hexa.local.yml ?
		// voir si y'a un fichier gitignore
		// - si oui, ajouter .hexa.local.yml dedans
		// - si non, creer un fichier gitignore avec .hexa.local.yml dedans
		// - ajouter les infos de base dans le .hexa.local.yml ( user.me: xxx)
		// - dire a l'utilisateur de modifier le fichier pour ajouter son jira token et son jira domain
		// - lui dire de relancer la commande pour initialiser le reste

	},
}
