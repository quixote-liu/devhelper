package api

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lcs/devhelper/internal/models"
	"github.com/lcs/devhelper/internal/repository"
	"github.com/lcs/devhelper/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminHandler struct {
	userRepo *repository.UserRepo
}

func NewAdminHandler(userRepo *repository.UserRepo) *AdminHandler {
	return &AdminHandler{userRepo: userRepo}
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	users, total, err := h.userRepo.List(page, pageSize, search)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"users": users, "total": total, "page": page, "page_size": pageSize})
}

func (h *AdminHandler) GetUser(c *gin.Context) {
	var id uint
	if _, err := parseUintParam(c, "id", &id); err != nil {
		return
	}
	user, err := h.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFound(c, "user not found")
		} else {
			utils.InternalError(c, err.Error())
		}
		return
	}
	utils.OK(c, user)
}

type updateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	var id uint
	if _, err := parseUintParam(c, "id", &id); err != nil {
		return
	}
	user, err := h.userRepo.FindByID(id)
	if err != nil {
		utils.NotFound(c, "user not found")
		return
	}
	var req updateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role == "admin" || req.Role == "user" {
		user.Role = req.Role
	}
	if err := h.userRepo.Update(user); err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, user)
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	var id uint
	if _, err := parseUintParam(c, "id", &id); err != nil {
		return
	}
	// Prevent deleting yourself
	selfID, _ := c.Get("user_id")
	if selfID.(uint) == id {
		utils.BadRequest(c, "cannot delete yourself")
		return
	}
	if err := h.userRepo.Delete(id); err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, nil)
}

type resetPasswordReq struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (h *AdminHandler) ResetPassword(c *gin.Context) {
	var id uint
	if _, err := parseUintParam(c, "id", &id); err != nil {
		return
	}
	user, err := h.userRepo.FindByID(id)
	if err != nil {
		utils.NotFound(c, "user not found")
		return
	}
	var req resetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	user.PasswordHash = string(hash)
	if err := h.userRepo.Update(user); err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"message": "password reset successfully"})
}

// UserProfile allows users to update their own profile
type updateProfileReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func UpdateProfile(userRepo *repository.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		user, err := userRepo.FindByID(userID.(uint))
		if err != nil {
			utils.NotFound(c, "user not found")
			return
		}
		var req updateProfileReq
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, err.Error())
			return
		}
		if req.Username != "" {
			user.Username = req.Username
		}
		if req.Email != "" {
			user.Email = req.Email
		}
		if err := userRepo.Update(user); err != nil {
			utils.InternalError(c, err.Error())
			return
		}
		utils.OK(c, user)
	}
}

// ChangePassword allows users to change their own password
func ChangePassword(userRepo *repository.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		user, err := userRepo.FindByID(userID.(uint))
		if err != nil {
			utils.NotFound(c, "user not found")
			return
		}
		var req struct {
			OldPassword string `json:"old_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required,min=6"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, err.Error())
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
			utils.BadRequest(c, "incorrect old password")
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
		if err != nil {
			utils.InternalError(c, err.Error())
			return
		}
		user.PasswordHash = string(hash)
		if err := userRepo.Update(user); err != nil {
			utils.InternalError(c, err.Error())
			return
		}
		utils.OK(c, gin.H{"message": "password changed successfully"})
	}
}

// DeleteAccount allows users to delete their own account
func DeleteAccount(userRepo *repository.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		var req struct {
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, err.Error())
			return
		}
		user, err := userRepo.FindByID(userID.(uint))
		if err != nil {
			utils.NotFound(c, "user not found")
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			utils.BadRequest(c, "incorrect password")
			return
		}
		if err := userRepo.Delete(userID.(uint)); err != nil {
			utils.InternalError(c, err.Error())
			return
		}
		utils.OK(c, gin.H{"message": "account deleted"})
	}
}

func init() {
	_ = models.User{}
}
