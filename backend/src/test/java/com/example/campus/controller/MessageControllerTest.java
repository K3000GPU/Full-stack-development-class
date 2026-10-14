package com.example.campus.controller;

import com.example.campus.config.JacksonConfig;
import com.example.campus.exception.ApiExceptionHandler;
import com.example.campus.model.Message;
import com.example.campus.repository.MessageRepository;
import com.example.campus.service.MessageService;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.context.annotation.Import;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.http.MediaType;
import org.springframework.jdbc.CannotGetJdbcConnectionException;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;

import java.time.LocalDateTime;
import java.util.List;

import static org.hamcrest.Matchers.not;
import static org.hamcrest.Matchers.containsString;
import static org.mockito.Mockito.verifyNoInteractions;
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.content;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

// 测试 HTTP 协议、JSON 类型和错误格式；真实 MySQL 另用运行验收确认。
@WebMvcTest(MessageController.class)
@Import({MessageService.class, ApiExceptionHandler.class, JacksonConfig.class})
class MessageControllerTest {
    @Autowired
    private MockMvc mvc;
    @MockitoBean
    private MessageRepository repository;

    private final Message saved = new Message(7L, "同学", "你好", LocalDateTime.of(2026, 10, 13, 10, 0));

    @Test
    void getReturnsAnArray() throws Exception {
        when(repository.findAll()).thenReturn(List.of(saved));
        mvc.perform(get("/api/messages"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$[0].id").value(7))
                .andExpect(jsonPath("$[0].createdAt").value("2026-10-13T10:00:00"));
    }

    @Test
    void postReturnsCreatedDatabaseRecord() throws Exception {
        when(repository.insert("同学", "你好")).thenReturn(saved);
        mvc.perform(post("/api/messages").contentType(MediaType.APPLICATION_JSON)
                        .content("{\"nickname\":\"  同学  \",\"content\":\"你好\"}"))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.id").value(7))
                .andExpect(jsonPath("$.nickname").value("同学"))
                .andExpect(jsonPath("$.content").value("你好"))
                .andExpect(jsonPath("$.createdAt").value("2026-10-13T10:00:00"));
    }

    @ParameterizedTest
    @ValueSource(strings = {
            "{\"nickname\":123,\"content\":\"你好\"}",
            "{\"nickname\":\"同学\",\"content\":1.5}",
            "{\"nickname\":true,\"content\":\"你好\"}",
            "{\"nickname\":\"同学\",\"content\":{}}",
            "{\"nickname\":[],\"content\":\"你好\"}",
            "null", "{bad json}"
    })
    void rejectsNonStringAndMalformedJson(String json) throws Exception {
        mvc.perform(post("/api/messages").contentType(MediaType.APPLICATION_JSON).content(json))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").isString());
        verifyNoInteractions(repository);
    }

    @Test
    void invalidMessageReturnsBusinessError() throws Exception {
        mvc.perform(post("/api/messages").contentType(MediaType.APPLICATION_JSON)
                        .content("{\"nickname\":null,\"content\":\"你好\"}"))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("昵称不能为空"));
        verifyNoInteractions(repository);
    }

    @Test
    void unavailableDatabaseReturns503WithoutDetails() throws Exception {
        when(repository.findAll()).thenThrow(new CannotGetJdbcConnectionException("secret-database-details"));
        mvc.perform(get("/api/messages"))
                .andExpect(status().isServiceUnavailable())
                .andExpect(jsonPath("$.message").value("数据库暂时不可用，请稍后重试"))
                .andExpect(content().string(not(containsString("secret-database-details"))));
    }

    @Test
    void otherDatabaseFailureReturns500WithoutDetails() throws Exception {
        when(repository.findAll()).thenThrow(new DataIntegrityViolationException("secret-sql-details"));
        mvc.perform(get("/api/messages"))
                .andExpect(status().isInternalServerError())
                .andExpect(jsonPath("$.message").isString())
                .andExpect(content().string(not(containsString("secret-sql-details"))));
    }
}
