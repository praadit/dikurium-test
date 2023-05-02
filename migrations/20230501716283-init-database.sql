
-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.users
(
    id uuid NOT NULL default uuid_generate_v4(),
    email text NOT NULL,
    password text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE TABLE public.todos
(
    id uuid NOT NULL default uuid_generate_v4(),
    user_id uuid NOT NULL,
    title text NOT NULL,
    is_completed boolean NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT todos_pkey PRIMARY KEY (id),
    CONSTRAINT user_todos_key FOREIGN KEY (user_id)
        REFERENCES public.users (id)
        ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- +migrate Down
DROP TABLE "public"."todos";
DROP TABLE "public"."users";