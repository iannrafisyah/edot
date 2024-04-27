-- +goose Up
create table warehouses (
   id            bigserial primary key,
   name          varchar(50),
   address       varchar,
   status        boolean default false,
   updated_at    timestamptz default now(),
   created_at    timestamptz default now(),
   deleted_at    timestamptz default null
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
drop table warehouses;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
