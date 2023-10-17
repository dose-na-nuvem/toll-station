package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/dose-na-nuvem/toll-station/config"
	"github.com/dose-na-nuvem/toll-station/pkg/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o pedágio",
	Long:  `Libera as faixas sem-parar e faz cobranças`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		ctx := context.Background()

		svc := service.New(cfg)

		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			<-sigint

			cfg.Logger.Info("finalizando o serviço")

			// TODO: colocar uma deadline para o shutdown
			if err := svc.Shutdown(ctx); err != nil {
				cfg.Logger.Error("erro ao finalizar o serviço: %w", zap.Error(err))
			}
			cfg.Logger.Info("serviço finalizado com sucesso")
		}()

		err := svc.Start(ctx)

		if err != nil {
			log.Println("Não foi possível inicializar o serviço! Abortando execução...")
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
