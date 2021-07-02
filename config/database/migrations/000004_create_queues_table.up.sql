CREATE SEQUENCE IF NOT EXISTS queues_id_seq;

CREATE TABLE IF NOT EXISTS "public"."queues" (
    "id" int8 NOT NULL DEFAULT nextval('queues_id_seq'::regclass),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "type" text,
    "object" int8 NOT NULL,
    PRIMARY KEY ("id")
);