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

func (st *Postgres) ReadPipeline(p *Pipeline) error {
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

func (st *Postgres) CreateRun(r *Run) error {
	logger := logger.WithFields(log.Fields{
		"pipeline_remote": r.PipelineRemote,
		"pipeline_name":   r.PipelineName,
	})

	sqlinsert := `
	INSERT INTO runs (start_time, end_time, success, pipeline_remote, pipeline_name)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING count
	`

	logger.Debug("saving pipeline run")

	// Using QueryRow because the insert is returning "count".
	err := st.db.QueryRow(
		sqlinsert, r.Start, r.End, r.Success, r.PipelineRemote, r.PipelineName).
		Scan(&r.Count)
	if err != nil {
		logger.WithField("error", err).Debug("unable to insert pipeline run")
		return err
	}

	logger.Debug("pipeline run saved")

	return nil
}

func (st *Postgres) CreateStep(s *Step) error {
	logger := logger.WithFields(log.Fields{
		"pipeline_remote": s.PipelineRemote,
		"pipeline_name":   s.PipelineName,
		"run_count":       s.RunCount,
		"name":            s.Name,
	})

	sqlinsert := `
	INSERT INTO steps (name, start_time, end_time, success, pipeline_remote, pipeline_name, run_count)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`

	logger.Debug("saving run step")

	// Using QueryRow because the insert is returning "id".
	err := st.db.QueryRow(
		sqlinsert, s.Name, s.Start, s.End, s.Success, s.PipelineRemote, s.PipelineName, s.RunCount).
		Scan(&s.ID)
	if err != nil {
		logger.WithField("error", err).Debug("unable to insert run step")
		return err
	}

	logger.Debug("run step saved")

	return nil
}

func (st *Postgres) CreateTask(t *Task) error {
	logger := logger.WithFields(log.Fields{
		"name":    t.Name,
		"step_id": t.StepID,
	})

	sqlinsert := `
	INSERT INTO tasks (name, start_time, end_time, success, step_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	logger.Debug("saving step task")

	// Using QueryRow because the insert is returning "id".
	err := st.db.QueryRow(
		sqlinsert, t.Name, t.Start, t.End, t.Success, t.StepID).
		Scan(&t.ID)
	if err != nil {
		logger.WithField("error", err).Debug("unable to insert step task")
		return err
	}

	logger.Debug("step task saved")

	return nil
}

func (st *Postgres) UpdateRun(r *Run) error {
	logger := logger.WithFields(log.Fields{
		"pipeline_remote": r.PipelineRemote,
		"pipeline_name":   r.PipelineName,
		"count":           r.Count,
		"end":             r.End,
		"success":         r.Success,
	})

	sqlupdate := `
	UPDATE runs
	SET success = $1, end_time = $2
	WHERE runs.pipeline_remote = $3 AND runs.pipeline_name = $4 AND runs.count = $5
	`

	logger.Debug("saving run step")

	st.db.Exec(sqlupdate, r.Success, r.End, r.PipelineRemote, r.PipelineName, r.Count)

	logger.Debug("run step saved")

	return nil
}

func (st *Postgres) UpdateStep(s *Step) error {
	logger := logger.WithFields(log.Fields{
		"pipeline_remote": s.PipelineRemote,
		"pipeline_name":   s.PipelineName,
		"run_count":       s.RunCount,
		"name":            s.Name,
		"id":              s.ID,
		"success":         s.Success,
		"end":             s.End,
	})

	sqlupdate := `
	UPDATE steps
	SET success = $1, end_time = $2
	WHERE steps.id = $3
	`

	logger.Debug("saving run step")

	st.db.Exec(sqlupdate, s.Success, s.End, s.ID)

	logger.Debug("run step saved")

	return nil
}

func (st *Postgres) UpdateTask(t *Task) error {
	logger := logger.WithFields(log.Fields{
		"name":    t.Name,
		"step_id": t.StepID,
		"success": t.Success,
		"id":      t.ID,
		"end":     t.End,
	})

	sqlupdate := `
	UPDATE tasks
	SET success = $1, end_time = $2
	WHERE tasks.id = $3
	`

	logger.Debug("saving step task")

	st.db.Exec(sqlupdate, t.Success, t.End, t.ID)

	logger.Debug("step task saved")

	return nil
}
