CREATE SCHEMA IF NOT EXISTS checklist AUTHORIZATION postgres;

DO
$do$
BEGIN
   IF EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE  rolname = 'checklist') THEN

      RAISE NOTICE 'Role "checklist" already exists. Skipping.';
   ELSE
      BEGIN
         CREATE ROLE checklist LOGIN PASSWORD 'checklist_pwd';
      EXCEPTION
         WHEN duplicate_object THEN
            RAISE NOTICE 'Role "checklist" was just created by a concurrent transaction. Skipping.';
      END;
   END IF;
END
$do$;

DO
$do$
BEGIN
   IF EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE  rolname = 'checklist_ro') THEN

      RAISE NOTICE 'Role "checklist_ro" already exists. Skipping.';
   ELSE
      BEGIN
         CREATE ROLE checklist_ro LOGIN PASSWORD 'checklist_pwd';
      EXCEPTION
         WHEN duplicate_object THEN
            RAISE NOTICE 'Role "checklist_ro" was just created by a concurrent transaction. Skipping.';
      END;
   END IF;
END
$do$;

GRANT CONNECT ON DATABASE refarch TO checklist;
GRANT CONNECT ON DATABASE refarch TO checklist_ro;
GRANT USAGE ON SCHEMA checklist TO checklist;
GRANT USAGE ON SCHEMA checklist TO checklist_ro;

GRANT ALL ON ALL TABLES IN SCHEMA checklist TO checklist;
GRANT SELECT ON ALL TABLES IN SCHEMA checklist TO checklist_ro;

