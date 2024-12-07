package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bootdotdev/learn-cicd-starter/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	// Decode request body
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters: "+err.Error())
		return
	}

	// Validate input
	if params.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Name cannot be empty")
		return
	}

	// Generate API key
	apiKey, err := generateRandomSHA256Hash()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate API key: "+err.Error())
		return
	}

	// Create user in database
	err = cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New().String(),
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		Name:      params.Name,
		ApiKey:    apiKey,
	})

	if err != nil {
		log.Printf("Error creating user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	// Retrieve created user
	user, err := cfg.DB.GetUser(r.Context(), apiKey)
	if err != nil {
		log.Printf("Error retrieving user after creation: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	// Convert user to response format
	userResp, err := databaseUserToUser(user)
	if err != nil {
		log.Printf("Error converting user to response format: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert user")
		return
	}

	// Respond with created user information
	respondWithJSON(w, http.StatusCreated, userResp)
}

func generateRandomSHA256Hash() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(randomBytes)
	hashString := hex.EncodeToString(hash[:])
	return hashString, nil
}

func (cfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {

	userResp, err := databaseUserToUser(user)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert user")
		return
	}

	respondWithJSON(w, http.StatusOK, userResp)
}
