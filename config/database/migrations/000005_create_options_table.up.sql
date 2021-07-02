CREATE TABLE IF NOT EXISTS public."options" (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	"name" text NOT NULL,
	value text NULL,
	CONSTRAINT options_pkey PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_options_deleted_at ON public.options USING btree (deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_options_name ON public.options USING btree (name);
