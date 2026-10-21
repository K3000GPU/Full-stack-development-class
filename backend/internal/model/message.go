package model

import "time"

// Message 对应数据库的一条记录；请求和响应分别描述 HTTP 数据。
type Message struct {
	ID        int64
	Nickname  string
	Content   string
	CreatedAt time.Time
}

type CreateMessageRequest struct {
	Nickname string `json:"nickname"`
	Content  string `json:"content"`
}

type MessageResponse struct {
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

func NewMessageResponse(message Message) MessageResponse {
	return MessageResponse{
		ID: message.ID, Nickname: message.Nickname, Content: message.Content,
		CreatedAt: message.CreatedAt.Format("2006-01-02T15:04:05"),
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}
