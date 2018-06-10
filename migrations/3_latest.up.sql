create schema if not exists latest;

create view latest.revisions as
select distinct on (p.resource_id)
    p.resource_id,
    p.revision_id,
    r.created_at,
    p.updated_at
from ledger.revisions p,
     ledger.resources r
where p.resource_id = r.resource_id AND
      r.resource_id NOT IN (
        select resource_id
        from ledger.deletes
        where deleted_at > p.updated_at
      )
order by p.resource_id, p.revision_id desc;

create view latest.devices as
select distinct on (p.resource_id)
    p.resource_id,
    p.revision_id,
    r.created_at,
    p.updated_at,
    p.name
from ledger.devices p,
     ledger.resources r
where p.resource_id = r.resource_id AND
      r.resource_id NOT IN (
        select resource_id
        from ledger.deletes
        where deleted_at > p.updated_at
      )
order by p.resource_id, p.revision_id desc;

insert into ledger.resources(resource_id, service_id)
values ('6ffbbdf4-b5e8-420d-901d-0134289ec3b0',
        (select service_id from ledger.services where name='devices'));

insert into ledger.devices(resource_id, name)
values ('6ffbbdf4-b5e8-420d-901d-0134289ec3b0', 'bob');

insert into ledger.deletes(resource_id)
values ('6ffbbdf4-b5e8-420d-901d-0134289ec3b0')
on conflict (resource_id) do update
    set deleted_at = now();

insert into ledger.devices(resource_id, name)
values ('6ffbbdf4-b5e8-420d-901d-0134289ec3b0', 'bob');
