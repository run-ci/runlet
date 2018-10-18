CREATE TABLE pipelines (
    remote varchar(255) NOT NULL UNIQUE,
    ref varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    
    CONSTRAINT remote_ref PRIMARY KEY(remote, ref)
);

CREATE TABLE runs (
    count SERIAL NOT NULL UNIQUE,
    start_time timestamp NOT NULL,
    end_time timestamp NOT NULL,
    success boolean NOT NULL,

    remote varchar(255) NOT NULL,
    ref varchar(255) NOT NULL,

    FOREIGN KEY (remote, ref) REFERENCES pipelines(remote, ref),
    PRIMARY KEY (remote, ref, count)
);

CREATE TABLE steps (
    id SERIAL NOT NULL UNIQUE,
    name varchar(255) NOT NULL,
    success boolean NOT NULL,

    remote varchar(255) NOT NULL,
    ref varchar(255) NOT NULL,
    count INTEGER NOT NULL,

    FOREIGN KEY (remote, ref, count) REFERENCES runs(remote, ref, count),
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
