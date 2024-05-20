create table if not exists drug_type
(
    id       serial primary key,
    name_    text,
    produced boolean
);


create table if not exists drug
(
    id             serial primary key,
    name_          text,
    type_id        integer not null
        references drug_type (id)
            ON DELETE CASCADE,
    critical_count integer not null
        constraint critical_count_not_negative
            check (critical_count >= 0),
    price          integer not null
        constraint price_positive
            check (price > 0)
);


create table if not exists drug_uses
(
    id       serial primary key,
    use_text text
);

create table if not exists doctor
(
    id    serial primary key,
    name_ text
);

create table if not exists consumers
(
    id         serial primary key,
    name_      text not null,
    phone      text not null,
    addres     text not null,
    birth_date date not null
);

create table if not exists technology_book
(
    id          serial primary key,
    description text not null
);

create table if not exists type_uses
(
    type_id integer
        references drug_type (id)
            ON DELETE CASCADE,
    use_id  integer
        references drug_uses (id)
            ON DELETE CASCADE
);


create table if not exists recipe
(
    id          serial primary key,
    consumer_id integer not null
        references consumers (id)
            ON DELETE CASCADE,
    doctor_id   integer not null
        references doctor (id)
            ON DELETE CASCADE,
    drug_id     integer not null
        references drug (id)
            ON DELETE CASCADE,
    use_id      integer not null
        references drug_uses (id)
            ON DELETE CASCADE,
    disease     text    not null,
    drug_count  integer not null
);

create table if not exists inventasrization
(
    drug_id integer             not null
        references drug (id)
            ON DELETE CASCADE,
    date_   timestamp with time zone not null,
    count_  integer             not null
        constraint count_positive
            check (count_ >= 0),
    constraint inventasrization_pk
        primary key (drug_id, date_)
);

create table if not exists storage
(
    comp_id integer not null
        constraint storage_pk
            primary key
        constraint storage_components_id_fk
            references drug (id)
            ON DELETE CASCADE,
    count_  integer not null
        constraint count_not_negative
            check (count_ >= 0)
);

create table if not exists technology_drug
(
    tech_id integer not null
        constraint technology_drug_technology_book_id_fk
            references technology_book (id)
            ON DELETE CASCADE,
    drug_id integer not null
        constraint technology_drug_drug_id_fk
            references drug (id)
            ON DELETE CASCADE,
    constraint technology_drug_pk
        primary key (drug_id)
);

create table if not exists technology_components
(
    comp_id integer not null
        constraint technology_components_components_id_fk
            references drug (id)
            ON DELETE CASCADE,
    tech_id integer not null
        constraint technology_components_technology_book_id_fk
            references technology_book (id)
            ON DELETE CASCADE,
    count_  integer not null
        constraint count_not_negative
            check (count_ >= 0),
    constraint technology_components_pk
        primary key (comp_id, tech_id)
);

create table if not exists consumer_order
(
    recipe_id     integer   not null
        constraint consumer_order_recipe_id_fk
            references recipe (id)
            ON DELETE CASCADE,
    status        text      not null,
    complete_time timestamp not null,
    order_date    timestamp,
    id            serial primary key
);

CREATE OR REPLACE FUNCTION insert_in_storage() RETURNS TRIGGER as
$$
BEGIN
    INSERT INTO storage (comp_id, count_) VALUES (NEW."id", 0);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER insert_in_storage
    AFTER INSERT
    ON drug
    FOR EACH ROW
EXECUTE FUNCTION insert_in_storage();

CREATE OR REPLACE FUNCTION check_tc() RETURNS TRIGGER as
$$
BEGIN
    IF NEW.comp_id in (select d.id
                       from drug d
                                inner join drug_type dt on dt.id = d.type_id
                       where produced) THEN
        RAISE EXCEPTION 'comp id cannot be produced';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER emp_stamp
    BEFORE INSERT OR UPDATE
    ON technology_components
    FOR EACH ROW
EXECUTE FUNCTION check_tc();


