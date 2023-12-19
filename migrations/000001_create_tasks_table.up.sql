CREATE TABLE IF NOT EXISTS tasks (
id bigserial PRIMARY KEY,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
title text NOT NULL,
description text NOT NULL,
status text NOT NULL,
version integer NOT NULL DEFAULT 1
);