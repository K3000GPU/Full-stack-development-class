package handler

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"log"
	"net"
	"net/http"

	"campus-message-board/internal/model"
	"campus-message-board/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type MessageHandler struct {
	service *service.MessageService
}

func New(service *service.MessageService) *MessageHandler {
	return &MessageHandler{service: service}
}

// Handler 负责接收 HTTP、解析 JSON、调用业务方法和返回响应。
func (handler *MessageHandler) List(c *gin.Context) {
	messages, err := handler.service.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	response := make([]model.MessageResponse, 0, len(messages))
	for _, message := range messages {
		response = append(response, model.NewMessageResponse(message))
	}
	c.JSON(http.StatusOK, response)
}

func (handler *MessageHandler) Create(c *gin.Context) {
	if c.ContentType() != "application/json" {
		c.JSON(http.StatusUnsupportedMediaType, model.ErrorResponse{Message: "请使用 application/json 发送留言"})
		return
	}
	var request model.CreateMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "请求格式不正确，nickname 和 content 必须为字符串"})
		return
	}
	message, err := handler.service.Create(c.Request.Context(), request.Nickname, request.Content)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, model.NewMessageResponse(message))
}

func respondError(c *gin.Context, err error) {
	var validationError *service.ValidationError
	if errors.As(err, &validationError) {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: validationError.Message})
		return
	}
	// 详细错误只留在服务端日志，响应不包含 SQL、凭据或异常堆栈。
	log.Printf("数据库操作失败：%v", err)
	var networkError net.Error
	if errors.As(err, &networkError) || errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, driver.ErrBadConn) || errors.Is(err, mysql.ErrInvalidConn) || errors.Is(err, sql.ErrConnDone) {
		c.JSON(http.StatusServiceUnavailable, model.ErrorResponse{Message: "数据库暂时不可用，请稍后重试"})
		return
	}
	c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "留言读取或保存失败，请稍后重试"})
}
