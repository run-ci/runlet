CREATE TABLE pipelines (
    remote varchar(255) NOT NULL UNIQUE,
    branch varchar(255) NOT NULL,
    tag varchar(255),
    name varchar(255) NOT NULL,
    
    CONSTRAINT remote_branch PRIMARY KEY(remote, branch)
);
