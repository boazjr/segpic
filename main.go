package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	fmt.Println("starting Segpic")

	p := NewPicSum(nil)
	records, err := p.ListImages()
	if err != nil {
		log.Panic("images website is unavailable", err)
	}
	db := NewDB()
	db.Start("")
	db.Seed(records)
	s := &Server{
		p:  p,
		db: db,
	}
	s.start()
}

// Server is the base object of the service
type Server struct {
	r  *chi.Mux
	db *DB
	p  *PicSum
}

func (s *Server) start() error {
	s.r = chi.NewRouter()

	s.r.Use(middleware.RequestID)
	s.r.Use(middleware.RealIP)
	s.r.Use(middleware.Logger)
	s.r.Use(middleware.Recoverer)

	s.r.Use(middleware.Timeout(60 * time.Second))

	s.r.Route("/api/images", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json; charset=utf-8"))
		r.Get("/", s.listImages)
		r.Patch("/{imageID}", s.flagImage)
	})

	s.fileServer(s.r)

	return http.ListenAndServe(":8081", s.r)
}

func (s *Server) fileServer(router *chi.Mux) {
	root := "./web/dist"
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}

func (s *Server) flagImage(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	imageID := chi.URLParam(r, "imageID")

	if imageID == "" {
		log.Print("missing imageID")
		http.Error(w, http.StatusText(422), http.StatusUnprocessableEntity)
		return
	}
	img, err := s.db.FlagImage(imageID)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(422), http.StatusUnprocessableEntity)
		return
	}

	b, err := json.Marshal(toImageRes(img))
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
}

func (s *Server) listImages(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	imgs, err := s.db.ListImages()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	b, err := json.Marshal(toImagesRes(imgs...))
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
}

type ImageRes struct {
	ID          string `json:"id"`
	Author      string `json:"author"`
	Flag        bool   `json:"flag"`
	DownloadURL string `json:"download_url"`
}

func toImagesRes(imgs ...DBImage) []ImageRes {
	res := make([]ImageRes, len(imgs))
	for i := range imgs {
		res[i] = toImageRes(imgs[i])
	}
	return res
}
func toImageRes(img DBImage) ImageRes {
	return ImageRes{
		ID:          img.ID,
		Author:      img.Author,
		Flag:        img.Flag,
		DownloadURL: img.MakeURL(500),
	}
}
