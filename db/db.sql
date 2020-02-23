drop table if exists request;
create table request(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    method text,
    proto_schema text,
    dest text,
    body text,
    header text,
    created text
);