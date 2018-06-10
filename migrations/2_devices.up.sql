create table if not exists ledger.devices (
    name text not null
) inherits (ledger.revisions);

insert into ledger.services (name)
values ('devices')
on conflict do nothing;
