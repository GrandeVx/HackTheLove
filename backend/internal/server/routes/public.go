package routes

import (
	"backend/config"
	"backend/internal/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
)

func registerPublicRoutes(r *gin.Engine, h *Handler) {
	public := r.Group("/")
	public.POST("/verifyGoogleJWT", h.VerifyGoogleToken)
	public.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"YES I AM UP <3!": 200}) })
	public.GET("/health", h.HealthHandler)
	public.GET("/questions", GetQuestions)
	public.GET("/time", h.GetTime)
}

func (h *Handler) GetTime(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"time": time.Now().Format(time.RFC3339)})
}

func GetQuestions(c *gin.Context) {
	var questionsModel models.Questions

	err := json.Unmarshal(config.Questions, &questionsModel)
	if err != nil {
		log.Error().Err(err).Msg("Errore nel parsing del JSON delle domande")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load questions"})
		return
	}

	c.JSON(http.StatusOK, questionsModel)
}

func (h *Handler) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, h.DB.Health())
}

func (h *Handler) VerifyGoogleToken(c *gin.Context) {
	var user models.User
	var body struct {
		IDToken string `json:"idToken"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		log.Warn().Err(err).Msg("Invalid JSON in Google token request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	log.Debug().Str("token", body.IDToken[:40]+"...").Msg("Verifying Google token")

	validator, err := idtoken.Validate(c, body.IDToken, config.GoogleClientID)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid Google JWT")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Google JWT"})
		return
	}

	email, emailOk := validator.Claims["email"].(string)
	familyName, familyNameOk := validator.Claims["family_name"].(string)
	givenName, givenNameOk := validator.Claims["given_name"].(string)

	if !emailOk || !familyNameOk || !givenNameOk {
		log.Warn().Msg("Missing required claims in Google JWT")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Google JWT payload"})
		return
	}

	user.Email = email
	user.FamilyName = familyName
	user.GivenName = givenName

	existingUser, errCheck := h.UserRepo.GetUser(user.Email)
	if errCheck != nil && existingUser == nil {
		if err := h.UserRepo.AddUser(&user); err != nil {
			log.Error().Err(err).Str("email", user.Email).Msg("Database error while adding user")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding user"})
			return
		}
	} else {
		log.Info().Str("email", user.Email).Msg("User already exists skipping creation")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":       email,
		"family_name": familyName,
		"given_name":  givenName,
		"exp":         time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	signedToken, err := token.SignedString(config.JwtSecret)
	if err != nil {
		log.Error().Err(err).Msg("Error signing JWT token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   signedToken,
	})
}
