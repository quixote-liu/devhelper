package api

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"devhelper/internal/config"
	"devhelper/internal/models"
	"devhelper/internal/repository"
	"devhelper/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthHandler(userRepo repository.UserRepository, cfg *config.Config) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, cfg: cfg}
}

func (h *AuthHandler) generateTokens(user *models.User) (*tokenPair, error) {
	access, err := utils.GenerateAccessToken(user.ID, user.Email, user.Role, h.cfg.JWTSecret, h.cfg.JWTAccessExpiry)
	if err != nil {
		return nil, err
	}
	refresh, err := utils.GenerateRefreshToken(user.ID, user.Email, user.Role, h.cfg.JWTSecret, h.cfg.JWTRefreshExpiry)
	if err != nil {
		return nil, err
	}
	return &tokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

type tokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) checkAndRegisterUser(username, email, password string) (*models.User, *tokenPair, error) {
	if _, err := h.userRepo.FindByEmail(email); err == nil {
		return nil, nil, errors.New("email already registered")
	}
	if _, err := h.userRepo.FindByUsername(username); err == nil {
		return nil, nil, errors.New("username already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, nil, err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		Role:         "user",
	}
	if err := h.userRepo.Create(user); err != nil {
		return nil, nil, err
	}

	tokens, err := h.generateTokens(user)
	return user, tokens, err
}

func getFieldName(field string) string {
	switch field {
	case "Username":
		return "用户名"
	case "Email":
		return "邮箱"
	case "Password":
		return "密码"
	default:
		return field
	}
}

// 将 validator 错误转换为友好的中文提示
func formatValidationError(err error) string {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range validationErrs {
			field := e.Field()
			tag := e.Tag()

			var msg string
			switch tag {
			case "required":
				msg = fmt.Sprintf("%s不能为空", getFieldName(field))
			case "email":
				msg = "邮箱格式不正确"
			case "min":
				msg = fmt.Sprintf("%s长度至少%s个字符", getFieldName(field), e.Param())
			case "max":
				msg = fmt.Sprintf("%s长度不能超过%s个字符", getFieldName(field), e.Param())
			default:
				msg = fmt.Sprintf("%s格式不正确", getFieldName(field))
			}
			messages = append(messages, msg)
		}
		return strings.Join(messages, "；")
	}
	return err.Error()
}

type registerReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, formatValidationError(err))
		return
	}
	user, tokens, err := h.checkAndRegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"user": user, "tokens": tokens})
}

func (h *AuthHandler) loginUser(username, password string) (*models.User, *tokenPair, bool, error) {
	user, err := h.userRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, true, errors.New("invalid credentials")
		}
		return nil, nil, false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		fmt.Println("compare password error = ", err)
		return nil, nil, true, errors.New("invalid credentials")
	}

	now := time.Now()
	user.LastLogin = &now
	_ = h.userRepo.Update(user)

	tokens, err := h.generateTokens(user)
	return user, tokens, false, err
}

type loginReq struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// 登录接口不暴露具体字段错误，统一返回通用错误
		fmt.Println("登录参数错误: ", err)
		utils.InternalError(c, "登录参数错误")
		return
	}
	user, tokens, userNotFound, err := h.loginUser(req.UserName, req.Password)
	if userNotFound {
		utils.BadRequest(c, "账号或密码错误")
		return
	}
	if err != nil {
		fmt.Println("登录用户发生错误: ", err)
		utils.InternalError(c, "登录用户发生错误")
		return
	}
	utils.OK(c, gin.H{"user": user, "tokens": tokens})
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	claims, err := utils.ParseToken(req.RefreshToken, h.cfg.JWTSecret)
	if err != nil || claims.Type != "refresh" {
		utils.Unauthorized(c, "invalid refresh token")
		return
	}
	user, err := h.userRepo.FindByID(claims.UserID)
	if err != nil {
		utils.Unauthorized(c, "user not found")
		return
	}
	tokenPair, err := h.generateTokens(user)
	if err != nil {
		utils.Unauthorized(c, "internal service failed")
		return
	}

	utils.OK(c, tokenPair)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.userRepo.FindByID(userID.(uint))
	if err != nil {
		utils.NotFound(c, "user not found")
		return
	}
	utils.OK(c, user)
}
