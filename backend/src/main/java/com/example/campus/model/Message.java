package com.example.campus.model;

import java.time.LocalDateTime;

// Model 对应数据库中的一条留言；record 自动提供只读字段和访问方法。
public record Message(Long id, String nickname, String content, LocalDateTime createdAt) {
}
