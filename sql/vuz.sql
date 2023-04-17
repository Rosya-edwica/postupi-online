BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS vuz(
    id character varying(255) NOT NULL,
    name text NOT NULL,
    description text,
    city character varying(255) NOT NULL,
    cost integer,
    budget_points double precision,
    payment_points double precision,
    budget_places integer,
    payment_places integer,
    image text,
    logo text,
    url text,

    primary key(id)
);

CREATE TABLE IF NOT EXISTS vuz_specialization(
    id character varying(100) NOT NULL,
    name text NOT NULL,
    description text,
    direction text NOT NULL,
    cost integer,
    budget_points double precision,
    payment_points double precision,
    budget_places integer,
    payment_places integer,
    image text,
    url text NOT NULL,

    primary key(id)
);

CREATE TABLE IF NOT EXISTS vuz_program (
    id character varying(100) NOT NULL,
    name text NOT NULL,
    description text,
    direction character varying(100) NOT NULL,
    form character varying(100) NOT NULL,
    subjects text[] DEFAULT '{}'::text[],
    cost integer,
    has_professions boolean DEFAULT false NOT NULL,
    budget_points double precision,
    payment_points double precision,
    budget_places integer,
    payment_places integer,
    image text,
    url text NOT NULL,

    primary key(id),
    constraint check_form check(form in ('Бакалавриат', 'Специалитет', 'Магистратура'))
);

CREATE TABLE IF NOT EXISTS vuz_profession(
    id serial,
    name text not null,
    image text not null,

    primary key(id),
    constraint unique_profession unique(name)
);

CREATE TABLE IF NOT EXISTS vuz_contacts (
    id serial,
    vuz_id character varying(255) NOT NULL,
    address text NOT NULL,
    phones text NOT NULL,
    email character varying(255) NOT NULL,
    website text NOT NULL,

    primary key(id),
    foreign key(vuz_id) references vuz(id),
    constraint unique_contacts unique(vuz_id)
);

CREATE TABLE IF NOT EXISTS vuz_to_specialization(
    id serial,
    vuz_id varchar(255) not null,
    spec_id varchar(255) not null,

    primary key(id),
    foreign key(vuz_id) references vuz(id) on delete cascade,
    foreign key(spec_id) references vuz_specialization(id) on delete cascade
);

CREATE TABLE IF NOT EXISTS vuz_specialization_to_program(
    id serial,
    vuz_specialization integer,
    program_id varchar(255) not null,

    primary key(id),
    foreign key(vuz_specialization) references vuz_to_specialization(id) on delete cascade,
    foreign key(program_id) references vuz_program(id) on delete cascade
);

CREATE TABLE IF NOT EXISTS vuz_specialization_program_to_profession(
    id serial,
    vuz_specialization_program integer,
    profession_id integer,

    primary key(id),
    foreign key(vuz_specialization_program) references vuz_specialization_to_program(id) on delete cascade,
    foreign key(profession_id) references vuz_profession(id) on delete cascade
);

END TRANSACTION;