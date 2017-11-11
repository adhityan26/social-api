create schema `social-api`;

use `social-api`;

create table users
(
	id int not null auto_increment
		primary key,
	name varchar(100) null,
	email varchar(100) not null,
	updated_at int null,
	created_at int null,
	constraint users_email_uindex
		unique (email)
);