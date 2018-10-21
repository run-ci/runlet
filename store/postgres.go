package store

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq" // load the postgres driver
	log "github.com/sirupsen/logrus"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(connstr string) (PipelineStore, error) {
	logger = logger.WithField("store", "postgres")

	logger.Debug("connecting to database")

	db, err := sql.Open("postgres", connstr)
	if err != nil {
		logger.WithField("error", err).Debug("unable to connect to database")
		return nil, err
	}

	return &Postgres{
		db: db,
	}, nil
}

func (st *Postgres) SavePipeline(p Pipeline) error {
	logger := logger.WithFields(log.Fields{
		"pipeline": p,
	})

	logger.Debug("saving pipeline")

	return nil
}

func (st *Postgres) LoadPipeline(p *Pipeline) error {
	logger := logger.WithFields(log.Fields{
		"pipeline": p,
	})

	q := `SELECT remote, ref, name FROM pipelines
	WHERE pipelines.remote = $1 AND pipelines.name = $2`
	logger = logger.WithField("query", q)

	logger.Debug("loading pipeline")

	row := st.db.QueryRow(q, p.Remote, p.Name)

	logger.Debug("scanning rows")

	err := row.Scan(&p.Remote, &p.Ref, &p.Name)
	if err == sql.ErrNoRows {
		return errors.New("pipeline not found")
	}

	return err
}
