drop table if exists ledger.devices cascade;

delete from ledger.deletes
where resource_id in (
    select resource_id
    from ledger.resources r,
         ledger.services s
    where s.name = 'devices' and
          s.service_id = r.service_id
);

delete from ledger.resources
where service_id in (select service_id from ledger.services where name = 'devices');

delete from ledger.services
where name = 'devices';
