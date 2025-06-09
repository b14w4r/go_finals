package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gorilla/mux"
    _ "github.com/lib/pq"

    "banking_service_project/handlers"
    "banking_service_project/middleware"
    "banking_service_project/repositories"
    "banking_service_project/services"
)

func main() {
    // Load environment variables
    dbURL := os.Getenv("DATABASE_URL")
    jwtSecret := os.Getenv("JWT_SECRET")
    pgpPrivateKeyPath := os.Getenv("PGP_PRIVATE_KEY_PATH")
    pgpPublicKeyPath := os.Getenv("PGP_PUBLIC_KEY_PATH")
    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")
    smtpUser := os.Getenv("SMTP_USER")
    smtpPass := os.Getenv("SMTP_PASS")

    if dbURL == "" || jwtSecret == "" {
        log.Fatal("DATABASE_URL and JWT_SECRET must be set")
    }

    // Connect to Postgres
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }
    defer db.Close()

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    // Initialize repositories
    userRepo := repositories.NewUserRepository(db)
    accountRepo := repositories.NewAccountRepository(db)
    cardRepo := repositories.NewCardRepository(db)
    transactionRepo := repositories.NewTransactionRepository(db)
    creditRepo := repositories.NewCreditRepository(db)
    scheduleRepo := repositories.NewPaymentScheduleRepository(db)

    // Initialize services
    authService := services.NewAuthService(userRepo, jwtSecret)
    accountService := services.NewAccountService(accountRepo, transactionRepo)
    cardService := services.NewCardService(cardRepo, accountRepo, pgpPublicKeyPath, pgpPrivateKeyPath)
    transferService := services.NewTransferService(accountRepo, transactionRepo)
    creditService := services.NewCreditService(creditRepo, scheduleRepo, accountRepo)
    analyticsService := services.NewAnalyticsService(transactionRepo)
    externalService := services.NewExternalService(smtpHost, smtpPort, smtpUser, smtpPass, pgpPublicKeyPath, pgpPrivateKeyPath)

    // Initialize handlers
    h := handlers.NewHandler(authService, accountService, cardService, transferService, creditService, analyticsService, externalService)

    // Setup router
    r := mux.NewRouter()

    // Public routes
    r.HandleFunc("/register", h.Register).Methods("POST")
    r.HandleFunc("/login", h.Login).Methods("POST")

    // Protected routes
    authRouter := r.PathPrefix("/").Subrouter()
    authRouter.Use(middleware.AuthMiddleware(jwtSecret))

    authRouter.HandleFunc("/accounts", h.CreateAccount).Methods("POST")
    authRouter.HandleFunc("/accounts", h.GetUserAccounts).Methods("GET")
    authRouter.HandleFunc("/cards", h.CreateCard).Methods("POST")
    authRouter.HandleFunc("/cards", h.GetUserCards).Methods("GET")
    authRouter.HandleFunc("/transfer", h.Transfer).Methods("POST")
    authRouter.HandleFunc("/analytics", h.GetAnalytics).Methods("GET")
    authRouter.HandleFunc("/credits/{creditId}/schedule", h.GetCreditSchedule).Methods("GET")
    authRouter.HandleFunc("/accounts/{accountId}/predict", h.PredictBalance).Methods("GET")
    authRouter.HandleFunc("/credits/apply", h.ApplyCredit).Methods("POST")

    // Start server
    serverPort := os.Getenv("PORT")
    if serverPort == "" {
        serverPort = "8080"
    }
    srv := &http.Server{
        Handler:      r,
        Addr:         ":" + serverPort,
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    log.Printf("Starting server on port %s...", serverPort)
    if err := srv.ListenAndServe(); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
