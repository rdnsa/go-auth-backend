package http

import (
	"go-auth-backend/internal/dto"
	"go-auth-backend/internal/usecase"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase usecase.UserUsecase
}

func NewHandler(usecase usecase.UserUsecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.usecase.Register(c.Request.Context(), &req)
	if err != nil {
		log.Printf("Register gagal | Email: %s | Error: %v", req.Email, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// INI YANG BARU: LOG SUKSES
	log.Printf("REGISTER BERHASIL | Name: %s | Email: %s | UserID: %s",
		res.Name, res.Email, res.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Register berhasil!",
		"data":    res,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.usecase.Login(c.Request.Context(), &req)
	if err != nil {
		log.Printf("LOGIN GAGAL | Email: %s | Error: %v", req.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// INI YANG BARU: LOG SUKSES
	log.Printf("LOGIN BERHASIL | Email: %s | UserID: %s", res.User.Email, res.User.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"data":    res,
	})
}
