CREATE TABLE IF NOT EXISTS regular_business_hours (
    regular_business_hours_id serial primary key,
    dogrun_id bigint not null,
    day int not null,
    open_time time,
    close_time time,
    is_all_day boolean DEFAULT FALSE,
    is_closed boolean DEFAULT FALSE,
    created_at timestamp,
    updated_at timestamp
);