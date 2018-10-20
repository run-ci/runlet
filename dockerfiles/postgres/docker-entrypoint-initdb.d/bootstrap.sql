CREATE TABLE pipelines (
    remote varchar(255) NOT NULL UNIQUE,
    ref varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    
    PRIMARY KEY(remote, name),
    UNIQUE(remote, name)
);

CREATE TABLE runs (
    count SERIAL NOT NULL UNIQUE,
    start_time timestamp NOT NULL,
    end_time timestamp NOT NULL,
    success boolean NOT NULL,

    pipeline_remote varchar(255) NOT NULL,
    pipeline_name varchar(255) NOT NULL,

    FOREIGN KEY (pipeline_remote, pipeline_name) REFERENCES pipelines(remote, name),
    PRIMARY KEY (pipeline_remote, pipeline_name, count)
);

CREATE TABLE steps (
    id SERIAL NOT NULL UNIQUE,
    name varchar(255) NOT NULL,
    start_time timestamp NOT NULL,
    end_time timestamp NOT NULL,
    success boolean NOT NULL,

    pipeline_remote varchar(255) NOT NULL,
    pipeline_name varchar(255) NOT NULL,
    run_count INTEGER NOT NULL,

    FOREIGN KEY (pipeline_remote, pipeline_name, run_count) REFERENCES runs(pipeline_remote, pipeline_name, count),
    PRIMARY KEY (id)
);

CREATE TABLE tasks (
    id SERIAL NOT NULL UNIQUE,
    name varchar(255) NOT NULL,
    start_time timestamp NOT NULL,
    end_time timestamp NOT NULL,
    success boolean NOT NULL,

    step_id INTEGER NOT NULL,

    FOREIGN KEY (step_id) REFERENCES steps(id),
    PRIMARY KEY (id)
);

INSERT INTO pipelines (remote, ref, name)
VALUES
    ('https://gitlab.com/run-ci/runlet', 'master', 'default');

-- TODO: clean this up, or don't even put it in here. It's just staying in here for now so I have somewhere to put it
-- and I can get some sleep.

-- Inserting run with its steps. This can be expanded for the tasks as well.
WITH run_data AS (
    INSERT INTO runs (start_time, end_time, success, pipeline_remote, pipeline_name)
    VALUES
        (current_timestamp, current_timestamp, true, 'https://gitlab.com/run-ci/runlet', 'default')
    RETURNING *
), step_data AS (
    INSERT INTO steps (name, start_time, end_time, success, pipeline_remote, pipeline_name, run_count)
        SELECT 'write', current_timestamp, current_timestamp, true, 'https://gitlab.com/run-ci/runlet', 'default', run_data.count
        FROM run_data
    RETURNING *
)
INSERT INTO tasks (name, start_time, end_time, success, step_id)
    SELECT 'write', current_timestamp, current_timestamp, true, step_data.id
    FROM step_data;
