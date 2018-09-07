package nats

import (
	"context"
	"encoding/json"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"

	"github.com/nats-io/go-nats"
)

// Subscriber wraps an endpoint and provides nats.MsgHandler.
type Subscriber struct {
	e            endpoint.Endpoint
	dec          DecodeRequestFunc
	enc          EncodeResponseFunc
	before       []RequestFunc
	after        []SubscriberResponseFunc
	errorEncoder ErrorEncoder
	logger       log.Logger
}

// NewSubscriber constructs a new subscriber, which provides nats.MsgHandler and wraps
// the provided endpoint.
func NewSubscriber(
	e endpoint.Endpoint,
	dec DecodeRequestFunc,
	enc EncodeResponseFunc,
	options ...SubscriberOption,
) *Subscriber {
	s := &Subscriber{
		e:            e,
		dec:          dec,
		enc:          enc,
		errorEncoder: DefaultErrorEncoder,
		logger:       log.NewNopLogger(),
	}
	for _, option := range options {
		option(s)
	}
	return s
}

// SubscriberOption sets an optional parameter for subscribers.
type SubscriberOption func(*Subscriber)

// SubscriberBefore functions are executed on the publisher request object before the
// request is decoded.
func SubscriberBefore(before ...RequestFunc) SubscriberOption {
	return func(s *Subscriber) { s.before = append(s.before, before...) }
}

// SubscriberAfter functions are executed on the subscriber reply after the
// endpoint is invoked, but before anything is published to the reply.
func SubscriberAfter(after ...SubscriberResponseFunc) SubscriberOption {
	return func(s *Subscriber) { s.after = append(s.after, after...) }
}

// SubscriberErrorEncoder is used to encode errors to the subscriber reply
// whenever they're encountered in the processing of a request. Clients can
// use this to provide custom error formatting. By default,
// errors will be published with the DefaultErrorEncoder.
func SubscriberErrorEncoder(ee ErrorEncoder) SubscriberOption {
	return func(s *Subscriber) { s.errorEncoder = ee }
}

// SubscriberErrorLogger is used to log non-terminal errors. By default, no errors
// are logged. This is intended as a diagnostic measure. Finer-grained control
// of error handling, including logging in more detail, should be performed in a
// custom SubscriberErrorEncoder which has access to the context.
func SubscriberErrorLogger(logger log.Logger) SubscriberOption {
	return func(s *Subscriber) { s.logger = logger }
}

// ServeMsg provides nats.MsgHandler.
func (s Subscriber) ServeMsg(nc *nats.Conn) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for _, f := range s.before {
			ctx = f(ctx, msg)
		}

		request, err := s.dec(ctx, msg)
		if err != nil {
			s.logger.Log("err", err)
			if msg.Reply == "" {
				return
			}
			s.errorEncoder(ctx, err, msg.Reply, nc)
			return
		}

		response, err := s.e(ctx, request)
		if err != nil {
			s.logger.Log("err", err)
			if msg.Reply == "" {
				return
			}
			s.errorEncoder(ctx, err, msg.Reply, nc)
			return
		}

		for _, f := range s.after {
			ctx = f(ctx, nc)
		}

		if msg.Reply == "" {
			return
		}

		if err := s.enc(ctx, msg.Reply, nc, response); err != nil {
			s.logger.Log("err", err)
			s.errorEncoder(ctx, err, msg.Reply, nc)
			return
		}
	}
}

// ErrorEncoder is responsible for encoding an error to the subscriber reply.
// Users are encouraged to use custom ErrorEncoders to encode errors to
// their replies, and will likely want to pass and check for their own error
// types.
type ErrorEncoder func(ctx context.Context, err error, reply string, nc *nats.Conn)

// NopRequestDecoder is a DecodeRequestFunc that can be used for requests that do not
// need to be decoded, and simply returns nil, nil.
func NopRequestDecoder(_ context.Context, _ *nats.Msg) (interface{}, error) {
	return nil, nil
}

// EncodeJSONResponse is a EncodeResponseFunc that serializes the response as a
// JSON object to the subscriber reply. Many JSON-over services can use it as
// a sensible default.
func EncodeJSONResponse(_ context.Context, reply string, nc *nats.Conn, response interface{}) error {
	b, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return nc.Publish(reply, b)
}

// DefaultErrorEncoder writes the error to the subscriber reply.
func DefaultErrorEncoder(_ context.Context, err error, reply string, nc *nats.Conn) {
	logger := log.NewNopLogger()

	type Response struct {
		Error string `json:"err"`
	}

	var response Response

	response.Error = err.Error()

	b, err := json.Marshal(response)
	if err != nil {
		logger.Log("err", err)
		return
	}

	if err := nc.Publish(reply, b); err != nil {
		logger.Log("err", err)
	}
}
