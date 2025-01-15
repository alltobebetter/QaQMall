package handlers

import (
	"net/http"
	"qaqmall/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartHandler struct {
	db *gorm.DB
}

func NewCartHandler(db *gorm.DB) *CartHandler {
	return &CartHandler{db: db}
}

// AddItemRequest 添加购物车项请求
type AddItemRequest struct {
	UserID uint `json:"user_id"`
	Item   struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	} `json:"item"`
}

// GetCart 获取购物车
func (h *CartHandler) GetCart(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Query("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的用户ID"})
		return
	}

	var items []models.CartItem
	if err := h.db.Preload("Product").Where("user_id = ?", userID).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取购物车失败"})
		return
	}

	c.JSON(http.StatusOK, models.Cart{
		UserID: uint(userID),
		Items:  items,
	})
}

// AddItem 添加商品到购物车
func (h *CartHandler) AddItem(c *gin.Context) {
	var req AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求参数"})
		return
	}

	// 检查商品是否存在
	var product models.Product
	if err := h.db.First(&product, req.Item.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "商品不存在"})
		return
	}

	// 检查是否已在购物车中
	var existingItem models.CartItem
	result := h.db.Where("user_id = ? AND product_id = ?", req.UserID, req.Item.ProductID).First(&existingItem)

	if result.Error == nil {
		// 更新数量
		existingItem.Quantity += req.Item.Quantity
		if err := h.db.Save(&existingItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "更新购物车失败"})
			return
		}
	} else {
		// 创建新项
		cartItem := models.CartItem{
			UserID:    req.UserID,
			ProductID: req.Item.ProductID,
			Quantity:  req.Item.Quantity,
		}
		if err := h.db.Create(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "添加到购物车失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "添加成功"})
}

// UpdateItem 更新购物车项数量
func (h *CartHandler) UpdateItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的购物车项ID"})
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求参数"})
		return
	}

	var cartItem models.CartItem
	if err := h.db.First(&cartItem, itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "购物车项不存在"})
		return
	}

	cartItem.Quantity = req.Quantity
	if err := h.db.Save(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// RemoveItem 从购物车中移除商品
func (h *CartHandler) RemoveItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的购物车项ID"})
		return
	}

	if err := h.db.Delete(&models.CartItem{}, itemID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// EmptyCart 清空购物车
func (h *CartHandler) EmptyCart(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Query("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的用户ID"})
		return
	}

	if err := h.db.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "清空购物车失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "购物车已清空"})
}
