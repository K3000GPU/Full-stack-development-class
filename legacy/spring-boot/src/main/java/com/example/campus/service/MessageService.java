package com.example.campus.service;

import com.example.campus.exception.InvalidMessageException;
import com.example.campus.model.Message;
import com.example.campus.repository.MessageRepository;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.regex.Pattern;

@Service
public class MessageService {
    private final MessageRepository repository;
    // 与前端 String.trim() 使用同一空白集合，包含 NBSP 和 FEFF。
    private static final Pattern EDGE_WHITESPACE = Pattern.compile(
            "^[\\u0009-\\u000D\\u0020\\u00A0\\u1680\\u2000-\\u200A\\u2028\\u2029\\u202F\\u205F\\u3000\\uFEFF]+|"
            + "[\\u0009-\\u000D\\u0020\\u00A0\\u1680\\u2000-\\u200A\\u2028\\u2029\\u202F\\u205F\\u3000\\uFEFF]+$");

    public MessageService(MessageRepository repository) {
        this.repository = repository;
    }

    public List<Message> listMessages() {
        return repository.findAll();
    }

    // Service 负责规则与业务过程；事务保证 INSERT 和回读属于同一次发布。
    @Transactional
    public Message createMessage(String nickname, String content) {
        if (nickname == null || trimWhitespace(nickname).isEmpty()) {
            throw new InvalidMessageException("昵称不能为空");
        }
        String cleanNickname = trimWhitespace(nickname);
        if (characterCount(cleanNickname) > 30) {
            throw new InvalidMessageException("昵称最多 30 个字符");
        }
        if (content == null || trimWhitespace(content).isEmpty()) {
            throw new InvalidMessageException("留言内容不能为空");
        }
        if (characterCount(content) > 200) {
            throw new InvalidMessageException("留言内容最多 200 个字符");
        }
        // 昵称去除首尾空白；内容保留用户输入的换行与空白。
        return repository.insert(cleanNickname, content);
    }

    private int characterCount(String value) {
        // 按 Unicode code point 计数，常见 emoji 不会被算作两个字符。
        return value.codePointCount(0, value.length());
    }

    private String trimWhitespace(String value) {
        return EDGE_WHITESPACE.matcher(value).replaceAll("");
    }
}
