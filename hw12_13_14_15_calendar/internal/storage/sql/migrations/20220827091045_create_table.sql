-- +goose Up
-- +goose StatementBegin
create table event
(
    id          uuid not null
        primary key,
    title       text,
    date_start  timestamp,
    date_end    timestamp,
    description text,
    owner_id    text,
    alarm_time  bigint
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists event
-- +goose StatementEnd
