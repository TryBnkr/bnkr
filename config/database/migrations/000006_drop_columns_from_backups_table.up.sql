ALTER TABLE backups 
DROP COLUMN IF EXISTS database_name,
DROP COLUMN IF EXISTS database_user,
DROP COLUMN IF EXISTS database_password,
DROP COLUMN IF EXISTS database_host,
DROP COLUMN IF EXISTS database_port,
DROP COLUMN IF EXISTS deployment;