package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"qaqmall/models"
)

// ProductHandler 商品处理器
type ProductHandler struct {
	db *gorm.DB
}

// NewProductHandler 创建商品处理器
func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

// ListProducts 获取商品列表
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var products []models.Product
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	category := c.Query("category")

	query := h.db.Model(&models.Product{})

	if category != "" {
		query = query.Where("categories LIKE ?", "%"+category+"%")
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"total":    total,
		"products": products,
	})
}

// GetProduct 获取单个商品
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := h.db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "商品不存在"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// SearchProducts 搜索商品
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "搜索关键词不能为空"})
		return
	}

	var products []models.Product
	h.db.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"results": products,
	})
}

// CreateProduct 创建商品
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求数据"})
		return
	}

	if err := h.db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建商品失败"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct 更新商品信息
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := h.db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "商品不存在"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求数据"})
		return
	}

	if err := h.db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新商品失败"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct 删除商品
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := h.db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "商品不存在"})
		return
	}

	if err := h.db.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除商品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "商品已删除"})
}
