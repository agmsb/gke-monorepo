package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var version = os.Getenv("VERSION")
var backend string

func main() {

	port := ":8080"
	backend = *flag.String("backend-service", "http://gke-monorepo-backend:8080", "hostname of backend server")
	flag.Parse()

	log.Printf("Frontend version: %s\n", version)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", handleIndex)
	r.GET("/version", handleVersion)
	r.GET("/healthz", handleHealthz)

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")

}

func handleIndex(c *gin.Context) {
	// Call backend for info
	resp, err := http.Get(backend)
	if err != nil {
		c.String(http.StatusInternalServerError, "Request to backend server ()%s failed:\n%v", backend, err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to read body from backend request:\n%v", err)
		return
	}
	var p PodMetadata
	log.Println(p)
	err = json.Unmarshal([]byte(body), &p)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to parse JSON from backend request:\n%v", err)
		return
	}
	c.HTML(http.StatusOK, "index.html", p)
}

func handleVersion(c *gin.Context) {
	c.String(http.StatusOK, "%s", c.Value("version"))
}

func handleHealthz(c *gin.Context) {
	c.String(http.StatusOK, "", "")
}

// PodMetadata represents info about an InstanceMetadata in GCE
type PodMetadata struct {
	Name       string
	Namespace  string
	HostIP     string
	PodIP      string
	StartTime  string
	RawRequest string
}
