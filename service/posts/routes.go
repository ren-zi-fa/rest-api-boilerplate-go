package posts

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/types"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/utils"
)

type Handler struct {
	store types.PostStore
}

func NewHandler(store types.PostStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoute(router chi.Router) {

	router.Get("/posts", h.handleGetPosts)
	router.Get("/posts/{id}", h.handleGetPost)
	router.Post("/posts", h.handleCreatePost)
	router.Delete("/posts/{id}", h.handleDeletePost)
}

func (h *Handler) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.store.GetPosts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, posts)

}

func (h *Handler) handleGetPost(w http.ResponseWriter, r *http.Request) {
	postIdStr := chi.URLParam(r, "id")
	postID, err := strconv.Atoi(postIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid post id"))
		return
	}

	post, err := h.store.GetPostByID(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("post not found with id %d", postID))
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, post)
}

func (h *Handler) handleDeletePost(w http.ResponseWriter, r *http.Request) {

	// chi.URL param will catch path param example /post/{id}
	postIdStr := chi.URLParam(r, "id")
	postID, err := strconv.Atoi(postIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid post id"))
		return
	}

	rowsAffected, err := h.store.DeletePostByID(postID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("post not found with id %d", postID))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "post deleted successfully"})
}

func (h *Handler) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	var post types.CreatePostPayload
	err := utils.ParseJSON(r, &post)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// this method will compare contents of struct post and  struct post payload 
	if err := utils.Validate.Struct(post); err != nil {
		fieldErrors := utils.FormatValidationError(err)
		utils.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":  "invalid payload",
			"fields": fieldErrors,
		})
		return
	}
	id, err := h.store.CreatePost(post)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]int64{"id": id})
}
