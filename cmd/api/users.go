package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/timmy1496/social/internal/store"
	"net/http"
	"strconv"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)

		return
	}

	ctx := r.Context()

	user, err := app.storage.Users.Get(ctx, userID)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.badRequestResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}
