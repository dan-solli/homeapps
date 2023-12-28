CREATE SCHEMA IF NOT EXISTS servicemesh AUTHORIZATION postgres;

DO
$do$
BEGIN
   IF EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE  rolname = 'servicemesh') THEN

      RAISE NOTICE 'Role "servicemesh" already exists. Skipping.';
   ELSE
      BEGIN
         CREATE ROLE servicemesh LOGIN PASSWORD 'servicemesh_pwd';
      EXCEPTION
         WHEN duplicate_object THEN
            RAISE NOTICE 'Role "servicemesh" was just created by a concurrent transaction. Skipping.';
      END;
   END IF;
END
$do$;

DO
$do$
BEGIN
   IF EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE  rolname = 'servicemesh_ro') THEN

      RAISE NOTICE 'Role "servicemesh_ro" already exists. Skipping.';
   ELSE
      BEGIN
         CREATE ROLE servicemesh_ro LOGIN PASSWORD 'servicemesh_pwd';
      EXCEPTION
         WHEN duplicate_object THEN
            RAISE NOTICE 'Role "servicemesh_ro" was just created by a concurrent transaction. Skipping.';
      END;
   END IF;
END
$do$;

GRANT CONNECT ON DATABASE refarch TO servicemesh;
GRANT CONNECT ON DATABASE refarch TO servicemesh_ro;
GRANT USAGE ON SCHEMA servicemesh TO servicemesh;
GRANT USAGE ON SCHEMA servicemesh TO servicemesh_ro;

GRANT ALL ON ALL TABLES IN SCHEMA servicemesh TO servicemesh;
GRANT SELECT ON ALL TABLES IN SCHEMA servicemesh TO servicemesh_ro;

