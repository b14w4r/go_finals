package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"

    "banking_service_project/services"
)

type Handler struct {
    authService     services.AuthService
    accountService  services.AccountService
    cardService     services.CardService
    transferService services.TransferService
    creditService   services.CreditService
    analyticsService services.AnalyticsService
    externalService services.ExternalService
}

func NewHandler(authS services.AuthService, accountS services.AccountService, cardS services.CardService, transferS services.TransferService, creditS services.CreditService, analyticsS services.AnalyticsService, externalS services.ExternalService) *Handler {
    return &Handler{
        authService:      authS,
        accountService:   accountS,
        cardService:      cardS,
        transferService:  transferS,
        creditService:    creditS,
        analyticsService: analyticsS,
        externalService:  externalS,
    }
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    type request struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    var req request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    user, err := h.authService.Register(req.Username, req.Email, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    type request struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    var req request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    token, err := h.authService.Login(req.Email, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
    userIDStr := r.Context().Value("userID").(string)
    userID, _ := strconv.Atoi(userIDStr)
    account, err := h.accountService.CreateAccount(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(account)
}

func (h *Handler) GetUserAccounts(w http.ResponseWriter, r *http.Request) {
    userIDStr := r.Context().Value("userID").(string)
    userID, _ := strconv.Atoi(userIDStr)
    accounts, err := h.accountService.GetUserAccounts(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(accounts)
}

func (h *Handler) CreateCard(w http.ResponseWriter, r *http.Request) {
    userIDStr := r.Context().Value("userID").(string)
    userID, _ := strconv.Atoi(userIDStr)
    // For simplicity, assume accountID is passed as query parameter
    accountIDStr := r.URL.Query().Get("account_id")
    if accountIDStr == "" {
        http.Error(w, "account_id is required", http.StatusBadRequest)
        return
    }
    accountID, _ := strconv.Atoi(accountIDStr)

    card, err := h.cardService.CreateCard(accountID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(card)
}

func (h *Handler) GetUserCards(w http.ResponseWriter, r *http.Request) {
    userIDStr := r.Context().Value("userID").(string)
    userID, _ := strconv.Atoi(userIDStr)
    // For simplicity, assume accountID is passed as query parameter
    accountIDStr := r.URL.Query().Get("account_id")
    if accountIDStr == "" {
        http.Error(w, "account_id is required", http.StatusBadRequest)
        return
    }
    accountID, _ := strconv.Atoi(accountIDStr)

    cards, err := h.cardService.GetCards(accountID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(cards)
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
    userIDStr := r.Context().Value("userID").(string)
    userID, _ := strconv.Atoi(userIDStr)

    type request struct {
        FromAccountID int     `json:"from_account_id"`
        ToAccountID   int     `json:"to_account_id"`
        Amount        float64 `json:"amount"`
    }
    var req request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    tx, err := h.transferService.Transfer(req.FromAccountID, req.ToAccountID, req.Amount)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(tx)
}

func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
    userIDStr := r.Context().Value("userID").(string)
    userID, _ := strconv.Atoi(userIDStr)
    // For simplicity, fixed current month and year
    stats, err := h.analyticsService.GetMonthlyStats(userID, int(time.Now().Month()), time.Now().Year())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(stats)
}

func (h *Handler) GetCreditSchedule(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    creditID, _ := strconv.Atoi(vars["creditId"])
    schedule, err := h.creditService.GetSchedule(creditID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(schedule)
}

func (h *Handler) PredictBalance(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    accountID, _ := strconv.Atoi(vars["accountId"])
    // For simplicity, days parameter as query
    daysStr := r.URL.Query().Get("days")
    days, _ := strconv.Atoi(daysStr)
    prediction, err := h.analyticsService.PredictBalance(accountID, days)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]float64{"prediction": prediction})
}

func (h *Handler) ApplyCredit(w http.ResponseWriter, r *http.Request) {
    type request struct {
        AccountID   int     `json:"account_id"`
        Principal   float64 `json:"principal"`
        AnnualRate  float64 `json:"annual_rate"`
        TermMonths  int     `json:"term_months"`
    }
    var req request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    credit, schedule, err := h.creditService.ApplyCredit(req.AccountID, req.Principal, req.AnnualRate, req.TermMonths)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "credit":   credit,
        "schedule": schedule,
    })
}
