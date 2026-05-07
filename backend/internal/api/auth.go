package api

import (
	"fmt"
	"strings"

	"devhelper/internal/service"
	"devhelper/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
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
	user, tokens, err := h.svc.Register(req.Username, req.Email, req.Password)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"user": user, "tokens": tokens})
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
	user, tokens, userNotFound, err := h.svc.Login(req.UserName, req.Password)
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
	tokens, err := h.svc.Refresh(req.RefreshToken)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}
	utils.OK(c, tokens)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.svc.GetUser(userID.(uint))
	if err != nil {
		utils.NotFound(c, "user not found")
		return
	}
	utils.OK(c, user)
}
