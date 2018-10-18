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
    success boolean NOT NULL,

    run_remote varchar(255) NOT NULL,
    run_name varchar(255) NOT NULL,
    run_count INTEGER NOT NULL,

    FOREIGN KEY (run_remote, run_name, run_count) REFERENCES runs(pipeline_remote, pipeline_name, count),
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
