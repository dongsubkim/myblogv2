drop table posts cascade;
drop table comments cascade;

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
  username   varchar(255),
  password   text,
  body       text,
  post_uuid  varchar(64) references posts(uuid),
  created_at timestamp not null       
);
