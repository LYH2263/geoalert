package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *Server) attachExtra() {
	s.mux.HandleFunc("/api/evaluate", s.handleEvaluate)
	s.mux.HandleFunc("/api/inside", s.handleInside)
	s.mux.HandleFunc("/api/distance", s.handleDistance)
	s.mux.HandleFunc("/api/purge", s.handlePurge)
}

func (s *Server) handleEvaluate(w http.ResponseWriter, r *http.Request) {
	fid := r.URL.Query().Get("fence")
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lng, _ := strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
	inside, err := s.eng.EvaluatePoint(fid, lat, lng)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, map[string]any{"inside": inside})
}

func (s *Server) handleInside(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, s.eng.ObjectsInside(r.URL.Query().Get("fence")))
}

func (s *Server) handleDistance(w http.ResponseWriter, r *http.Request) {
	d, err := s.eng.DistanceToFence(r.URL.Query().Get("object"), r.URL.Query().Get("fence"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, map[string]any{"meters": d})
}

func (s *Server) handlePurge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
	var body struct {
		ObjectID string `json:"object_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	s.eng.PurgeObject(body.ObjectID)
	writeJSON(w, map[string]string{"ok": "1"})
}
