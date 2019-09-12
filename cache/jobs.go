package cache

import (
	"database/sql"
	"sync"

	"github.com/Hao1995/104hackathon/glob"
	"github.com/astaxie/beego/logs"
)

var (
	jobsIns *jobs
	once    sync.Once
)

func GetJobsInstance() *jobs {
	once.Do(func() {
		jobsIns = &jobs{mu: sync.Mutex{}}
	})
	return jobsIns
}

type jobs struct {
	mu        sync.Mutex
	countJobs sql.NullInt64
}

func (c *jobs) CountJobs() (int64, error) {
	if !c.countJobs.Valid {
		c.mu.Lock()
		defer c.mu.Unlock()
		if !c.countJobs.Valid {
			// Query data
			if err := glob.DB.QueryRow("SELECT COUNT(1) FROM `jobs`").Scan(&c.countJobs); err != nil {
				return c.countJobs.Int64, err
			}
			logs.Debug("All job = %v", c.countJobs.Int64)
		}
	}
	return c.countJobs.Int64, nil
}

func (c *jobs) Refresh() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Clear data
	c.countJobs = sql.NullInt64{}

	return
}
