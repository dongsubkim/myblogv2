drop table posts;
drop table comments;

create table posts (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  category   varchar(255)[],
  content    text,
  created_at timestamp not null  
);

create table comments (
  id         serial primary key,
  username   varchar(255),
  password   text,
  body       text,
  post_id    integer references posts(id),
  created_at timestamp not null       
);
