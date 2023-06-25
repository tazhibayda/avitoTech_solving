create table users
(
    user_id integer not null unique,
    balance float   not null
);


create table transactions
(
    transaction_id bigserial
            primary key,
    user_id        integer                            not null
        constraint transactions___fk
            references "users" (user_id),
    amount         float                              not null,
    operation      varchar(55)                        not null,
    date           timestamp  not null DEFAULT NOW()
);
