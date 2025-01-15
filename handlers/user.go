package handlers

import (
	"net/http"
	"qaqmall/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type LoginResponse struct {
	UserID uint `json:"user_id"`
	Role   int  `json:"role"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求参数"})
		return
	}

	var user models.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户名或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户名或密码错误"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		UserID: user.ID,
		Role:   user.Role,
	})
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求参数"})
		return
	}

	var existingUser models.User
	if err := h.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "用户名已存在"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器错误"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     0,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"user_id": user.ID,
	})
}
