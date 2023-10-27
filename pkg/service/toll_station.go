package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dose-na-nuvem/toll-station/config"
	// "github.com/dose-na-nuvem/toll-station/pkg/model"
	"github.com/dose-na-nuvem/toll-station/pkg/server"
	"github.com/dose-na-nuvem/toll-station/pkg/telemetry"

	// "github.com/dose-na-nuvem/toll-station/pkg/telemetry"
	// "go.opentelemetry.io/otel"
	"go.uber.org/zap"
	// "gorm.io/driver/sqlite"
	// "gorm.io/gorm"
)

const MaxServers = 2

type TollStation struct {
	cfg *config.Cfg
	srv *server.HTTP
	// grpc              *server.GRPC
	telemetry         *telemetry.Telemetry
	asyncErrorChannel chan error
	signalsChannel    chan os.Signal
}

func New(cfg *config.Cfg, tm *telemetry.Telemetry) *TollStation {
	return &TollStation{
		cfg:               cfg,
		asyncErrorChannel: make(chan error, MaxServers), // buffered
		signalsChannel:    make(chan os.Signal),
		telemetry:         tm,
	}
}

// func (c *Customer) bootstrap(ctx context.Context) (server.CustomerStore, error) {
// 	// este rastreador é somente para o processo de bootstrapping
//
// 	tr := otel.GetTracerProvider().Tracer("bootstrap")
//
// 	ctx, rootSpan := tr.Start(ctx, "bootstrap")
//
// 	_, span := tr.Start(ctx, "db/open")
// 	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
// 	if err != nil {
// 		return nil, fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
// 	}
// 	span.End()
//
// 	// Migrate the schema
// 	_, span = tr.Start(ctx, "db/migrate")
// 	if err := db.AutoMigrate(&model.Customer{}); err != nil {
// 		return nil, fmt.Errorf("falha ao migrar o esquema do banco de dados: %w", err)
// 	}
// 	span.End()
//
// 	rootSpan.End()
//
// 	return model.NewStore(db), nil
// }

func (t *TollStation) Start(ctx context.Context) error {

	t.cfg.Logger.Info("🚗💨💰💰💰 pedágio funcionando...")
	var err error

	// tp, err := telemetry.NewTracerProvider()
	// if err != nil {
	// 	return fmt.Errorf("falha ao iniciar os rastreadores: %w", err)
	// }
	// otel.SetTracerProvider(tp)
	//
	// store, err := c.bootstrap(ctx)
	// if err != nil {
	// 	return err
	// }

	ch := server.NewTollStationHandler(t.cfg.Logger, t.telemetry /*store*/)

	// t.grpc, err = server.NewGRPC(c.cfg, store)
	// if err != nil {
	// 	return fmt.Errorf("falha ao iniciar o servidor GRPC: %w", err)
	// }
	// c.grpc.Start(ctx, c.asyncErrorChannel)

	t.srv, err = server.NewHTTP(t.cfg, ch)
	if err != nil {
		return fmt.Errorf("falha ao iniciar o servidor HTTP: %w", err)
	}
	t.srv.Start(ctx, t.asyncErrorChannel)

	signal.Notify(t.signalsChannel, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(t.signalsChannel)

LOOP:
	for {
		select {
		case err := <-t.asyncErrorChannel:
			t.cfg.Logger.Error("falha ao iniciar o servidor: %w", zap.Error(err))
			break LOOP
		case signal := <-t.signalsChannel:
			t.cfg.Logger.Debug("signal received", zap.Any("signal", signal.String()))
			err := t.Shutdown(ctx)
			if err != nil {
				t.cfg.Logger.Error("falha ao finalizar servidor: %w", zap.Error(err))
			}
			return nil
		}
	}

	return nil
}

func (t *TollStation) Shutdown(ctx context.Context) error {
	if err := t.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("erro ao finalizar o serviço: %w", err)
	}

	t.telemetry.Shutdown(ctx)

	// if err := t.grpc.Shutdown(ctx); err != nil {
	// 	return fmt.Errorf("erro ao finalizar o serviço: %w", err)
	// }

	return nil
}
