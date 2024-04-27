-- +goose Up
create table transactions (
   id                bigserial primary key,
   invoice           varchar(50) unique not null,
   grand_total       float,
   tax               float,
   amount            float,
   type              int,
   status            int,
   user_id           int,
   updated_at  timestamptz default now(),
   created_at  timestamptz default now(),
   foreign key (user_id) references users (id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
drop table transactions;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
