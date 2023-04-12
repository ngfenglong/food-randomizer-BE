-- Table: public.place_category

-- DROP TABLE public.place_category;

CREATE TABLE public.place_category
(
    id integer NOT NULL,
    place_id integer,
    category_id integer,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT place_category_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.place_category
    OWNER to root;
	
	
-- Table: public.place

-- DROP TABLE public.place;

CREATE TABLE public.place
(
    id integer NOT NULL DEFAULT nextval('place_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default",
    description text COLLATE pg_catalog."default",
    is_halal boolean,
    is_vegetarian boolean,
    location character varying COLLATE pg_catalog."default",
    lat character varying COLLATE pg_catalog."default",
    "long" character varying COLLATE pg_catalog."default",
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    category character varying COLLATE pg_catalog."default",
    CONSTRAINT "Place_pkey" PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.place
    OWNER to root;
	
	
-- Table: public.category

-- DROP TABLE public.category;

CREATE TABLE public.category
(
    id integer NOT NULL,
    category_name character varying COLLATE pg_catalog."default",
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT category_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.category
    OWNER to root;
	
	
-- SEQUENCE: public.place_id_seq

-- DROP SEQUENCE public.place_id_seq;

CREATE SEQUENCE public.place_id_seq
    INCREMENT 1
    START 4
    MINVALUE 1
    MAXVALUE 100000
    CACHE 1;

ALTER SEQUENCE public.place_id_seq
    OWNER TO root;