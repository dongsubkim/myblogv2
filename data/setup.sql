drop table admins cascade;
drop table sessions cascade;
drop table posts cascade;
drop table comments cascade;
drop table images cascade;

create table admins (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  name       varchar(255),
  email      varchar(255) not null unique,
  password   varchar(255) not null,
  created_at timestamp not null   
);

create table sessions (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  email      varchar(255),
  admin_id   integer references admins(id),
  created_at timestamp not null   
);

create table posts (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  title      varchar(255),
  category   varchar(255)[],
  content    text,
  created_at timestamp not null  
);

create table comments (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  username   varchar(255) not null,
  password   varchar(255) not null,
  body       text not null,
  post_uuid  varchar(64) references posts(uuid),
  created_at timestamp not null       
);

create table images (
  id          serial primary key,
  uuid        varchar(255) not null unique,
  filename    varchar(255) not null,
  file_url    text not null,
  post_uuid   varchar(64) references posts(uuid),
  created_at  timestamp not null
)
