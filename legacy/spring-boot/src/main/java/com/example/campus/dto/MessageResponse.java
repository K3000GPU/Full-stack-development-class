package com.example.campus.dto;

import com.example.campus.model.Message;

import java.time.LocalDateTime;

// 响应 DTO 定义前端可以看到的 JSON 字段。
public record MessageResponse(Long id, String nickname, String content, LocalDateTime createdAt) {
    public static MessageResponse from(Message message) {
        return new MessageResponse(message.id(), message.nickname(), message.content(), message.createdAt());
    }
}
