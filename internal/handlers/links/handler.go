package links

import (
	"github.com/gin-gonic/gin"
	"github.com/kurnosovmak/short-link/internal/services/link_service"
	httpserver "github.com/kurnosovmak/short-link/pkg/http-server"
	"github.com/kurnosovmak/short-link/pkg/logging"
	"net/http"
)

const (
	urlCreateLink   = "/api/link"
	urlRedirectLink = "/api/link/:key"
)

type Handler struct {
	LinkService link_service.LinkService
	Logger      logging.Logger
}

func NewHandler(linkService link_service.LinkService, logger logging.Logger) *Handler {
	handler := Handler{
		LinkService: linkService,
		Logger:      logger,
	}
	return &handler
}

func (h *Handler) Register(router *httpserver.HttpServer) {

	router.Handle(http.MethodGet, urlRedirectLink, h.RedirectLink)
	router.Handle(http.MethodPost, urlCreateLink, h.CreateShortLink)
}

func (h *Handler) RedirectLink(r *gin.Context) {
	var dto = link_service.RedirectLinkDTO{
		Key: r.Param("key"),
	}
	if dto.Key == "" {
		r.JSON(http.StatusOK, gin.H{
			"message": "key in empty",
		})
		return
	}

	link, err := h.LinkService.RedirectLink(r.Request.Context(), dto)
	if err != nil {
		h.Logger.Error(err)
		r.JSON(http.StatusOK, gin.H{
			"message": "link not found",
		})
		return
	}
	r.Redirect(http.StatusFound, link)
}

func (h *Handler) CreateShortLink(r *gin.Context) {
	var dto link_service.CreateLinkDTO

	err := r.ShouldBindJSON(&dto)
	if err != nil || dto.Link == "" {
		h.Logger.Info(err)
		r.JSON(http.StatusOK, gin.H{
			"message": "data wrong",
		})
		return
	}

	key, err := h.LinkService.CreateLink(r.Request.Context(), dto)
	if err != nil {
		h.Logger.Error(err)
		return
	}
	r.JSON(http.StatusOK, gin.H{
		"link": urlCreateLink + "/" + key,
	})
}
