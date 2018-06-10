create schema if not exists ledger;

create extension if not exists "uuid-ossp";

create table ledger.services(
    service_id uuid default uuid_generate_v4(),
    name text unique,
    primary key(service_id)
);

create table ledger.resources(
    resource_id uuid,
    service_id uuid references ledger.services,
	created_at timestamp not null default now(),
	primary key(resource_id)
);

create table ledger.revisions(
    revision_id serial,
    resource_id uuid not null references ledger.resources,
    updated_at timestamp not null default now(),
    primary key(revision_id)
);
create index revisions_updated_at on ledger.revisions(updated_at);

create table ledger.deletes(
    resource_id uuid references ledger.resources,
    deleted_at timestamp not null default now(),
    primary key(resource_id)
);

create index deletes_deleted_at on ledger.deletes(deleted_at);
