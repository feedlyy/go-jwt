-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id serial not null ,
    username text not null ,
    name text not null ,
    password text not null ,
    role text not null ,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP ,
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
