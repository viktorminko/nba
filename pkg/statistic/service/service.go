package service

import (
	"context"
	"github.com/viktorminko/nba/pkg/statistic/frontend"
	"github.com/viktorminko/nba/pkg/statistic/stats"
	"github.com/viktorminko/nba/pkg/statistic/subscriber"
	"log"
	"net/http"
	"strconv"
	"time"
)

func startServer(ctx context.Context, port int, st *stats.Stats, frontend frontend.Displayer) {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := frontend.Display(w, st)
		if err != nil {
			log.Println("error in frontend", err)
			http.Error(w, "error displaying data", http.StatusInternalServerError)
		}
	}))

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Print("stopping server")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Print("Server stopped")
}

func startErrorHandler(errCh <-chan error) {
	go func() {
		for {
			select {
			case err, ok := <-errCh:
				if !ok {
					return
				}
				log.Println(err)
			}
		}
	}()
}

func Start(ctx context.Context, eventsSubscriber subscriber.Subscriber, port int, frontend frontend.Displayer) error {
	log.Println("Starting service")

	st := stats.New()

	startErrorHandler(st.StartUpdater(ctx, eventsSubscriber.Subscribe()))

	startServer(ctx, port, st, frontend)

	return nil
}
