package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o microserviço",
	Long: `Permite gerenciar os veículos e as tags de customers:

Será possível adicionar carros com multiplas tags e relacionar os carros com clientes.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚗💨 vehicles running")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
