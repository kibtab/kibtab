// Command kibtab starts the Kibtab instance.
//
// Kibtab turns a spreadsheet into a client for a relational database.
// Release v0.1.0 starts the instance and it answers a health check.
//
// The command reads PORT for the listen address and DATABASE_URL for the
// database. Release v0.3.0 connects to the database. This release reads
// the address and it reports when the address is empty.
package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

// version is the Kibtab version. The plan in plan.md holds the releases.
// The build stamps a release binary; a plain build answers "devel".
var version = "devel"

// healthResponse is the body of GET /healthz.
type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// healthzHandler answers the health check with the status and the version.
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	body := healthResponse{Status: "ok", Version: version}
	if err := encoder.Encode(body); err != nil {
		log.Printf("healthz: write response: %v", err)
	}
}

// listenAddress reads the listen address from PORT. It defaults to 8080.
func listenAddress(getenv func(string) string) (string, error) {
	raw := getenv("PORT")
	if raw == "" {
		return "127.0.0.1:8080", nil
	}
	port, err := strconv.Atoi(raw)
	if err != nil {
		return "", err
	}
	if port < 1 || port > 65535 {
		return "", errInvalidPort{port: port}
	}
	return "127.0.0.1:" + strconv.Itoa(port), nil
}

// errInvalidPort reports a PORT value outside the TCP port range.
type errInvalidPort struct {
	port int
}

func (e errInvalidPort) Error() string {
	return "PORT must be a TCP port number, got " + strconv.Itoa(e.port)
}

// checkDatabaseURL reads the database address. Release v0.3.0 uses it.
func checkDatabaseURL(getenv func(string) string) error {
	if getenv("DATABASE_URL") == "" {
		return errMissingDatabaseURL{}
	}
	return nil
}

// errMissingDatabaseURL reports an empty DATABASE_URL.
type errMissingDatabaseURL struct{}

func (e errMissingDatabaseURL) Error() string {
	return "DATABASE_URL is not set"
}

// newMux builds the HTTP routes of the instance.
func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzHandler)
	return mux
}

// run wires the adapters and it serves the instance. It returns the error
// from serve. The command logs the error and it exits.
func run(getenv func(string) string, stderr io.Writer) error {
	log.SetOutput(stderr)
	address, err := listenAddress(getenv)
	if err != nil {
		return err
	}
	if err := checkDatabaseURL(getenv); err != nil {
		log.Printf("config: %v", err)
	}
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	log.Printf("kibtab %s listens on http://%s", version, listener.Addr().String())
	return serve(listener, newMux())
}

// serve holds the serve call of the standard library. The test replaces it.
var serve = http.Serve

// exit ends the process. The test replaces it.
var exit = os.Exit

func main() {
	if err := run(os.Getenv, os.Stderr); err != nil {
		log.Printf("serve: %v", err)
		exit(1)
	}
}
