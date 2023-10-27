package config

import (
	"time"

	"go.uber.org/zap"
)

type Cfg struct {
	Server   ServerSettings `mapstructure:"server"`
	Telemetry TelemetrySettings `mapstructure:"telemetry"`

	Logger *zap.Logger
}

type ServerSettings struct {
	HTTP HTTPServerSettings `mapstructure:"http"`
	// TLS  *TLSSettings       `mapstructure:"tls"`
}

type HTTPServerSettings struct {
	Endpoint          string        `mapstructure:"endpoint"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
}

type TelemetrySettings struct {
	Metrics MetricsSettings `mapstructure:"metrics"`
}

type MetricsSettings struct {
	Endpoint string `mapstructure:"endpoint"`
}

// type GRPCServerSettings struct {
// 	Endpoint string `mapstructure:"endpoint"`
// }

// type TLSSettings struct {
// 	CertFile    string `mapstructure:"cert_file"`
// 	CertKeyFile string `mapstructure:"cert_key_file"`
// 	Insecure    bool   `mapstructure:"insecure"`
// }

func New() *Cfg {
	cfg := &Cfg{
		Logger: zap.Must(zap.NewDevelopment()),
		Server: ServerSettings{
			//TLS: &TLSSettings{},
		},
		Telemetry: TelemetrySettings{
			Metrics: MetricsSettings{
				Endpoint: ":8888",
			},
		},
	}

	return cfg
}
