package server

//TODO: review the http handling (return code, ...)
import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/mch1307/gomotics/config"
	"github.com/mch1307/gomotics/log"
	"github.com/mch1307/gomotics/nhc"
	"github.com/mch1307/gomotics/ws"

	"github.com/gorilla/mux"
)

// HealthMsg static alive json for health endpoint
const HealthMsg = `{"alive":true}`

//var Clients *ClientPool

// Server holds the gomotics app definition
type Server struct {
	Router     *mux.Router
	ListenPort string
	LogLevel   string
	LogFile    string
}

func init() {
	//log.Init()
	log.Info("Starting gomotics")
}

// Initialize initialize the server
// also calls the internal in mem db
func (s *Server) Initialize() {
	s.ListenPort = ":" + strconv.Itoa(config.Conf.ServerConfig.ListenPort)
	s.LogLevel = config.Conf.ServerConfig.LogLevel
	s.LogFile = config.Conf.ServerConfig.LogPath
	s.Router = mux.NewRouter().StrictSlash(true)
	s.intializeRoutes()
	// Initialize NHC in memory db
	nhc.Init(&config.Conf.NhcConfig)
}

func (s *Server) intializeRoutes() {
	s.Router.HandleFunc("/health", Health).Methods("GET")
	s.Router.HandleFunc("/api/v1/nhc/", nhc.GetNhcItems).Methods("GET")
	s.Router.HandleFunc("/api/v1/nhc/{id}", nhc.GetNhcItem).Methods("GET")
	s.Router.HandleFunc("/api/v1/nhc/action", nhc.NhcCmd).Methods("PUT")
	s.Router.HandleFunc("/events", ws.ServeWebSocket).Methods("GET")

	s.Router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		// p will contain regular expression is compatible with regular expression in Perl, Python, and other languages.
		// for instance the regular expression for path '/articles/{id}' will be '^/articles/(?P<v0>[^/]+)$'
		p, err := route.GetPathRegexp()
		if err != nil {
			return err
		}
		m, err := route.GetMethods()
		if err != nil {
			return err
		}
		fmt.Println(strings.Join(m, ","), t, p)
		log.Info(strings.Join(m, ","), t, p)
		return nil
	})
}

// Run starts the server process
func (s *Server) Run() {
	go nhc.Listener()
	log.Fatal(http.ListenAndServe(s.ListenPort, s.Router))
}

// Health endpoint for health monitoring
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(HealthMsg))
}
