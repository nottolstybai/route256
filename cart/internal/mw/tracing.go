package mw

import (
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func HandleWithSpan(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("default").Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL), trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		fn(w, r.WithContext(ctx))
	}
}
