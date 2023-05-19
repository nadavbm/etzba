DO $$ DECLARE
BEGIN
--
-- migrations pattern: if exist move on
--
IF EXISTS(SELECT 1 FROM pg_tables WHERE tablename = 'migrations') THEN
  RAISE NOTICE 'migrations table exists, skipping initial table creation';
  RETURN;
END IF;

--
-- Name: migrations; Type: TABLE; Schema: public;
--
CREATE TABLE migrations (
    name text PRIMARY KEY,
    time TIMESTAMP DEFAULT NOW()
);

END $$;

--
-- Name: results table; Type: TABLE; Description: results table
--
DO $$ BEGIN
IF EXISTS(SELECT 1 FROM migrations WHERE name = 'create-results-table') THEN RETURN;
END IF;

CREATE TABLE results (
	id SERIAL PRIMARY KEY NOT NULL,
    created timestamptz NOT NULL DEFAULT NOW(),
	type TEXT NOT NULL,
	job_duration DOUBLE PRECISION NOT NULL,
    avg_duration DOUBLE PRECISION NOT NULL,
	min_duration DOUBLE PRECISION NOT NULL,
    med_duration DOUBLE PRECISION NOT NULL,
    max_duration DOUBLE PRECISION NOT NULL,
    total INT NOT NULL
);

INSERT INTO migrations (name) VALUES ('create-results-table');
END $$;


--
-- Name: data; Description: add data to results table
--

DO $$ BEGIN
IF
EXISTS(SELECT 1 FROM migrations WHERE name = 'add-data') THEN RETURN;
END IF;

INSERT INTO results (type, job_duration, avg_duration, min_duration, med_duration, max_duration, total)
VALUES ('api', 65.65,12.32,56.32,31.14,99.90,10);

INSERT INTO results (type, job_duration, avg_duration, min_duration, med_duration, max_duration, total)
VALUES ('api', 45.45,11.12,49.19,32.34,90.91,21);

INSERT INTO results (type, job_duration, avg_duration, min_duration, med_duration, max_duration, total)
VALUES ('api', 51.52,13.12,54.43,32.32,87.98,109);

INSERT INTO results (type, job_duration, avg_duration, min_duration, med_duration, max_duration, total)
VALUES ('api', 54.54,14.14,56.29,35.87,86.21,230);

INSERT INTO results (type, job_duration, avg_duration, min_duration, med_duration, max_duration, total)
VALUES ('api', 56.56,11.23,51.15,38.12,84.32,509);

INSERT INTO results (type, job_duration, avg_duration, min_duration, med_duration, max_duration, total)
VALUES ('api', 45.45,14.22,54.11,36.11,81.22,15);

INSERT INTO results (type, job_duration, avg_duration, min_duration, med_duration, max_duration, total)
VALUES ('sql', 66.66,12.32,56.32,31.14,99.90,67);

INSERT INTO results (type, job_duration, avg_duration, min_duration, med_duration, max_duration, total)
VALUES ('sql', 55.55,11.12,49.19,32.34,90.91,107);

INSERT INTO results (type, job_duration, avg_duration, min_duration, med_duration, max_duration, total)
VALUES ('sql', 44.44,13.12,54.43,32.32,87.98,279);

INSERT INTO results (type, job_duration, avg_duration, min_duration, med_duration, max_duration, total)
VALUES ('sql', 34.34,14.14,56.29,35.87,86.21,78);

INSERT INTO migrations (name) VALUES ('add-data');
END $$;
