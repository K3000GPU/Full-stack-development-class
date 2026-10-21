package service

import (
	"context"
	"strings"
	"unicode/utf8"

	"campus-message-board/internal/model"
	"campus-message-board/internal/repository"
)

type ValidationError struct {
	Message string
}

func (err *ValidationError) Error() string { return err.Message }

type MessageService struct {
	repository *repository.MessageRepository
}

func New(repository *repository.MessageRepository) *MessageService {
	return &MessageService{repository: repository}
}

func (service *MessageService) List(ctx context.Context) ([]model.Message, error) {
	return service.repository.FindAll(ctx)
}

// Service 处理规则与业务过程；每层只连接一个具体实现。
func (service *MessageService) Create(ctx context.Context, nickname, content string) (model.Message, error) {
	cleanNickname, err := validateInput(nickname, content)
	if err != nil {
		return model.Message{}, err
	}
	return service.repository.Insert(ctx, cleanNickname, content)
}

func validateInput(nickname, content string) (string, error) {
	cleanNickname := trimWhitespace(nickname)
	if cleanNickname == "" {
		return "", &ValidationError{Message: "昵称不能为空"}
	}
	if utf8.RuneCountInString(cleanNickname) > 30 {
		return "", &ValidationError{Message: "昵称最多 30 个字符"}
	}
	if trimWhitespace(content) == "" {
		return "", &ValidationError{Message: "留言内容不能为空"}
	}
	if utf8.RuneCountInString(content) > 200 {
		return "", &ValidationError{Message: "留言内容最多 200 个字符"}
	}
	return cleanNickname, nil
}

func trimWhitespace(value string) string {
	// 与前端 String.trim() 的空白集合一致，包含 NBSP 和 FEFF。
	const whitespace = "\t\n\v\f\r \u00a0\u1680\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200a\u2028\u2029\u202f\u205f\u3000\ufeff"
	return strings.Trim(value, whitespace)
}
