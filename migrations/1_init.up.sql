create table if not exists users (
    id integer primary key,
    email text not null unique,
    pass_hash blob not null,
    is_admin boolean not null default false
);
create index if not exists idx_email on users(email);
create table if not exists apps (
    id integer primary key,
    name text not null unique,
    secret text not null unique
);