package jwt

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"

	ssogrpc "restrocheck/internal/app/grpc/sso"
	"restrocheck/pkg/jwt"
	"restrocheck/pkg/logger/sl"
	resp "restrocheck/pkg/response"
)

type ContextKey string

const ClaimsKey ContextKey = "claims"
const IsAdminKey ContextKey = "is_admin"

// JWTAuthIsAdminMiddleware проверяет JWT, запрашивает SSO и добавляет данные в контекст.
func JWTAuthIsAdminMiddleware(log *slog.Logger, ssoapi *ssogrpc.Client, appSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const fn = "middleware.authentication.JWTAuthMiddleware"

			log = log.With(
				slog.String("fn", fn),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Warn("Authorization header is missing")
				resp.RespondWithError(log, w, r, http.StatusUnauthorized, errors.New("missing authorization header"), "missing authorization header")
				return
			}

			// Разбираем заголовок "Bearer <TOKEN>"
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" || headerParts[1] == "" {
				log.Warn("Invalid Authorization header format")
				resp.RespondWithError(log, w, r, http.StatusUnauthorized, errors.New("invalid authorization header"), "invalid authorization header")
				return
			}

			tokenString := headerParts[1]

			claims, err := jwt.ParseToken(tokenString, appSecret)
			if err != nil {
				log.Warn("Invalid JWT token", sl.Err(err))
				resp.RespondWithError(log, w, r, http.StatusUnauthorized, errors.New("invalid token"), "invalid token")
				return
			}

			log.Debug("JWT parsed successfully", slog.Any("claims", claims))

			// Проверяем, является ли пользователь админом через SSO
			isAdmin, err := ssoapi.IsAdmin(r.Context(), claims.ID)
			if err != nil {
				log.Error("Failed to verify admin status", sl.Err(err))
				resp.RespondWithError(log, w, r, http.StatusInternalServerError, errors.New("failed to verify admin status"), "internal server error")
				return
			}

			log.Debug("Admin status retrieved", slog.Bool("is_admin", isAdmin))

			// Если пользователь не админ, возвращаем ошибку
			if !isAdmin {
				log.Warn("Access denied: user is not an admin", slog.Int64("user_id", claims.ID))
				resp.RespondWithError(log, w, r, http.StatusForbidden, errors.New("access denied"), "forbidden")
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			ctx = context.WithValue(ctx, IsAdminKey, isAdmin)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
