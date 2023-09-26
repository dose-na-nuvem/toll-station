package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "vehicles",

	Short: "Controle de veículos e tags",

	Long: `Permite gerenciar os veículos e as tags de customers:
Será possível adicionar carros com multiplas tags e relacionar os carros com os donos atuais.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
