package com.example.campus.exception;

import com.example.campus.dto.ErrorResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.dao.DataAccessException;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.HttpStatusCode;
import org.springframework.http.ResponseEntity;
import org.springframework.jdbc.CannotGetJdbcConnectionException;
import org.springframework.transaction.CannotCreateTransactionException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;
import org.springframework.web.context.request.WebRequest;
import org.springframework.web.servlet.mvc.method.annotation.ResponseEntityExceptionHandler;

// 统一错误处理：对外只返回可读提示，不发送 SQL、密码或异常堆栈。
@RestControllerAdvice
public class ApiExceptionHandler extends ResponseEntityExceptionHandler {
    private static final Logger LOGGER = LoggerFactory.getLogger(ApiExceptionHandler.class);

    @ExceptionHandler(InvalidMessageException.class)
    public ResponseEntity<ErrorResponse> handleInvalidMessage(InvalidMessageException exception) {
        return ResponseEntity.badRequest().body(new ErrorResponse(exception.getMessage()));
    }

    @ExceptionHandler({CannotGetJdbcConnectionException.class, CannotCreateTransactionException.class})
    public ResponseEntity<ErrorResponse> handleDatabaseUnavailable(Exception exception) {
        LOGGER.error("教学数据库连接失败，异常类型：{}", exception.getClass().getSimpleName());
        return ResponseEntity.status(HttpStatus.SERVICE_UNAVAILABLE)
                .body(new ErrorResponse("数据库暂时不可用，请稍后重试"));
    }

    @ExceptionHandler(DataAccessException.class)
    public ResponseEntity<ErrorResponse> handleDatabaseError(DataAccessException exception) {
        LOGGER.error("数据库操作失败，异常类型：{}", exception.getClass().getSimpleName());
        return ResponseEntity.internalServerError().body(new ErrorResponse("留言读取或保存失败，请稍后重试"));
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<ErrorResponse> handleUnexpectedError(Exception exception) {
        LOGGER.error("处理留言时发生异常，异常类型：{}", exception.getClass().getSimpleName());
        return ResponseEntity.internalServerError().body(new ErrorResponse("服务暂时出错，请稍后重试"));
    }

    @Override
    protected ResponseEntity<Object> handleExceptionInternal(Exception exception, Object body,
            HttpHeaders headers, HttpStatusCode status, WebRequest request) {
        String message = switch (status.value()) {
            case 400 -> "请求格式不正确，nickname 和 content 必须为字符串";
            case 404 -> "接口不存在";
            case 405 -> "请求方法不支持";
            case 415 -> "请使用 application/json 发送留言";
            default -> "请求无法处理，请检查后重试";
        };
        return super.handleExceptionInternal(exception, new ErrorResponse(message), headers, status, request);
    }
}
