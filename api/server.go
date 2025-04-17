package rest

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sugyk/rest_vpn/api/handlers"
	"github.com/sugyk/rest_vpn/lib/configs"
	"github.com/sugyk/rest_vpn/lib/database"
)

type Server struct {
	service *handlers.Service
	router  *mux.Router
}

func NewServer() *Server {

	err := configs.LoadConfigs()
	if err != nil {
		log.Fatalln("Error parsing db config file", err)
	}

	// Parse the db config file
	db_cnf := configs.NewDbConfigs()
	outline_config := configs.NewOutlineConfigs()

	log.Println("Successfully parsed db config file")

	db, err := database.NewPostgresConnection(database.ConfingDB(db_cnf))
	if err != nil {
		log.Fatalln("Error connecting to database", err)
	}

	log.Println("Successfully connected to database")

	router := mux.NewRouter()
	handlers.Register(router, db, outline_config)

	return &Server{
		router: router,
	}
}

func (s *Server) Run(port string) error {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: cors.Default().Handler(s.router),
	}

	stopServer := make(chan os.Signal, 1)
	signal.Notify(stopServer, syscall.SIGINT, syscall.SIGTERM)

	servChannel := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Server is running on port", port)
		servChannel <- server.ListenAndServe()
	}()

	select {
	case err := <-servChannel:
		log.Fatalln("Error starting server", err)
	case <-stopServer:
		err := server.Shutdown(nil)
		if err != nil {
			return fmt.Errorf("Graceful shutdown did not complete: %e", err)
		}
		log.Println("Shuting down the server")
		wg.Wait()
	}

	return nil
}
