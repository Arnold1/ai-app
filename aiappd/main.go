package main

import (
    "log"
    "net/http"
    "html/template"
    "context"
    "time"
    "os/signal"
    "os"
    "syscall"

    "github.com/Arnold1/ai-app/build"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	build.Init()
}

func startWeb(webport string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)

	srv := &http.Server{
		Addr:    webport,
		Handler: mux,
	}

	go func() {
		// returns ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("ListenAndServe(): %s\n", err)
		}
	}()

	log.Printf("DebugUI on %s", webport[1:])

	// returning reference so caller can call Shutdown()
	return srv
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("view/index.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	type ServerStatus struct {
		BuildTimestamp      string
		BuildVersion        string
		StartedAt           string
		Uptime              string
		HTTPProfileEndPoint string
	}

	var toDisplay ServerStatus
	info := build.Info()
	toDisplay.BuildTimestamp = info.Time
	toDisplay.BuildVersion = info.SHA
	toDisplay.StartedAt = info.StartedAt
	toDisplay.Uptime = info.Uptime
	t.Execute(w, toDisplay)
}

// waitForSignals blocks until we get a signal telling the app to stop
func waitForSignals() {
	var signalHandler = make(chan os.Signal, 1)
	signal.Notify(signalHandler, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Running PID of", os.Getpid())

	s := <-signalHandler
	log.Println("Signal received", s)
}

func main() {
    debugUISvr := startWeb(":8080")

    waitForSignals()

    healthCTX, healthCancel := context.WithTimeout(context.Background(), time.Second*2)
	defer healthCancel()

    if err := debugUISvr.Shutdown(healthCTX); err != nil {
		log.Println(err)
	}
}