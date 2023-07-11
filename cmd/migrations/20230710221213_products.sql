-- +goose Up
-- +goose StatementBegin
CREATE TABLE products(
    id serial not null ,
    name text not null ,
    qty int not null ,
    description text not null ,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP ,
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products;
-- +goose StatementEnd
