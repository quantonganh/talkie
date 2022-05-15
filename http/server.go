package http

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/quantonganh/talkie"
	"github.com/quantonganh/talkie/ui"
)

// Server represents an HTTP server
type Server struct {
	logger zerolog.Logger
	router *gin.Engine

	UserService    talkie.UserService
	TokenService   talkie.TokenService
	CommentService talkie.CommentService
}

// NewServer creates new HTTP server
func NewServer() *Server {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	zLogger := log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)
	log.Logger = zLogger

	s := &Server{
		logger: zLogger,
		router: gin.New(),
	}

	s.router.Use(logger.SetLogger())

	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"},
		AllowMethods:     []string{http.MethodPost, http.MethodGet, http.MethodPatch},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}))

	publicFS, _ := fs.Sub(ui.Public, "public")

	_ = fs.WalkDir(publicFS, ".", func(path string, d fs.DirEntry, err error) error {
		fmt.Println(path)
		return nil
	})

	s.router.GET("/", gin.WrapH(http.FileServer(http.FS(publicFS))))
	s.router.GET("/favicon.png", gin.WrapH(http.FileServer(http.FS(publicFS))))
	s.router.GET("/global.css", gin.WrapH(http.FileServer(http.FS(publicFS))))
	s.router.GET("/build/bundle.css", gin.WrapH(http.FileServer(http.FS(publicFS))))
	s.router.GET("/build/bundle.js", gin.WrapH(http.FileServer(http.FS(publicFS))))
	s.router.GET("/build/bundle.js.map", gin.WrapH(http.FileServer(http.FS(publicFS))))

	s.router.POST("/auth/google", s.Error(s.authGoogle()))

	s.router.GET("/token/refresh", s.Error(s.refreshTokenHandler()))

	s.router.POST("/comments", authJWT(), s.Error(s.createCommentHandler()))
	s.router.GET("/comments/:post_slug", s.Error(s.getCommentsHandler()))
	s.router.PATCH("/comments/:id", authJWT(), s.Error(s.updateCommentHandler()))
	s.router.GET("/comments", s.Error(s.listCommentsHandler()))

	return s
}

// Run starts listening and serving HTTP requests
func (s *Server) Run(port string) error {
	return s.router.Run(fmt.Sprintf(":%s", port))
}
