create table if not exists public.user (
    id serial not null primary key,
    name varchar(30) not null,
    password_hash varchar (255) not null,
    enable boolean not null,
    role varchar(30) not null,
    date_create timestamp not null
);
CREATE UNIQUE INDEX user_name_uindex ON public.user(name);