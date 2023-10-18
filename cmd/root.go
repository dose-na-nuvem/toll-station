package cmd

import (
	"os"
	"time"

	"github.com/dose-na-nuvem/toll-station/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	configFile     string
	cfg            = config.New()
	defaultTimeout = 5 * time.Second
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "pedagio",

	Short: "Central de Pedágio",

	Long: `Simula um pedágio aguardando carros, realizando cobranças e controlando cancelas`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(startCmd)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "config.yaml",
		"Define o arquivo de configuração a utilizar.")

	startCmd.Flags().StringVar(&cfg.Server.HTTP.Endpoint, "server.http.endpoint", "localhost:56433",
		"Endereço HTTP onde o serviço vai servir requisições.")

	startCmd.Flags().DurationVar(&cfg.Server.HTTP.ReadHeaderTimeout, "server.http.readheadertimeout",
		defaultTimeout,
		"Tempo máximo de leitura dos headers de uma requisição (Duration)")

	if err := viper.BindPFlags(startCmd.Flags()); err != nil {
		cfg.Logger.Error("falha ao ligar as flags", zap.Error(err))
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetConfigFile(configFile)

	// Tenta ler o arquivo de configuração, ignorando erros caso o mesmo não seja encontrado
	// Retorna um erro se não conseguirmos analisar o arquivo de configuração encontrado.
	if err := viper.ReadInConfig(); err != nil {
		// Não há problemas se não existir um arquivo de configuração.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			cfg.Logger.Error("arquivo não encontrado",
				zap.String("arquivo", configFile),
				zap.Error(err),
			)

			return
		}

		cfg.Logger.Error("falha na leitura do arquivo de configuração", zap.Error(err))
	} else {
		cfg.Logger.Info("arquivo de configuração lido", zap.String("config", configFile))
	}

	// convert Viper's internal state into our configuration object
	if err := viper.Unmarshal(&cfg); err != nil {
		cfg.Logger.Error("falhou ao converter o arquivo de configuração", zap.Error(err))

		return
	}
}
