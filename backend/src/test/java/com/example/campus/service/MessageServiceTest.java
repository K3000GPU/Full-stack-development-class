package com.example.campus.service;

import com.example.campus.exception.InvalidMessageException;
import com.example.campus.model.Message;
import com.example.campus.repository.MessageRepository;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.NullAndEmptySource;
import org.junit.jupiter.params.provider.ValueSource;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.time.LocalDateTime;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.verifyNoInteractions;
import static org.mockito.Mockito.when;

// 单元测试只验证业务规则；Mock Repository 不代表真实 MySQL 已通过验收。
@ExtendWith(MockitoExtension.class)
class MessageServiceTest {
    @Mock
    private MessageRepository repository;
    @InjectMocks
    private MessageService service;

    @Test
    void trimsNicknameButPreservesContent() {
        String content = "  第一行\n第二行  ";
        var saved = new Message(1L, "新同学", content, LocalDateTime.of(2026, 10, 13, 10, 0));
        when(repository.insert("新同学", content)).thenReturn(saved);

        assertThat(service.createMessage(" \u00A0\uFEFF新同学\u3000\n", content)).isEqualTo(saved);
        verify(repository).insert("新同学", content);
    }

    @Test
    void acceptsExactUnicodeCharacterLimits() {
        String nickname = "😀".repeat(30);
        String content = "😀".repeat(200);
        var saved = new Message(1L, nickname, content, LocalDateTime.now());
        when(repository.insert(nickname, content)).thenReturn(saved);

        assertThat(service.createMessage(" " + nickname + " ", content)).isEqualTo(saved);
    }

    @ParameterizedTest
    @NullAndEmptySource
    @ValueSource(strings = {" ", "\t\n", "\u3000", "\u00A0\uFEFF"})
    void rejectsBlankNicknameBeforeWriting(String nickname) {
        assertThatThrownBy(() -> service.createMessage(nickname, "你好"))
                .isInstanceOf(InvalidMessageException.class).hasMessage("昵称不能为空");
        verifyNoInteractions(repository);
    }

    @ParameterizedTest
    @NullAndEmptySource
    @ValueSource(strings = {" ", "\t\n", "\u3000", "\u00A0\uFEFF"})
    void rejectsBlankContentBeforeWriting(String content) {
        assertThatThrownBy(() -> service.createMessage("同学", content))
                .isInstanceOf(InvalidMessageException.class).hasMessage("留言内容不能为空");
        verifyNoInteractions(repository);
    }

    @Test
    void rejectsNicknameOverThirtyCodePoints() {
        assertThatThrownBy(() -> service.createMessage("😀".repeat(31), "你好"))
                .isInstanceOf(InvalidMessageException.class).hasMessage("昵称最多 30 个字符");
        verifyNoInteractions(repository);
    }

    @Test
    void rejectsContentOverTwoHundredCodePointsIncludingSpaces() {
        assertThatThrownBy(() -> service.createMessage("同学", "😀".repeat(200) + " "))
                .isInstanceOf(InvalidMessageException.class).hasMessage("留言内容最多 200 个字符");
        verifyNoInteractions(repository);
    }
}
