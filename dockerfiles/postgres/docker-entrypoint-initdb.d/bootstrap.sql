CREATE TABLE pipelines (
    remote varchar(255) NOT NULL UNIQUE,
    ref varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    
    PRIMARY KEY(remote, name),
    UNIQUE(remote, name)
);

CREATE TABLE runs (
    -- TODO: fix this. This count won't be accurate. This should be an ID and another
    -- way needs to be found to determine the count.
    count SERIAL NOT NULL UNIQUE,
    start_time timestamp NOT NULL,
    end_time timestamp,
    success boolean,

    pipeline_remote varchar(255) NOT NULL,
    pipeline_name varchar(255) NOT NULL,

    FOREIGN KEY (pipeline_remote, pipeline_name) REFERENCES pipelines(remote, name),
    PRIMARY KEY (pipeline_remote, pipeline_name, count)
);

CREATE TABLE steps (
    id SERIAL NOT NULL UNIQUE,
    name varchar(255) NOT NULL,
    start_time timestamp NOT NULL,
    end_time timestamp,
    success boolean,

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
    end_time timestamp,
    success boolean,

    step_id INTEGER NOT NULL,

    FOREIGN KEY (step_id) REFERENCES steps(id),
    PRIMARY KEY (id)
);

INSERT INTO pipelines (remote, ref, name)
VALUES
    ('https://github.com/run-ci/runlet', 'master', 'default');
