package com.example.campus.repository;

import com.example.campus.model.Message;
import org.springframework.dao.DataRetrievalFailureException;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.jdbc.support.GeneratedKeyHolder;
import org.springframework.stereotype.Repository;

import java.sql.Statement;
import java.util.List;

@Repository
public class MessageRepository {
    private final JdbcTemplate jdbcTemplate;

    private static final RowMapper<Message> ROW_MAPPER = (resultSet, rowNumber) -> new Message(
            resultSet.getLong("id"),
            resultSet.getString("nickname"),
            resultSet.getString("content"),
            resultSet.getTimestamp("created_at").toLocalDateTime()
    );

    public MessageRepository(JdbcTemplate jdbcTemplate) {
        this.jdbcTemplate = jdbcTemplate;
    }

    // Repository 负责执行 SQL，把数据库行转换为 Model。
    public List<Message> findAll() {
        return jdbcTemplate.query("""
                SELECT id, nickname, content, created_at
                FROM messages
                ORDER BY created_at DESC, id DESC
                """, ROW_MAPPER);
    }

    public Message insert(String nickname, String content) {
        var keyHolder = new GeneratedKeyHolder();
        jdbcTemplate.update(connection -> {
            // ? 是参数占位符：用户输入作为数据传入，不拼接进 SQL。
            var statement = connection.prepareStatement("""
                    INSERT INTO messages (nickname, content) VALUES (?, ?)
                    """, Statement.RETURN_GENERATED_KEYS);
            statement.setString(1, nickname);
            statement.setString(2, content);
            return statement;
        }, keyHolder);

        Number id = keyHolder.getKey();
        if (id == null) {
            throw new DataRetrievalFailureException("新增留言未返回 id");
        }
        // 再查询一次，返回数据库实际生成的 id 和 created_at。
        return jdbcTemplate.queryForObject("""
                SELECT id, nickname, content, created_at FROM messages WHERE id = ?
                """, ROW_MAPPER, id.longValue());
    }
}
