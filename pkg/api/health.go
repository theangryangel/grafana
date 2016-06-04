package api

import (
	"github.com/grafana/grafana/pkg/middleware"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/setting"

	"time"

	"gopkg.in/macaron.v1"
)

// GET /api/health
func GetHealth(c *middleware.Context) {
	var err error
	status := 200
	res := map[string]interface{}{}
	res["current_timestamp"] = time.Now().Format(time.RFC3339)
	res["database_ok"] = true
	res["session_ok"] = true

	res["version"] = map[string]interface{}{
		"version": setting.BuildVersion,
		"commit": setting.BuildCommit,
		"built": setting.BuildStamp,
	}

	// Check database
	statsQuery := m.GetSystemStatsQuery{}
	err = bus.Dispatch(&statsQuery);

	if err != nil {
		status = 500
		res["database_ok"] = false
	}

	// Check session
	// TODO: This is even sane?
	sessionWrapper := middleware.GetSession()
	err = sessionWrapper.Start(c)
	if err != nil {
		status = 500
		res["session_ok"] = false
	}
	sessionWrapper.Destory(c)

	c.JSON(status, res)
}
