package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o pedágio",
	Long:  `Libera as faixas sem-parar e faz cobranças`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("🚗💨💰💰💰 pedágio funcionando...")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
