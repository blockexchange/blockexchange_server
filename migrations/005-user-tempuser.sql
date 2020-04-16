

-- create temporary user with default password "temp"
insert into public.user(created, name, hash)
  values(
    extract(epoch from now()) * 1000,
    'temp',
    '$2a$10$g.6pRR93BwXfsMnPLWIKgOfIBDOcc48wJPCDtfNfzJbD/7zE2xgtm'
  );
