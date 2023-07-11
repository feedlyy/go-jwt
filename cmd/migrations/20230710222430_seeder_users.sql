-- +goose Up
-- +goose StatementBegin
insert into users values
(1, 'admin', 'admin', '5cd753728ee2cc02ea3816e4b229a780b01560c3155cb25610f6c3681efb3b01', 'admin'), -- pwd: admin
(2, 'fadli', 'fadli', '925e23673d7cb416da43057958491693264db5bcd3bf519c56ec97ac917988de', 'user'); -- pwd: user
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- SELECT 'down SQL query';
-- +goose StatementEnd
