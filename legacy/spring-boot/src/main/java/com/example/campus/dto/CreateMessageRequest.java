package com.example.campus.dto;

// 请求 DTO 只接收允许用户填写的字段，id 和时间由数据库生成。
public record CreateMessageRequest(String nickname, String content) {
}
