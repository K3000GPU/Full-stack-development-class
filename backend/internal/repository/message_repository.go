package repository

import (
	"context"
	"database/sql"

	"campus-message-board/internal/model"
)

type MessageRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Repository 只执行 SQL，并把数据库行读入 Message。
func (repository *MessageRepository) FindAll(ctx context.Context) ([]model.Message, error) {
	rows, err := repository.db.QueryContext(ctx, `
		SELECT id, nickname, content, created_at FROM messages ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages := make([]model.Message, 0)
	for rows.Next() {
		var message model.Message
		if err := rows.Scan(&message.ID, &message.Nickname, &message.Content, &message.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, rows.Err()
}

func (repository *MessageRepository) Insert(ctx context.Context, nickname, content string) (model.Message, error) {
	// 同一次发布中完成 INSERT 和回读；回读失败时不留下不确定的提交。
	transaction, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Message{}, err
	}
	defer transaction.Rollback()
	// ? 是参数占位符，输入作为值传入，不拼接到 SQL 中。
	result, err := transaction.ExecContext(ctx, `INSERT INTO messages (nickname, content) VALUES (?, ?)`, nickname, content)
	if err != nil {
		return model.Message{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return model.Message{}, err
	}
	var message model.Message
	err = transaction.QueryRowContext(ctx, `
		SELECT id, nickname, content, created_at FROM messages WHERE id = ?
	`, id).Scan(&message.ID, &message.Nickname, &message.Content, &message.CreatedAt)
	if err != nil {
		return model.Message{}, err
	}
	if err := transaction.Commit(); err != nil {
		return model.Message{}, err
	}
	return message, nil
}
