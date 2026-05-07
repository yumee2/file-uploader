package main

import (
	"context"
	"errors"
	server "file-uploader/internal/http"
	"file-uploader/internal/repository/sqlite"
	"file-uploader/internal/service"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"
)

func main() {
	setupAuth := flag.Bool("setup-auth", false, "set up authentication password")
	flag.Parse()

	dbConn, err := sqlite.NewDBConnection()
	if err != nil {
		fmt.Printf("Error during creating a database: %s", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	fileService := service.NewFileService(dbConn)
	if *setupAuth {
		password, err := readPassword()
		if err != nil {
			fmt.Printf("Error reading password: %s\n", err)
			os.Exit(1)
		}

		if err := fileService.CreatePassword(password); err != nil {
			fmt.Printf("Error creating password: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("Password created")
		return
	}
	fileHandler := server.NewFileHandler(fileService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /files", fileHandler.GetFiles)
	mux.HandleFunc("GET /files/{id}", fileHandler.DownloadFile)
	mux.HandleFunc("DELETE /files/{id}", fileHandler.DeleteFile)
	mux.HandleFunc("POST /files", fileHandler.AddFile)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	serverErr := make(chan error, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()
	log.Println("server started on :8000")

	select {
	case err := <-serverErr:
		log.Printf("server stopped unexpectedly: %v", err)
		return
	case <-quit:
		log.Println("server shutting down")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		return
	}
	log.Println("http server stopped")
}

func readPassword() (string, error) {
	fmt.Print("Enter password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}

	fmt.Print("Confirm password: ")
	confirmBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}

	password := string(passwordBytes)
	confirm := string(confirmBytes)

	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}
	if password != confirm {
		return "", fmt.Errorf("passwords do not match")
	}

	return password, nil
}

// TODO:
// user authentication
// 	1. CLI cmd
// 	2. check if password exists
// 	3. protect routes
// do chunking on a client side so i can upload large files

//FIXES:
// creating files.db
