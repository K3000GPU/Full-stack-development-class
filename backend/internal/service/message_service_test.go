package service

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestValidateInput(t *testing.T) {
	tests := []struct {
		name, nickname, content, wantNickname, wantError string
	}{
		{"trim", "  同学  ", " 留言\n", "同学", ""},
		{"unicode whitespace", "\u00a0\ufeff同学\u3000", "留言", "同学", ""},
		{"emoji boundary", strings.Repeat("😀", 30), strings.Repeat("🌱", 200), strings.Repeat("😀", 30), ""},
		{"Chinese boundary", strings.Repeat("名", 30), strings.Repeat("字", 200), strings.Repeat("名", 30), ""},
		{"empty nickname", "", "留言", "", "昵称不能为空"},
		{"blank nickname", "\t \u00a0\ufeff", "留言", "", "昵称不能为空"},
		{"long nickname", strings.Repeat("名", 31), "留言", "", "昵称最多 30 个字符"},
		{"long emoji nickname", strings.Repeat("😀", 31), "留言", "", "昵称最多 30 个字符"},
		{"empty content", "同学", "", "", "留言内容不能为空"},
		{"blank content", "同学", " \n\u00a0\ufeff", "", "留言内容不能为空"},
		{"long content", "同学", strings.Repeat("字", 201), "", "留言内容最多 200 个字符"},
		{"long emoji content", "同学", strings.Repeat("🌱", 201), "", "留言内容最多 200 个字符"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nickname, err := validateInput(test.nickname, test.content)
			if test.wantError != "" {
				var invalid *ValidationError
				if !errors.As(err, &invalid) || invalid.Message != test.wantError {
					t.Fatalf("want %q, got %v", test.wantError, err)
				}
				return
			}
			if err != nil || nickname != test.wantNickname {
				t.Fatalf("nickname=%q, err=%v", nickname, err)
			}
		})
	}
}

func TestInvalidInputDoesNotAccessRepository(t *testing.T) {
	// 不提供 Repository：无效请求必须在进入数据库之前被拒绝。
	service := New(nil)
	if _, err := service.Create(context.Background(), " ", "留言"); err == nil {
		t.Fatal("expected validation error")
	}
}
