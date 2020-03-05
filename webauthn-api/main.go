package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func main() {
	static := os.Getenv("STATIC_DIR")
	if static == "" {
		static = "static"
	}

	origin := os.Getenv("ORIGIN")
	if origin == "" {
		origin = "http://localhost:3000"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	sesskey := os.Getenv("SESSION_KEY")
	if sesskey == "" {
		sesskey = "SESSION_KEY"
	}

	var sesssecure bool
	if os.Getenv("SESSION_SECURE") == "true" {
		sesssecure = true
	}

	userStore := &UserInMemoryStore{
		store: NewInMemoryStore(),
	}

	webauthnStore := &WebauthnStore{
		store:  sessions.NewCookieStore([]byte(sesskey)),
		secure: sesssecure,
	}

	uri, err := url.ParseRequestURI(origin)
	if err != nil {
		log.Fatalf("parse origin: %s", err)
	}

	authn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "WebAuthn Experiments",
		RPID:          uri.Hostname(),
		RPOrigin:      uri.String(),
	})
	if err != nil {
		log.Fatalf("new webauthn: %s", err)
	}

	startRegisterHandler := StartRegister(authn, webauthnStore, userStore)
	completeRegisterHandler := CompleteRegister(authn, webauthnStore, userStore)

	startLoginHandler := StartLogin(authn, webauthnStore, userStore)
	completeLoginHandler := CompleteLogin(authn, webauthnStore, userStore)

	listHandler := List()

	router := mux.NewRouter()

	midleware := Auth(userStore)

	router.Handle("/api/register/start", startRegisterHandler).Methods(http.MethodOptions, http.MethodPost)
	router.Handle("/api/register/complete", completeRegisterHandler).Methods(http.MethodOptions, http.MethodPost)

	router.Handle("/api/login/start", startLoginHandler).Methods(http.MethodOptions, http.MethodPost)
	router.Handle("/api/login/complete", completeLoginHandler).Methods(http.MethodOptions, http.MethodPost)

	router.Handle("/api/list", midleware(listHandler)).Methods(http.MethodOptions, http.MethodGet)

	router.PathPrefix("/public").Handler(http.StripPrefix("/public", http.FileServer(http.Dir(static))))
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(static, "index.html"))
	})

	server := http.Server{
		Addr:         host + ":" + port,
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	// We don't support the graceful shutdown of the server. For a "real-world" program,
	// we have to close the server with a cancelable context. We also have to provide
	// a TimeoutHandler to exit as soon as we hit the timeout.
	// https://ieftimov.com/post/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/
	log.Printf("Server listen on port %s", host+":"+port)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Server stopped: %v", err)
	}
}
