create table metric_type (
    id  bigserial unique primary key,
    name varchar unique not null
);
create table metrics
(
    id  bigserial not null unique ,
    name  varchar,
    type  int8 REFERENCES metric_type(id),
    delta bigint,
    value double precision
);

insert into metric_type (name) values ('gauge');
insert into metric_type (name) values ('counter');
