-- +goose Up
create table products (
    id            bigserial primary key,
    name          varchar(50),
    price         float,
    updated_at    timestamptz default now(),
    created_at    timestamptz default now(),
    deleted_at    timestamptz default null
);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
drop table products;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
