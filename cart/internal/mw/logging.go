package mw

import (
	"go.uber.org/zap"
	"net/http"
	"route256.ozon.ru/project/cart/pkg/logger"
)

func Log(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("got http request",
			logger.FieldsWithTraceID(
				r.Context(),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
			)...,
		)
		next(w, r)
	}
}
