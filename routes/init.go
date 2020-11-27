package routes

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/powerqueue/fitque-users-api/services"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
)

//LoginServer - struct for Login server
type LoginServer struct {
	loginService services.LoginService
}

func (ls *LoginServer) defineRoutes(port string, metricsPort string) {
	apiPrefix := "/fitqueue-login-api/v1"

	ar := mux.NewRouter().PathPrefix(apiPrefix).Subrouter()

	//configure routes
	// ls.addLoginRoutes(ar)
	ar.HandleFunc("/retrieve-login", ls.RetrieveLoginHandler).Methods("POST")
	ar.HandleFunc("/create-login", ls.CreateLogin).Methods("POST")
	ar.HandleFunc("/term-login", ls.TermLogin).Methods("POST")

	r := mux.NewRouter()

	//server swagger
	swaggerUIBox := packr.New("SwaggerUiBox", "../resources/swagger-ui")
	swaggerSpecBox := packr.New("SwaggerSpecBox", "../resources/api/server")
	r.PathPrefix("/swagger-ui").Handler(http.StripPrefix("/swagger-ui", http.FileServer(swaggerUIBox)))
	r.PathPrefix("/swagger").Handler(http.StripPrefix("/swagger", http.FileServer(swaggerSpecBox)))

	// serve api with metrics
	r.PathPrefix(apiPrefix)
	http.Handle("/", r)
	// rest.InitServer().Add(r, port, rest.RouteInfo{SubRoutes: ar, IsSecure: true, ApiPrefix: apiPrefix}).
	// 	Metrics(metricsPort).Start()
	srv := &http.Server{
		Handler: ar,
		Addr:    "127.0.0.1:8095",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("About to start server at port :8095")

	log.Fatal(srv.ListenAndServe())
}

//NewLoginServer - initialize server
func NewLoginServer(loginServer services.LoginService) *LoginServer {
	return &LoginServer{
		loginService: loginServer,
	}
}

// InitServer -- initialize web server
func InitServer(port string, metricsPort string, loginService services.LoginService) {
	lServer := NewLoginServer(loginService)
	lServer.defineRoutes(port, metricsPort)
}
