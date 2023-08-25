SET DATABASE = defaultdb;
DROP DATABASE IF EXISTS rotavator;
-- read script

CREATE DATABASE IF NOT EXISTS rotavator;
SET DATABASE = rotavator;
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.login (
    id UUID PRIMARY KEY  NOT NULL DEFAULT gen_random_uuid(),
    email_md5 TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp(),
    code TEXT NOT NULL,
    CONSTRAINT email_md5_unique UNIQUE (email_md5)
);

CREATE TABLE users (
    id UUID PRIMARY KEY  NOT NULL DEFAULT gen_random_uuid(),
    name TEXT,
    email TEXT NOT NULL,
    unavailable TEXT[],
    CONSTRAINT email_unique UNIQUE (email)
);

CREATE TABLE orgs (
    id UUID PRIMARY KEY  NOT NULL DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    admin_id UUID NOT NULL
);

CREATE TABLE org_rotas (
    id UUID PRIMARY KEY  NOT NULL DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL,
    name TEXT NOT NULL,
    roles TEXT[],
    slots TEXT[],
    CONSTRAINT fk_org_id FOREIGN KEY(org_id) REFERENCES orgs(id),
    CONSTRAINT org_name_id_unique UNIQUE (org_id, name)
);

CREATE TABLE org_users (
    org_id UUID NOT NULL,
    user_id UUID NOT NULL,
    roles TEXT [],
    CONSTRAINT fk_org_id FOREIGN KEY(org_id) REFERENCES orgs(id),
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id)
);

INSERT INTO auth.login
    (email_md5, code)
VALUES
    ('some md5 hash', 'code 123');

-- alice     '10000000-0000-4000-0000-000000000001'
-- bob       '10000000-0000-4000-0000-000000000002'
-- charley   '10000000-0000-4000-0000-000000000003'
-- ccc       '20000000-0000-4000-0000-000000000001'
-- boathouse '30000000-0000-4000-0000-000000000001'
insert into users (id, name,email, unavailable) values ('10000000-0000-4000-0000-000000000001', 'alice','alice@example.com'
                                          ,ARRAY ['wed','2023/11/01','2023/12/21-2023/12/22']);
insert into users (id, name,email) values ('10000000-0000-4000-0000-000000000002', 'bob','bob@example.com');
insert into users (id, name,email) values ('10000000-0000-4000-0000-000000000003', 'charley','charley@example.com');

insert into orgs (id, name,admin_id) values ('20000000-0000-4000-0000-000000000001', 'ccc','10000000-0000-4000-0000-000000000001');

insert into org_rotas (id, org_id, name, roles, slots) values (
    '30000000-0000-4000-0000-000000000001'
    ,'20000000-0000-4000-0000-000000000001'
    ,'Club House'
    ,ARRAY ['key_holder','admin','quartermaster']
    ,ARRAY ['wed 18:00-20:00','sat 10:00-12:00','sun 10:00-12:00']
);

insert into org_users (org_id, user_id, roles) values (
    '20000000-0000-4000-0000-000000000001',
    '10000000-0000-4000-0000-000000000001',
    ARRAY ['key_holder','quartermaster']
);
insert into org_users (org_id, user_id, roles) values (
      '20000000-0000-4000-0000-000000000001',
      '10000000-0000-4000-0000-000000000002',
      ARRAY ['key_holder']
  );

-----
-- all users for a rota
select * from users;

select u.name, o.name from org_users
    inner join users u on org_users.user_id = u.id
    inner join orgs o on org_users.org_id = o.id;

select o.name, org_rotas.name, roles, slots from org_rotas
    inner join orgs o on org_rotas.org_id = o.id;

select u.name, orot.name, ou.roles, u.unavailable
from users u
    inner join org_users ou on ou.user_id = u.id
    inner join org_rotas orot on orot.org_id = ou.org_id;


