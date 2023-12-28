CREATE TABLE "servicemesh"."service" (
  "id" serial,
  "external_id" uuid NOT NULL,
  "name" varchar(255) NOT NULL,
  "version" varchar(255) NOT NULL,
  "port" int4 NOT NULL,
  "active" bool NOT NULL DEFAULT true,
  "started" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "stopped" timestamp DEFAULT NULL,
  PRIMARY KEY ("id")
)
;

ALTER TABLE "servicemesh"."service" 
  OWNER TO "servicemesh";

CREATE UNIQUE INDEX "service_idx1" ON "servicemesh"."service" USING btree (
  "port"
);