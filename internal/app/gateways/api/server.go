package api

import (
	"log"
	"net/http"

	"github.com/literalog/library/internal/app/domain/author"
	"github.com/literalog/library/internal/app/domain/book"
	"github.com/literalog/library/internal/app/domain/genre"
	"github.com/literalog/library/internal/app/domain/series"
	"github.com/literalog/library/internal/app/gateways/apis"
	"github.com/literalog/library/internal/app/gateways/database/mongodb"

	"github.com/gorilla/mux"
)

type Server struct {
	port     string
	logLevel int
	router   *mux.Router
}

func NewServer(port string) Server {

	router := mux.NewRouter()

	storage, err := mongodb.NewMongoStorage()
	if err != nil {
		log.Fatal(err)
	}

	db := storage.Client.Database("library")

	authorRepository := mongodb.NewAuthorRepository(db.Collection("authors"))
	authorService := author.NewService(authorRepository)
	authorHandler := author.NewHandler(authorService)

	seriesRepository := mongodb.NewSeriesRepository(db.Collection("series"))
	seriesService := series.NewService(seriesRepository)
	seriesHandler := series.NewHandler(seriesService)

	genreRepository := mongodb.NewGenreRepository(db.Collection("genre"))
	genreService := genre.NewService(genreRepository)
	genreHandler := genre.NewHandler(genreService)

	isbnRepository, err := apis.NewGBooksAPI("", "https://www.googleapis.com/books/v1")
	if err != nil {
		log.Fatal(err)
	}

	bookRepository := mongodb.NewBookRepository(db.Collection("books"))
	bookService := book.NewService(bookRepository, isbnRepository, authorService, seriesService, genreService)
	bookHandler := book.NewHandler(bookService)

	router.PathPrefix("/authors").Handler(authorHandler.Routes())
	router.PathPrefix("/series").Handler(seriesHandler.Routes())
	router.PathPrefix("/genres").Handler(genreHandler.Routes())
	router.PathPrefix("/books").Handler(bookHandler.Routes())

	return Server{
		port:     port,
		logLevel: 1,
		router:   router,
	}

}

func (s *Server) ServeHttp() error {
	log.Println("Server listening on", s.port)
	return http.ListenAndServe(s.port, s.router)
}
