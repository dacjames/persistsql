drop schema if exists ledger cascade;
create schema if not exists ledger;

drop schema if exists deleted cascade;
create schema if not exists deleted;

drop schema if exists latest cascade;
create schema if not exists latest;

drop schema if exists history cascade;
create schema if not exists history;

drop schema if exists bindings cascade;
create schema if not exists bindings;

drop schema if exists revision cascade;
create schema if not exists revision;

create table ledger.services(
	service_id serial,
	name text unique,
	primary key(service_id)
);
create index services_service_id on ledger.services(service_id);

create table ledger.resources(
	resource_id serial,
	service_id integer references ledger.services(service_id),
	created_at timestamp not null default now(),
	primary key(resource_id)
);

create table ledger.products(
	revision_id serial,
	updated_at timestamp not null default now(),
	resource_id integer references ledger.resources,
	product_type text,
	serial_number text,
	primary key(revision_id)
);
create index products_resource_id on ledger.products(resource_id);

create view revision.products as
select resource_id, max(revision_id) as revision_id
from ledger.products
group by resource_id;


create table ledger.devices(
	revision_id serial,
	updated_at timestamp not null default now(),
	resource_id integer references ledger.resources,
	name text,
	primary key(revision_id)
);
create index devices_resource_id on ledger.devices(resource_id);

create view revision.devices as
select resource_id, max(revision_id) as revision_id
from ledger.devices
group by resource_id;

create table deleted.products(
	resource_id integer references ledger.resources,
	deleted_at timestamp not null default now(),
	primary key(resource_id)
);
create index products_deleted_at on deleted.products(deleted_at);

create table deleted.devices(
	resource_id integer references ledger.resources,
	deleted_at timestamp not null default now(),
	primary key(resource_id)
);
create index devices_deleted_at on deleted.devices(deleted_at);

insert into ledger.services(service_id, name) values (1, 'products');
insert into ledger.services(service_id, name) values (2, 'devices');

insert into ledger.resources(resource_id, service_id) values (1, 1);
insert into ledger.resources(resource_id, service_id) values (2, 1);
insert into ledger.resources(resource_id, service_id) values (3, 1);

insert into ledger.resources(resource_id, service_id) values (4, 2);
insert into ledger.resources(resource_id, service_id) values (5, 2);
insert into ledger.resources(resource_id, service_id) values (6, 2);

insert into ledger.products(resource_id, product_type, serial_number, updated_at) values
	(1, 'fancy', '123', now());
insert into ledger.products(resource_id, product_type, serial_number, updated_at) values
	(1, 'fancy', '124', now() + interval '1 second');
insert into ledger.products(resource_id, product_type, serial_number, updated_at) values
	(2, 'cool', '125', now() + interval '2 second');
insert into ledger.products(resource_id, product_type, serial_number, updated_at) values
	(3, 'cooler', '126', now() + interval '3 second');

insert into deleted.products(resource_id, deleted_at) values
    (2, now() + interval '4 seconds');


insert into ledger.devices(resource_id, name, updated_at) values
	(4, 'bob', now());
insert into ledger.devices(resource_id, name, updated_at) values
	(4, 'bobby', now() + interval '1 second');
insert into ledger.devices(resource_id, name, updated_at) values
	(5, 'jill', now() + interval '2 second');
insert into ledger.devices(resource_id, name, updated_at) values
	(6, 'alf', now() + interval '3 second');

insert into deleted.devices(resource_id, deleted_at) values
	(5, now() + interval '4 seconds');


create view latest.products as
select distinct on (p.resource_id)
    p.resource_id,
    r.created_at,
    p.updated_at,
    p.product_type,
    p.serial_number
from ledger.products p,
        ledger.resources r
where p.resource_id = r.resource_id AND
        r.resource_id NOT IN (
            select resource_id
            from deleted.products
            where deleted_at > p.updated_at
    )
order by p.resource_id, p.revision_id desc;

create view latest.devices as
select distinct on (p.resource_id)
    p.resource_id,
    r.created_at,
    p.updated_at,
    p.name
from ledger.devices p,
        ledger.resources r
where p.resource_id = r.resource_id AND
        r.resource_id NOT IN (
            select resource_id
            from deleted.devices
            where deleted_at > p.updated_at
    )
order by p.resource_id, p.revision_id desc;

create table bindings.devices_products(
	binding_id serial,
	from_id integer not null references ledger.resources(resource_id),
	from_revision integer not null references ledger.devices(revision_id),
	to_id integer not null references ledger.resources(resource_id),
	to_revision integer not null references ledger.products(revision_id),
	bound_at timestamp not null default now(),
	primary key(binding_id)
);


insert into bindings.devices_products(from_id, from_revision, to_id, to_revision, bound_at) values
	(4, (select revision_id from revision.devices where resource_id=4), 2, (select revision_id from revision.products where resource_id=2), now() + interval '10 seconds');

insert into bindings.devices_products(from_id, from_revision, to_id, to_revision, bound_at) values
	(4, (select revision_id from revision.devices where resource_id=4), 1, (select revision_id from revision.products where resource_id=1), now() + interval '11 seconds');

insert into bindings.devices_products(from_id, from_revision, to_id, to_revision, bound_at) values
	(4, (select revision_id from revision.devices where resource_id=4), 2, (select revision_id from revision.products where resource_id=2), now() + interval '12 seconds');

insert into ledger.products(resource_id, product_type, serial_number, updated_at) values
	(1, 'fancier', '124', now() + interval '12 second');


create view latest.devices_products as
select distinct on (f.resource_id)
	t.resource_id,
	t.product_type,
	t.serial_number,
	b.bound_at
from ledger.devices f,
     ledger.products t,
	 bindings.devices_products b
where f.resource_id = b.from_id AND
	  t.resource_id = b.to_id
order by f.resource_id, b.binding_id desc
;

create view latest.products_devices as
select distinct on (f.resource_id)
	f.resource_id,
	f.name,
	b.bound_at
from ledger.devices f,
	 ledger.products t,
	 bindings.devices_products b
where f.resource_id = b.from_id AND
	  t.resource_id = b.to_id
order by f.resource_id, b.binding_id desc
;



create view history.products as
select r.resource_id,
	   'deleted' as action,
	   d.deleted_at as "at",
	   '{}'::json as resource
from ledger.resources r, deleted.products d
where r.resource_id = d.resource_id
union all
select p.resource_id,
	   'added' as action,
	   p.updated_at as "at",
       row_to_json((select d from (select p.resource_id, p.product_type, p.serial_number) d)) as resource
from ledger.products p
union all
select r.resource_id,
	   'bound from resource: devices' as action,
	   b.bound_at as "at",
	   row_to_json((select d from (select l.resource_id, l.name) d)) as resource
from ledger.resources r, bindings.devices_products b, ledger.devices l
where r.resource_id = b.to_id and
      l.resource_id = b.from_id and
	  l.revision_id = b.from_revision
order by "at" asc;

create view history.devices as
select r.resource_id,
	   'deleted' as action,
	   d.deleted_at as "at",
	   '{}'::json as resource
from ledger.resources r, deleted.devices d
where r.resource_id = d.resource_id
union all
select p.resource_id,
	   'added' as action,
	   p.updated_at as "at",
	   row_to_json((select d from (select p.resource_id, p.name) d)) as resource
from ledger.devices p
union all
select r.resource_id,
	   'bound to resource: product' as action,
	   b.bound_at as "at",
	   row_to_json((select d from (select l.resource_id, l.product_type, l.serial_number) d)) as resource
from ledger.resources r, bindings.devices_products b, ledger.products l
where r.resource_id = b.from_id and
      l.resource_id = b.to_id and
	  l.revision_id = b.to_revision
union all
select r.resource_id,
	   'bound resource updated: product' as action,
	   l.updated_at as "at",
	   row_to_json((select d from (select l.resource_id, l.product_type, l.serial_number) d)) as resource
from ledger.resources r, bindings.devices_products b, ledger.products l
where r.resource_id = b.from_id and
	  l.resource_id = b.to_id and
	  l.updated_at > b.bound_at
order by "at" asc;

create view history.resources as
select 'products' as service, h.* from history.products h
union all
select 'devices' as service, h.* from history.devices h;


select * from history.resources order by "at" asc;


