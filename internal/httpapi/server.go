package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LYH2263/go-geoalert"
)

// Server HTTP 管理 API。
type Server struct {
	eng *geoalert.Engine
	web string
	mux *http.ServeMux
}

func New(eng *geoalert.Engine, webDir string) *Server {
	s := &Server{eng: eng, web: webDir, mux: http.NewServeMux()}
	s.routes()
	s.attachExtra()
	return s
}

func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) routes() {
	if s.web != "" {
		s.mux.Handle("/", http.FileServer(http.Dir(s.web)))
	}
	s.mux.HandleFunc("/api/stats", s.handleStats)
	s.mux.HandleFunc("/api/fences", s.handleFences)
	s.mux.HandleFunc("/api/fence", s.handleFence)
	s.mux.HandleFunc("/api/ingest", s.handleIngest)
	s.mux.HandleFunc("/api/ingest/batch", s.handleIngestBatch)
	s.mux.HandleFunc("/api/alerts", s.handleAlerts)
	s.mux.HandleFunc("/api/tracks", s.handleTracks)
	s.mux.HandleFunc("/api/snapshot", s.handleSnapshot)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, s.eng.Stats())
}

func (s *Server) handleFences(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, s.eng.ListFences())
}

func (s *Server) handleFence(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var f geoalert.Fence
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if err := s.eng.RegisterFence(f); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		writeJSON(w, map[string]string{"ok": "1", "id": f.ID})
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if err := s.eng.UnregisterFence(id); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		writeJSON(w, map[string]string{"ok": "1"})
	default:
		id := r.URL.Query().Get("id")
		v, ok := s.eng.GetFence(id)
		if !ok {
			http.Error(w, "not found", 404)
			return
		}
		writeJSON(w, v)
	}
}

func (s *Server) handleIngest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
	var p geoalert.TrackPoint
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if p.At.IsZero() {
		p.At = time.Now().UTC()
	}
	if err := s.eng.Ingest(p); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, map[string]string{"ok": "1"})
}

func (s *Server) handleIngestBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
	var pts []geoalert.TrackPoint
	if err := json.NewDecoder(r.Body).Decode(&pts); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := s.eng.IngestBatch(r.Context(), pts); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, map[string]any{"ok": "1", "n": len(pts)})
}

func (s *Server) handleAlerts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("drain") == "1" {
		writeJSON(w, s.eng.DrainAlerts())
		return
	}
	writeJSON(w, s.eng.Alerts())
}

func (s *Server) handleTracks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, s.eng.ListTracks())
}

func (s *Server) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, s.eng.Snapshot())
}
