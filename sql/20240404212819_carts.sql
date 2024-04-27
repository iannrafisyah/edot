-- +goose Up
create table carts (
    id              bigserial primary key,
    name            varchar(50),
    price           float,
    qty             int,
    product_id      int,
    warehouse_id    int,
    to_warehouse_id int default null,
    user_id         int,
    updated_at      timestamptz default now(),
    created_at      timestamptz default now(),
    foreign key (product_id) references products (id),
    foreign key (user_id) references users (id),
    foreign key (warehouse_id) references warehouses (id),
    foreign key (to_warehouse_id) references warehouses (id)
);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
drop table carts;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
