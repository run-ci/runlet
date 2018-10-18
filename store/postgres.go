package store

import (
	"errors"

	sql "github.com/jmoiron/sqlx"

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

	logger.Debug("loading pipeline")

	rows, err := st.db.NamedQuery("SELECT * FROM pipelines WHERE pipelines.remote = :remote", p)
	if err != nil {
		logger.WithField("error", err).Debug("unable to load pipeline")

		return err
	}

	for rows.Next() {
		return rows.StructScan(p)
	}

	return errors.New("pipeline not found")
}
