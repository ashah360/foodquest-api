create extension if not exists "uuid-ossp";

create table if not exists users (
                                     id uuid default uuid_generate_v4(),
                                     email text not null,
                                     password text not null,

                                     first_name text not null,
                                     last_name text not null,
                                     phone_number text not null,

                                     address_line_1 text,
                                     address_line_2 text,
                                     address_line_3 text,
                                     city text,
                                     state text,
                                     postal_code text,
                                     country text,

                                     created_at timestamp without time zone default now(),
                                     last_login timestamp without time zone,

                                     primary key (id)
)