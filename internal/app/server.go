package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

var (
	Repo *repository.Repository
	Hand *httphandler.Handler
	Cnfg *models.Config
)

func init() {
	var err error
	Cnfg, err = initconfig.InitConfig("config.json")
	if err != nil {
		log.Error("Failed to initialize the config:", err)
		return
	}
	Repo = repository.NewRepository(Cnfg.ConnectionString)
	Hand = httpdhandler.NewHandler(Repo, Cnfg)
	go Hand.StartScheduler(context.TODO())
}

func (a *Application) StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/currency/save/{date}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.SaveCurrencyHandler(w, r.WithContext(ctx), ctx)

	})

	r.HandlerFunc("/currency/{date}/{code}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.GetCurrencyHandler(w, r.WithContext(ctx), ctx)
	})

	r.HandleFunc("/currency/{date}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.GetCurrencyHandler(w, r.WithContext(ctx), ctx)

	})

	r.HandlerFunc("/delete/{date}/{code}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.DeleteCurrencyHandler(w, r.WithContext(ctx), ctx)
	})

	r.HandlerFunc("/delete/{date}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.DeleteCurrencyHandler(w, r.WithContext(ctx), ctx)
	})

	go func() {
		if err := http.ListenAndServe(":8081", r); err != nil {
			fmt.Println("Failed to start the metrics server:", err)
		}
	}()
	server := &http.Server{
		Addr:         ":" + Cnfg.ListenPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}
	quit := make(chan os.Signal, 1)
	go shutdown(quit)
	log.Println("Listening on port", Cnfg.ListenPort, "...")
	log.Fatal(server.ListenAndServe())
}

func shutdown(quit chan os.Signal) {
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	log.Error("caught signal", "signal", s.String())
	os.Exit(0)
}
