package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/timmy1496/social/internal/store"
	"net/http"
	"strconv"
)

type postKey string

const postCtx postKey = "posts"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)

		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)

		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	}

	ctx := r.Context()

	if err := app.storage.Posts.Create(ctx, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())

		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())

		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.storage.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)

		return
	}

	post.Comments = comments

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)

		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)

		return
	}

	if err := app.storage.Posts.Delete(ctx, postID); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	payload := UpdatePostPayload{}

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)

		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)

		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.storage.Posts.Update(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)

		return
	}

}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		postID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)

			return
		}

		ctx := r.Context()

		post, err := app.storage.Posts.GetByID(ctx, postID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}

			return
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)

	return post
}
