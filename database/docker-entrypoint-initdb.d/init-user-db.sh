#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-SQL

CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS \$\$

    DECLARE
        data json;
        notification json;

    BEGIN

        -- Convert the old or new row to JSON, based on the kind of action.
        -- Action = DELETE?             -> OLD row
        -- Action = INSERT or UPDATE?   -> NEW row
        IF (TG_OP = 'DELETE') THEN
            data = row_to_json(OLD);
        ELSE
            data = row_to_json(NEW);
        END IF;

        -- Contruct the notification as a JSON string.
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'action', TG_OP,
                          'data', data);


        -- Execute pg_notify(channel, notification)
        PERFORM pg_notify('events',notification::text);

        -- Result is ignored since this is an AFTER trigger
        RETURN NULL;
    END;

\$\$ LANGUAGE plpgsql;



CREATE TABLE requests(
 id             serial NOT NULL PRIMARY KEY,
 method         VARCHAR (10) NOT NULL,
 url            VARCHAR NOT NULL,
 headers        JSONB NOT NULL,
 body           TEXT NOT NULL,
 content_length BIGINT NOT NULL,
 host           VARCHAR NOT NULL,

 requested_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER requests_notify_event
AFTER INSERT OR UPDATE ON requests
    FOR EACH ROW EXECUTE PROCEDURE notify_event();


CREATE TABLE responses
(
  id serial NOT NULL PRIMARY KEY,
  request_id INTEGER NOT NULL references requests(id),

  status        VARCHAR(50) NOT NULL,
  status_code    SMALLINT NOT NULL,
  headers       JSONB NOT NULL,
  body          TEXT NOT NULL,
  content_length BIGINT NOT NULL,
  mime_type     VARCHAR (50) NULL,

  responded_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER responses_notify_event
AFTER INSERT OR UPDATE ON responses
    FOR EACH ROW EXECUTE PROCEDURE notify_event();

SQL