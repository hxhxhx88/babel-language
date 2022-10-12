package backend

// func (s *mServer) serveFrontend(w http.ResponseWriter, r *http.Request) {
// 	if h.options.frontend == "" {
// 		http.Error(w, "not implemented", http.StatusNotImplemented)
// 		return
// 	}

// 	path := r.URL.Path

// 	// pkger requires absolute path
// 	if !strings.HasPrefix(path, "/") {
// 		path = "/" + path
// 	}

// 	if path == "/" || strings.HasPrefix(path, "/_") {
// 		path = "/index.html"
// 	}

// 	file, err := h.options.frontend.Open(path)
// 	if err != nil {
// 		logging.Error(err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	mimeType := mime.TypeByExtension(filepath.Ext(path))
// 	if mimeType != "" {
// 		w.Header().Set("Content-Type", mimeType)
// 	}

// 	io.Copy(w, file)
// }
