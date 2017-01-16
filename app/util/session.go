package util

import "github.com/astaxie/beego/session"

var SessionManager *session.Manager

func RegisterSessionManager(cfg *session.ManagerConfig) error {
	// Create new session manager
	globalSessions, err := session.NewManager("memory", cfg)

	if err != nil {
		return err
	}

	// Start the session garbage collector
	go globalSessions.GC()

	SessionManager = globalSessions

	return nil
}
