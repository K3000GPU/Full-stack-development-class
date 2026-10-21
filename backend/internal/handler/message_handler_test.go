package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"campus-message-board/internal/model"
	"campus-message-board/internal/service"
	"github.com/gin-gonic/gin"
)

func TestCreateRejectsInvalidRequestsBeforeDatabaseAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// 单元测试只检查协议与规则；真实 MySQL 由启动后的接口验收覆盖。
	handler := New(service.New(nil))
	router := gin.New()
	router.POST("/api/messages", handler.Create)
	tests := []struct {
		name, body, contentType string
		status                  int
	}{
		{"missing fields", `{}`, "application/json", 400},
		{"null nickname", `{"nickname":null,"content":"留言"}`, "application/json", 400},
		{"number nickname", `{"nickname":123,"content":"留言"}`, "application/json", 400},
		{"boolean content", `{"nickname":"同学","content":true}`, "application/json", 400},
		{"object content", `{"nickname":"同学","content":{}}`, "application/json", 400},
		{"array nickname", `{"nickname":[],"content":"留言"}`, "application/json", 400},
		{"null body", `null`, "application/json", 400},
		{"invalid JSON", `{bad json`, "application/json", 400},
		{"empty body", ``, "application/json", 400},
		{"blank content", `{"nickname":"同学","content":" "}`, "application/json", 400},
		{"wrong media type", `nickname=test`, "text/plain", 415},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/messages", strings.NewReader(test.body))
			request.Header.Set("Content-Type", test.contentType)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			var body model.ErrorResponse
			if response.Code != test.status || json.Unmarshal(response.Body.Bytes(), &body) != nil || body.Message == "" {
				t.Fatalf("status=%d, body=%s", response.Code, response.Body.String())
			}
		})
	}
}
