create table organizations(
    id serial primary key not null,    
    name varchar(250) not null,
    created_at date default now(),
    updated_at date default now() 
);