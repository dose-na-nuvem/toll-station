package server

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	// "github.com/dose-na-nuvem/customers/pkg/model"
	// "github.com/dose-na-nuvem/customers/pkg/telemetry"
	// "go.opentelemetry.io/otel/attribute"
	// "go.opentelemetry.io/otel/trace"
	"github.com/dose-na-nuvem/toll-station/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

var _ http.Handler = (*TollStationHandler)(nil)

// type CustomerStore interface {
// 	CreateCustomer(string) (*model.Customer, error)
// 	ListCustomers() ([]model.Customer, error)
// }

type TollStationHandler struct {
	logger    *zap.Logger
	telemetry *telemetry.Telemetry
	// store  CustomerStore
}

func NewTollStationHandler(logger *zap.Logger, tm *telemetry.Telemetry /*, store CustomerStore*/) *TollStationHandler {
	return &TollStationHandler{
		logger:    logger,
		telemetry: tm,
		//store:  store,
	}
}

func (h *TollStationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()
	// ctx := r.Context()
	ctx := context.Background()

	// Get the tag from the POST request.
	tag := r.FormValue("tag")

	gateOpenState := shouldOpenGate(tag)

	attrs := otelmetric.WithAttributes(
		attribute.Key("open").Bool(gateOpenState),
	)
	h.telemetry.TrafficCounter.Add(ctx, 1, attrs)

	h.logger.Info("Estado do portão", zap.Bool("aberto", gateOpenState))

	b, err := json.Marshal(gateOpenState)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	duration := time.Since(startTime).Abs().Milliseconds()
	h.telemetry.GateHistogram.Record(r.Context(), duration, attrs)

	// switch r.Method {
	// case http.MethodPost:
	// 	h.createCustomer(w, r)
	// case http.MethodGet:
	// 	h.listCustomers(w, r)
	// default:
	// 	w.WriteHeader(http.StatusNotImplemented)
	// }
}

func randomDelayMs() {
	r := rand.Intn(300) + 10
	time.Sleep(time.Duration(r) * time.Millisecond)
}

func shouldOpenGate(tag string) bool {
	// Todo: Consultar a tag/cliente e fazer pagamentos
	r := rand.Intn(100)
	if r < 50 {
		randomDelayMs()
		return true
	}

	// Otherwise, keep the gate closed.
	return false
}

// func (h *CustomerHandler) createCustomer(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseForm()
// 	if err != nil {
// 		h.logger.Warn("erro ao varrer dados do post")
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}
// 	name := r.PostForm.Get("name")
//
// 	_, span := telemetry.GetTracer().Start(r.Context(), "create-customer")
// 	defer span.End()
// 	_, err = h.store.CreateCustomer(name)
// 	if err != nil {
// 		h.logger.Warn("falha ao criar um customer", zap.Error(err))
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}
// }
//
// func (h *CustomerHandler) listCustomers(w http.ResponseWriter, r *http.Request) {
// 	_, span := telemetry.GetTracer().Start(r.Context(), "list-customers")
// 	defer span.End()
//
// 	c, err := h.store.ListCustomers()
// 	if err != nil {
// 		span.RecordError(err, trace.WithAttributes(
// 			attribute.String("error.message", "Falha ao consultar customers"),
// 		))
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}
//
// 	b, err := json.Marshal(c)
// 	if err != nil {
// 		span.RecordError(err, trace.WithAttributes(
// 			attribute.String("error.message", "Falha ao serializar customers"),
// 		))
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}
//
// 	w.Header().Set("Content-Type", "application/json")
// 	_, err = w.Write(b)
// 	if err != nil {
// 		span.RecordError(err, trace.WithAttributes(
// 			attribute.String("error.message", "Falha escrever a resposta da requisição"),
// 		))
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}
// }
