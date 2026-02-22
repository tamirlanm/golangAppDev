create table if not exists users (
    id serial primary key,
    name varchar(255) not null,
    email varchar(255) not null unique,
    age int not null,
    is_employed boolean not null default false
);

insert into users (name, email, age, is_employed) values ('John Due', 'john.due@example.com', 30, true);
