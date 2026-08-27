package console

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"wtc/internal/audit"
	"wtc/internal/blade"
	"wtc/internal/brake"
	"wtc/internal/cable"
	"wtc/internal/conv"
	"wtc/internal/grid"
	"wtc/internal/ns"
	"wtc/internal/pitch"
	"wtc/internal/store"
	"wtc/internal/tower"
	"wtc/internal/wind"
	"wtc/internal/yaw"
)

type Server struct {
	addr   string
	router *chi.Mux
	st     *store.Store
	limits ns.Limits
	audit  *audit.Recorder
	blade  *blade.Blade
	brake  *brake.System
	pitch  *pitch.System
	yaw    *yaw.System
	wind   *wind.System
	conv   *conv.System
	tower  *tower.System
	cable  *cable.System
	grid   *grid.System
	drive  *yaw.Drive
}

func NewServer(
	addr string,
	st *store.Store,
	limits ns.Limits,
	recorder *audit.Recorder,
	bladeSys *blade.Blade,
	brakeSys *brake.System,
	pitchSys *pitch.System,
	yawSys *yaw.System,
	windSys *wind.System,
	convSys *conv.System,
	towerSys *tower.System,
	cableSys *cable.System,
	gridSys *grid.System,
	drive *yaw.Drive,
) *Server {
	server := &Server{
		addr:   addr,
		st:     st,
		limits: limits,
		audit:  recorder,
		blade:  bladeSys,
		brake:  brakeSys,
		pitch:  pitchSys,
		yaw:    yawSys,
		wind:   windSys,
		conv:   convSys,
		tower:  towerSys,
		cable:  cableSys,
		grid:   gridSys,
		drive:  drive,
	}
	server.router = chi.NewRouter()
	server.router.Use(chimiddleware.Recoverer, chimiddleware.RealIP, server.withLogging)
	server.routes()
	return server
}

func (s *Server) Start() error {
	return s.StartWithContext(context.Background())
}

func (s *Server) StartWithContext(ctx context.Context) error {
	go s.loop()
	httpServer := &http.Server{Addr: s.addr, Handler: s.router}
	errCh := make(chan error, 1)
	go func() {
		errCh <- httpServer.ListenAndServe()
	}()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpServer.Shutdown(shutdownCtx)
	}
}

func (s *Server) loop() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if err := s.applyDrive(); err != nil {
			s.audit.Append("console", "loop", err.Error())
		}
		if s.conv.Closed() {
			angle, _, _ := s.pitch.Setpoint().Current()
			_ = s.pitch.Move(angle)
		}
	}
}

func (s *Server) applyDrive() error {
	target, _, _ := s.drive.Current()
	return s.yaw.ApplyTarget(target)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
