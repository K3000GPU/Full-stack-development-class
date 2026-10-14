package com.example.campus.controller;

import com.example.campus.dto.CreateMessageRequest;
import com.example.campus.dto.MessageResponse;
import com.example.campus.service.MessageService;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("/api/messages")
public class MessageController {
    private final MessageService service;

    public MessageController(MessageService service) {
        this.service = service;
    }

    // Controller 负责 HTTP 入口和响应转换，不执行 SQL 或业务校验。
    @GetMapping
    public List<MessageResponse> listMessages() {
        return service.listMessages().stream().map(MessageResponse::from).toList();
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public MessageResponse createMessage(@RequestBody CreateMessageRequest request) {
        return MessageResponse.from(service.createMessage(request.nickname(), request.content()));
    }
}
