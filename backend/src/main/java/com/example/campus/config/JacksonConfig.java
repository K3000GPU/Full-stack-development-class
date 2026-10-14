package com.example.campus.config;

import com.fasterxml.jackson.databind.cfg.CoercionAction;
import com.fasterxml.jackson.databind.cfg.CoercionInputShape;
import com.fasterxml.jackson.databind.type.LogicalType;
import org.springframework.boot.autoconfigure.jackson.Jackson2ObjectMapperBuilderCustomizer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class JacksonConfig {
    // JSON 中的昵称与内容必须是字符串，不能把数字或布尔值自动转成文字。
    @Bean
    Jackson2ObjectMapperBuilderCustomizer strictStrings() {
        return builder -> builder.postConfigurer(mapper -> {
            var strings = mapper.coercionConfigFor(LogicalType.Textual);
            strings.setCoercion(CoercionInputShape.Integer, CoercionAction.Fail);
            strings.setCoercion(CoercionInputShape.Float, CoercionAction.Fail);
            strings.setCoercion(CoercionInputShape.Boolean, CoercionAction.Fail);
        });
    }
}
