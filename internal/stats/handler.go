package stats

import (
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

var _ stats.Handler = (*Handler)(nil)

type Handler struct {
	logger zerolog.Logger
}

func NewHandler(logger zerolog.Logger) *Handler {
	return &Handler{logger}
}

func (h *Handler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	h.logger.Debug().
		Str("FullMethodName", info.FullMethodName).
		Bool("FailFast", info.FailFast).
		Msg("TagRPC")
	return ctx
}

func (h *Handler) HandleRPC(ctx context.Context, s stats.RPCStats) {
	switch st := s.(type) {
	case *stats.Begin:
		h.logger.Debug().
			Bool("IsClient", st.IsClient()).
			Time("BeginTime", st.BeginTime).
			Bool("FailFast", st.FailFast).
			Msg("stats Begin")
	case *stats.OutHeader:
		h.logger.Debug().
			Bool("IsClient", st.IsClient()).
			Str("Compression", st.Compression).
			Interface("Header", st.Header).
			Str("FullMethod", st.FullMethod).
			Stringer("RemoteAddr", st.RemoteAddr).
			Stringer("LocalAddr", st.LocalAddr).
			Msg("stats OutHeader")
	case *stats.OutPayload:
		h.logger.Debug().
			Bool("IsClient", st.IsClient()).
			Interface("Payload", st.Payload).
			Bytes("Data", st.Data).
			Int("Length", st.Length).
			Int("WireLength", st.WireLength).
			Time("SentTime", st.SentTime).
			Msg("stats OutPayload")
	case *stats.OutTrailer:
		h.logger.Debug().
			Bool("IsClient", st.IsClient()).
			Int("WireLength", st.WireLength).
			Interface("Trailer", st.Trailer).
			Msg("stats OutTrailer")
	case *stats.InHeader:
		h.logger.Debug().
			Bool("IsClient", st.IsClient()).
			Int("WireLength", st.WireLength).
			Str("Compression", st.Compression).
			Interface("Header", st.Header).
			Str("FullMethod", st.FullMethod).
			Stringer("RemoteAddr", st.RemoteAddr).
			Stringer("LocalAddr", st.LocalAddr).
			Msg("stats InHeader")
	case *stats.InPayload:
		h.logger.Debug().
			Bool("IsClient", st.IsClient()).
			Interface("Payload", st.Payload).
			Bytes("Data", st.Data).
			Int("Length", st.Length).
			Int("WireLength", st.WireLength).
			Time("RecvTime", st.RecvTime).
			Msg("stats InPayload")
	case *stats.InTrailer:
		h.logger.Debug().
			Bool("IsClient", st.IsClient()).
			Int("WireLength", st.WireLength).
			Interface("Trailer", st.Trailer).
			Msg("stats InTrailer")
	case *stats.End:
		ev := h.logger.Debug().
			Bool("IsClient", st.IsClient()).
			Time("BeginTime", st.BeginTime).
			Time("EndTime", st.EndTime).
			Interface("Trailer", st.Trailer).
			Err(st.Error)
		gRPCst, ok := status.FromError(st.Error)
		if ok {
			ev = ev.Interface("details", gRPCst.Details())
		}
		ev.Msg("stats End")
	default:
		h.logger.Error().
			Interface("stats", st).Msg("unknwon")
	}
}

func (h *Handler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	h.logger.Debug().
		Stringer("RemoteAddr", info.RemoteAddr).
		Stringer("LocalAddr", info.LocalAddr).
		Msg("TagConn")
	return context.Background()
}

func (h *Handler) HandleConn(ctx context.Context, st stats.ConnStats) {
	h.logger.Debug().
		Bool("IsClient", st.IsClient()).
		Msg("HandleConn")
}
