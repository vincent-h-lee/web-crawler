create table crawls(
  id bigserial primary key, 
  url varchar(2048) not null,
  title varchar(256),
  status_code smallint not null,
  timestamp timestamp not null
);

create table crawl_headings(
  id bigserial primary key,
  crawl_id bigint not null,

  text varchar(256),
  tag char(2),

  CONSTRAINT fk_crawl FOREIGN KEY(crawl_id) REFERENCES crawls(id)
);

create table crawl_links(
  id bigserial primary key,
  crawl_id bigint not null,

  text varchar(256),
  url varchar(2048) not null,

  CONSTRAINT fk_crawl FOREIGN KEY(crawl_id) REFERENCES crawls(id)
);