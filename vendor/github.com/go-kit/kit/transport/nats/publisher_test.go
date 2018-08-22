package nats_test

import (
	"context"
	"strings"
	"testing"
	"time"

	natstransport "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/go-nats"
)

func TestPublisher(t *testing.T) {
	var (
		testdata = "testdata"
		encode   = func(context.Context, *nats.Msg, interface{}) error { return nil }
		decode   = func(_ context.Context, msg *nats.Msg) (interface{}, error) {
			return TestResponse{string(msg.Data), ""}, nil
		}
	)

	nc := newNatsConn(t)
	defer nc.Close()

	sub, err := nc.QueueSubscribe("natstransport.test", "natstransport", func(msg *nats.Msg) {
		if err := nc.Publish(msg.Reply, []byte(testdata)); err != nil {
			t.Fatal(err)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	publisher := natstransport.NewPublisher(
		nc,
		"natstransport.test",
		encode,
		decode,
	)

	res, err := publisher.Endpoint()(context.Background(), struct{}{})
	if err != nil {
		t.Fatal(err)
	}

	response, ok := res.(TestResponse)
	if !ok {
		t.Fatal("response should be TestResponse")
	}
	if want, have := testdata, response.String; want != have {
		t.Errorf("want %q, have %q", want, have)
	}

}

func TestPublisherBefore(t *testing.T) {
	var (
		testdata = "testdata"
		encode   = func(context.Context, *nats.Msg, interface{}) error { return nil }
		decode   = func(_ context.Context, msg *nats.Msg) (interface{}, error) {
			return TestResponse{string(msg.Data), ""}, nil
		}
	)

	nc := newNatsConn(t)
	defer nc.Close()

	sub, err := nc.QueueSubscribe("natstransport.test", "natstransport", func(msg *nats.Msg) {
		if err := nc.Publish(msg.Reply, msg.Data); err != nil {
			t.Fatal(err)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	publisher := natstransport.NewPublisher(
		nc,
		"natstransport.test",
		encode,
		decode,
		natstransport.PublisherBefore(func(ctx context.Context, msg *nats.Msg) context.Context {
			msg.Data = []byte(strings.ToUpper(string(testdata)))
			return ctx
		}),
	)

	res, err := publisher.Endpoint()(context.Background(), struct{}{})
	if err != nil {
		t.Fatal(err)
	}

	response, ok := res.(TestResponse)
	if !ok {
		t.Fatal("response should be TestResponse")
	}
	if want, have := strings.ToUpper(testdata), response.String; want != have {
		t.Errorf("want %q, have %q", want, have)
	}

}

func TestPublisherAfter(t *testing.T) {
	var (
		testdata = "testdata"
		encode   = func(context.Context, *nats.Msg, interface{}) error { return nil }
		decode   = func(_ context.Context, msg *nats.Msg) (interface{}, error) {
			return TestResponse{string(msg.Data), ""}, nil
		}
	)

	nc := newNatsConn(t)
	defer nc.Close()

	sub, err := nc.QueueSubscribe("natstransport.test", "natstransport", func(msg *nats.Msg) {
		if err := nc.Publish(msg.Reply, []byte(testdata)); err != nil {
			t.Fatal(err)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	publisher := natstransport.NewPublisher(
		nc,
		"natstransport.test",
		encode,
		decode,
		natstransport.PublisherAfter(func(ctx context.Context, msg *nats.Msg) context.Context {
			msg.Data = []byte(strings.ToUpper(string(msg.Data)))
			return ctx
		}),
	)

	res, err := publisher.Endpoint()(context.Background(), struct{}{})
	if err != nil {
		t.Fatal(err)
	}

	response, ok := res.(TestResponse)
	if !ok {
		t.Fatal("response should be TestResponse")
	}
	if want, have := strings.ToUpper(testdata), response.String; want != have {
		t.Errorf("want %q, have %q", want, have)
	}

}

func TestPublisherTimeout(t *testing.T) {
	var (
		encode = func(context.Context, *nats.Msg, interface{}) error { return nil }
		decode = func(_ context.Context, msg *nats.Msg) (interface{}, error) {
			return TestResponse{string(msg.Data), ""}, nil
		}
	)

	nc := newNatsConn(t)
	defer nc.Close()

	ch := make(chan struct{})
	defer close(ch)

	sub, err := nc.QueueSubscribe("natstransport.test", "natstransport", func(msg *nats.Msg) {
		<-ch
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	publisher := natstransport.NewPublisher(
		nc,
		"natstransport.test",
		encode,
		decode,
		natstransport.PublisherTimeout(time.Second),
	)

	_, err = publisher.Endpoint()(context.Background(), struct{}{})
	if err != context.DeadlineExceeded {
		t.Errorf("want %s, have %s", context.DeadlineExceeded, err)
	}
}

func TestEncodeJSONRequest(t *testing.T) {
	var data string

	nc := newNatsConn(t)
	defer nc.Close()

	sub, err := nc.QueueSubscribe("natstransport.test", "natstransport", func(msg *nats.Msg) {
		data = string(msg.Data)

		if err := nc.Publish(msg.Reply, []byte("")); err != nil {
			t.Fatal(err)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	publisher := natstransport.NewPublisher(
		nc,
		"natstransport.test",
		natstransport.EncodeJSONRequest,
		func(context.Context, *nats.Msg) (interface{}, error) { return nil, nil },
	).Endpoint()

	for _, test := range []struct {
		value interface{}
		body  string
	}{
		{nil, "null"},
		{12, "12"},
		{1.2, "1.2"},
		{true, "true"},
		{"test", "\"test\""},
		{struct {
			Foo string `json:"foo"`
		}{"foo"}, "{\"foo\":\"foo\"}"},
	} {
		if _, err := publisher(context.Background(), test.value); err != nil {
			t.Fatal(err)
			continue
		}

		if data != test.body {
			t.Errorf("%v: actual %#v, expected %#v", test.value, data, test.body)
		}
	}

}
