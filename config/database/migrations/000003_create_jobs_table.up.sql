CREATE TABLE IF NOT EXISTS public.jobs (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	file text NULL,
	status text NULL,
	backup int8 NOT NULL,
	CONSTRAINT jobs_pkey PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_jobs_deleted_at ON public.jobs USING btree (deleted_at);
