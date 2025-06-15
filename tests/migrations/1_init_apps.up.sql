insert into apps(id, name, secret)
values(1, 'test-app', 'test-secret') on conflict do nothing;    