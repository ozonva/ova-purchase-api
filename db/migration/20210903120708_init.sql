-- +goose Up
-- +goose StatementBegin
create table purchases
(
    id         bigserial primary key,
    user_id    bigint,
    total      float8,
    created_at timestamp default now() not null,
    updated_at timestamp               not null,
    status     varchar(255)            not null
);
-- +goose StatementEnd
-- +goose StatementBegin
create table purchase_items
(
    id          bigserial primary key,
    purchase_id bigint not null,
    name        varchar(255),
    price       float8,
    quantity    integer
);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE purchase_items
    ADD CONSTRAINT fk_purchase_items_purchase_id FOREIGN KEY (purchase_id) REFERENCES purchases (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table purchase_items;
drop table purchases;
-- +goose StatementEnd
