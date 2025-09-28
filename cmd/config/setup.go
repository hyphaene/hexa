package config

import (
	_ "embed"
	"os"

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

		_, userErr := os.Stat(os.ExpandEnv("$HOME/.hexa.yml"))
		isUserConfigFileExists := userErr == nil
		_, localErr := os.Stat(".hexa.local.yml")
		isLocalConfigFileExists := localErr == nil

		if isUserConfigFileExists {
			println("User config file already exists at $HOME/.hexa.yml")
		} else {
			println("Creating user config file at $HOME/.hexa.yml")

			os.WriteFile(os.ExpandEnv("$HOME/.hexa.yml"), []byte(template), 0644)
		}

		if isLocalConfigFileExists {
			println("Local config file already exists at ./hexa.local.yml")
		} else {
			println("Creating local config file at ./hexa.local.yml")
		}
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
