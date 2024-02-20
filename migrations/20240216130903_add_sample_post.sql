-- +goose Up
-- +goose StatementBegin
INSERT INTO posts(title, content, excerpt) VALUES(
    'My Very First Post',
    'This is the content',
    'Excerpt01'
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM posts ORDER BY id DESC LIMIT 1;
-- +goose StatementEnd
