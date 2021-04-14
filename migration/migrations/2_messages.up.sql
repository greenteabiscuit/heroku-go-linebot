create table if not exists messages (
    id integer auto_increment primary key,
    name varchar(40),
    user_id integer,
    message varchar(40),
)