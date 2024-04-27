-- +goose Up
create table stocks (
   id             bigserial primary key,
   qty            int,
   latest         int,
   previous       int,
   operator       int,
   warehouse_id   int,
   product_id     int,
   transaction_id int default null,
   updated_at     timestamptz default now(),
   created_at     timestamptz default now(),
   foreign key (product_id) references products (id),
   foreign key (warehouse_id) references warehouses (id),
   foreign key (transaction_id) references transactions (id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
drop table stocks;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
