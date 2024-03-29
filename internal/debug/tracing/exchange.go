package tracing

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/stealthrocket/timecraft/internal/stream"
	"gopkg.in/yaml.v3"
)

type Request struct {
	Time time.Time
	Span time.Duration
	Err  error
	msg  ConnMessage
}

func (r Request) Format(w fmt.State, v rune) {
	if r.msg != nil {
		r.msg.Format(w, v)
	} else if r.Err != nil {
		fmt.Fprintf(w, "Request error: %s", r.Err.Error())
	} else {
		fmt.Fprint(w, "unknown")
	}
}

type Response struct {
	Time time.Duration
	Span time.Duration
	Err  error
	msg  ConnMessage
}

func (r Response) Format(w fmt.State, v rune) {
	if r.msg != nil {
		r.msg.Format(w, v)
	} else if r.Err != nil {
		fmt.Fprintf(w, "Response error: %s", r.Err.Error())
	} else {
		fmt.Fprint(w, "unknown")
	}
}

// Exchange values represent the exchange of a request and response between
// network peers.
type Exchange struct {
	Link Link
	Req  Request
	Res  Response
}

func (e Exchange) Format(w fmt.State, v rune) {
	fmt.Fprintf(w, "%s %s %s > %s",
		formatTime(e.Req.Time.Add(e.Res.Time)),
		e.Req.msg.Conn().Protocol().Name(),
		socketAddressString(e.Link.Src),
		socketAddressString(e.Link.Dst))

	if w.Flag('+') {
		fmt.Fprintf(w, "\n")
		e.Req.Format(w, v)
		e.Res.Format(w, v)
	} else {
		fmt.Fprintf(w, ": ")
		e.Req.Format(w, v)

		fmt.Fprintf(w, " => ")
		e.Res.Format(w, v)
	}
}

func (e Exchange) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.marshal())
}

func (e Exchange) MarshalYAML() (any, error) {
	return e.marshal(), nil
}

func (e *Exchange) marshal() *exchange {
	return &exchange{
		Link: e.Link,
		Req: request{
			Time: e.Req.Time,
			Span: e.Req.Span,
			Err:  errorString(e.Req.Err),
			Msg:  e.Req.msg.Marshal(),
		},
		Res: response{
			Time: e.Req.Time.Add(e.Res.Time),
			Span: e.Res.Span,
			Err:  errorString(e.Res.Err),
			Msg:  e.Res.msg.Marshal(),
		},
	}
}

type request struct {
	Time time.Time     `json:"time"            yaml:"time"`
	Span time.Duration `json:"span"            yaml:"span"`
	Err  string        `json:"error,omitempty" yaml:"error,omitempty"`
	Msg  any           `json:"message"         yaml:"message"`
}

type response struct {
	Time time.Time     `json:"time"            yaml:"time"`
	Span time.Duration `json:"span"            yaml:"span"`
	Err  string        `json:"error,omitempty" yaml:"error,omitempty"`
	Msg  any           `json:"message"         yaml:"message"`
}

type exchange struct {
	Link Link     `json:"link"     yaml:"link"`
	Req  request  `json:"request"  yaml:"request"`
	Res  response `json:"response" yaml:"response"`
}

var (
	_ fmt.Formatter  = Exchange{}
	_ json.Marshaler = Exchange{}
	_ yaml.Marshaler = Exchange{}
)

// ExchangeReader is a reader of Exchange values. Instances of ExchangeReader
// consume network messages from a reader of Message values and reconstruct the
// relationship between requests and responses.
type ExchangeReader struct {
	Messages stream.Reader[Message]

	inflight  map[int64]Message
	messages  []Message
	exchanges []Exchange
	offset    int
}

func (r *ExchangeReader) Read(exchanges []Exchange) (n int, err error) {
	if r.inflight == nil {
		r.inflight = make(map[int64]Message)
	}
	if len(r.messages) == 0 {
		r.messages = make([]Message, 1000)
	}

	for {
		if r.offset < len(r.exchanges) {
			n = copy(exchanges, r.exchanges[r.offset:])
			if r.offset += n; r.offset == len(r.exchanges) {
				r.offset, r.exchanges = 0, r.exchanges[:0]
			}
			return n, nil
		}

		numMessages, err := stream.ReadFull(r.Messages, r.messages)
		if numMessages == 0 {
			if len(r.inflight) > 0 {
				// The stream was interrupted before receiving a response for
				// some of the inflight requests; generate exchanges with the
				// response error set to io.ErrUnexpectedEOF so we don't mask
				// them from the application.
				for _, req := range r.inflight {
					r.exchanges = append(r.exchanges, Exchange{
						Link: req.Link,
						Req: Request{
							Time: req.Time,
							Span: req.Span,
							Err:  req.Err,
							msg:  req.msg,
						},
						Res: Response{
							Err: io.ErrUnexpectedEOF,
						},
					})
				}
				for id := range r.inflight {
					delete(r.inflight, id)
				}
				sortExchanges(r.exchanges)
				continue
			}
			switch err {
			case nil:
				err = io.ErrNoProgress
			case io.ErrUnexpectedEOF:
				err = io.EOF
			}
			return 0, err
		}

		for i := range r.messages[:numMessages] {
			m := &r.messages[i]

			req, ok := r.inflight[m.id]
			if !ok {
				r.inflight[m.id] = *m
			} else {
				delete(r.inflight, m.id)
				r.exchanges = append(r.exchanges, Exchange{
					Link: req.Link,
					Req: Request{
						Time: req.Time,
						Span: req.Span,
						Err:  req.Err,
						msg:  req.msg,
					},
					Res: Response{
						Time: m.Time.Sub(req.Time),
						Span: m.Span,
						Err:  m.Err,
						msg:  m.msg,
					},
				})
			}
		}

		sortExchanges(r.exchanges)
	}
}

func sortExchanges(exchanges []Exchange) {
	slices.SortFunc(exchanges, func(e1, e2 Exchange) int {
		t1 := e1.Req.Time.Add(e1.Res.Time + e1.Res.Span)
		t2 := e2.Req.Time.Add(e2.Res.Time + e2.Res.Span)
		return cmp.Compare(t1.UnixNano(), t2.UnixNano())
	})
}
