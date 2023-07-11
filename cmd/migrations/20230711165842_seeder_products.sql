-- +goose Up
-- +goose StatementBegin
insert into products values
(1, 'cap kaki tiga', 5, 'melegakan panas dalam'),
(2, 'fanta', 3, 'soda');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- SELECT 'down SQL query';
-- +goose StatementEnd
