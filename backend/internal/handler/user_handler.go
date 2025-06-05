package handler

import (
	"backend/internal/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 定义请求体结构
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// 处理函数
// 创建user
func HandleCreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUserID, err := repository.CreateUser(req.Username, req.Password, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user_id": newUserID,
	})
}

// 从id获取user
func HandleGetUserID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := repository.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user", "details": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// 获取所有用户
func HandleGetAllUsers(c *gin.Context) {
	user, err := repository.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can not get users"})
		return
	}
	if user == nil {
		c.JSON(http.StatusOK, gin.H{})
	}
	c.JSON(http.StatusOK, user)
}

// 更新用户邮箱
func HandleUpdateUser(c *gin.Context) {
	type UpdateUserEmailRequest struct {
		Email string `json:"email" binding:"required,email"`
	}

	strID := c.Param("id")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var update UpdateUserEmailRequest
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入合法的邮箱地址"})
		return
	}
	
	rowsAffected, err := repository.UpdateUserEmail(id,update.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在或邮箱未修改"})
		return
	}

	c.JSON(http.StatusOK, rowsAffected)

}

func HandleDeleteUser(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.ParseInt(strID, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	rowsAffected, err := repository.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "删除失败,请输入正确的ID"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在或邮箱未修改"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("用户 %d 已删除", id)})
}