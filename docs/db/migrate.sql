-- このファイルはmigrationが安定するまでのメモ
-- migrationが安定したら、すべてmigrationで管理する
\c dog_runner;

create schema dev;

create role developer login SUPERUSER password 'developer';

set search_path = 'dev';

CREATE TABLE IF NOT EXISTS dog_owners (
    dog_owner_id serial primary key,
    name varchar(128) not null,
    email varchar(255) not null,
    image text,
    sex char(1),
    reg_at timestamp not null,
    upd_at timestamp not null
);

CREATE TABLE IF NOT EXISTS dogs (
    dog_id serial primary key,
    dog_owner_id bigint not null,
    name varchar(128) not null,
    dog_type_id int,
    weight int,
    sex char(1),
    image text,
    reg_at timestamp not null,
    upd_at timestamp not null
);

CREATE TABLE IF NOT EXISTS dog_type_mst (
    dog_type_id serial primary key,
    name varchar(64) not null
);

CREATE TABLE IF NOT EXISTS injection_certifications (
    injection_certification_id serial primary key,
    dog_id bigint not null,
    type int not null,
    file text not null,
    reg_at timestamp not null,
    upd_at timestamp not null
);

CREATE TABLE IF NOT EXISTS dogrun_managers (
    dogrun_manager_id serial primary key,
    name varchar(128),
    email varchar(255) not null,
    reg_at timestamp not null,
    upd_at timestamp not null
);

CREATE TABLE IF NOT EXISTS dogruns (
    dogrun_id serial primary key,
    dogrun_manager_id bigint,
    name varchar(256) not null,
    address varchar(256),
    postcode varchar(8),
    business_day int,
    holiday int,
    open_time time,
    close_time time,
    description text,
    reg_at timestamp not null,
    upd_at timestamp not null
);

CREATE TABLE IF NOT EXISTS dogrun_images (
    dogrun_image_id serial primary key,
    dogrun_id bigint not null,
    image text not null,
    sort_order int,
    upload_at timestamp
);

CREATE TABLE IF NOT EXISTS dogrun_tags (
    dogrun_tag_id serial primary key,
    dogrun_id bigint not null,
    tag_id int not null
);

CREATE TABLE IF NOT EXISTS tag_mst (
    tag_id serial primary key,
    tag_name varchar(64) not null,
    description text
);


CREATE TABLE IF NOT EXISTS auth_dog_owners (
    auth_dog_owner_id serial primary key,
    dog_owner_id bigint not null,
    password varchar(256),
    grant_type int not null,
    login_at timestamp
);


CREATE TABLE IF NOT EXISTS auth_dogrun_managers (
    auth_dogrun_manager_id serial primary key,
    dogrun_manager_id bigint not null,
    password varchar(256),
    grant_type int not null,
    login_at timestamp
);


-- add foreign key 

alter table dogs add foreign key (dog_owner_id) references dog_owners (dog_owner_id);
alter table dogs add foreign key (dog_type_id) references dog_type_mst (dog_type_id);

alter table injection_certifications add foreign key (dog_id) references dogs (dog_id);

alter table dogruns add foreign key (dogrun_manager_id) references dogrun_managers (dogrun_manager_id);

alter table dogrun_images add foreign key (dogrun_id) references dogruns (dogrun_id);

alter table dogrun_tags add foreign key (dogrun_id) references dogruns (dogrun_id);
alter table dogrun_tags add foreign key (tag_id) references tag_mst (tag_id);

alter table auth_dog_owners add foreign key (dog_owner_id) references dog_owners (dog_owner_id);

alter table auth_dogrun_managers add foreign key (dogrun_manager_id) references dogrun_managers (dogrun_manager_id);

-- drop foreign key

alter table dogs drop constraint dev_dogs_dog_owner_id_fkey;
alter table dogs drop constraint dev_dogs_dog_type_id_fkey;

alter table injection_certifications drop constraint dev_injection_certifications_dog_id_fkey;

alter table dogruns drop constraint dev_dogruns_dogrun_manager_id_fkey;

alter table dogrun_images drop constraint dev_dogrun_images_dogrun_id_fkey;

alter table dogrun_tags drop constraint dev_dogrun_tags_dogrun_id_fkey;
alter table dogrun_tags drop constraint dev_dogrun_tags_tag_id_fkey;

alter table auth_dog_owners drop constraint dev_auth_dog_owners_dog_owner_id_fkey;

alter table auth_dogrun_managers drop constraint dev_auth_dogrun_managers_dogrun_manager_id_fkey;

