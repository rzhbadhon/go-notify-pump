package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gonotifysys.com/internal/handler"
	"gonotifysys.com/internal/worker"
)

func main() {

	wp := worker.NewWorkerPool(5, 10)

	wp.Run()

	fmt.Println("Worker Pool Started with 5 workers")

	notifHandler := handler.NewNotificationHandler(wp)

	r := gin.Default()

	r.POST("/send", notifHandler.SendNotification)

	srv := &http.Server{
		Addr: ":7090",
		Handler: r,
	}

	go func ()  {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed{
			log.Fatalf("listen: %s\n", err)
		}	
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil{
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Finishing pending jobs")
	wp.ShutDown()

	log.Println("Server exiting")

}
