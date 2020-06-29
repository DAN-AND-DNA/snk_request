package snk_request

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	test_instances := []struct {
		description string
		route       string
		req_header  []string
		resp_header map[string]string
		request     []byte
		response    []byte
	}{
		{
			description: "text response",
			route:       "/text-resp",
			req_header:  nil,
			resp_header: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			request:     nil,
			response:    []byte("text response"),
		},
		{
			description: "default text request",
			route:       "/default-text-req",
			req_header:  nil,
			resp_header: nil,
			request:     []byte("default text request"),
			response:    nil,
		},
		{
			description: "text request and response",
			route:       "/text-req",
			req_header:  nil,
			resp_header: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			request:     []byte("text request"),
			response:    []byte("text response"),
		},
		{
			description: "set header",
			route:       "/set-header",
			req_header:  []string{"snk", "src"},
			resp_header: nil,
			request:     nil,
			response:    nil,
		},
		{
			description: "default json request",
			route:       "/json-req",
			req_header:  nil,
			resp_header: nil,
			request:     []byte(`{"name": "dan"}`),
			response:    nil,
		},
		{
			description: "default json reponse",
			route:       "/json-resp",
			req_header:  nil,
			resp_header: map[string]string{"Content-Type": "application/json"},
			request:     nil,
			response:    []byte(`{"age": 28}`),
		},
		{
			description: "json request reponse",
			route:       "/json-req-resp",
			req_header:  nil,
			resp_header: map[string]string{"Content-Type": "application/json"},
			request:     []byte(`{"name": "dan"}`),
			response:    []byte(`{"age": 28}`),
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, r.Method, "GET", "should get method")

		switch r.URL.Path {
		case test_instances[0].route:
			_, _ = w.Write(test_instances[0].response)
		case test_instances[1].route:
			// check header
			assert.Equalf(t, "text/plain", r.Header.Get("Content-Type"), test_instances[1].description)
			body, _ := ioutil.ReadAll(r.Body)

			// check body
			assert.Equalf(t, (string)(test_instances[1].request), (string)(body), test_instances[1].description)

			_, _ = w.Write(test_instances[1].response)
		case test_instances[2].route:
			// check header
			assert.Equalf(t, "text/plain", r.Header.Get("Content-Type"), test_instances[2].description)
			body, _ := ioutil.ReadAll(r.Body)

			// check body
			assert.Equalf(t, (string)(test_instances[2].request), (string)(body), test_instances[2].description)

			_, _ = w.Write(test_instances[2].response)
		case test_instances[3].route:
			// check header
			assert.Equalf(t, "src", r.Header.Get("snk"), test_instances[3].description)

			_, _ = w.Write(test_instances[3].response)

		case test_instances[4].route:
			// check header
			assert.Equalf(t, "application/json", r.Header.Get("Content-Type"), test_instances[4].description)
			body, _ := ioutil.ReadAll(r.Body)

			// check body
			assert.Equalf(t, (string)(test_instances[4].request), (string)(body), test_instances[4].description)

			_, _ = w.Write(test_instances[4].response)

		case test_instances[5].route:
			// check header
			assert.Equalf(t, "", r.Header.Get("Content-Type"), test_instances[5].description)
			body, _ := ioutil.ReadAll(r.Body)

			// check body
			assert.Equalf(t, (string)(test_instances[5].request), (string)(body), test_instances[5].description)

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(test_instances[5].response)

		case test_instances[6].route:
			// check header
			assert.Equalf(t, "application/json", r.Header.Get("Content-Type"), test_instances[6].description)
			body, _ := ioutil.ReadAll(r.Body)

			// check body
			assert.Equalf(t, (string)(test_instances[6].request), (string)(body), test_instances[6].description)

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(test_instances[6].response)
		}

	}))

	defer ts.Close()

	request := New()
	for _, test := range test_instances {

		if resp, body, err := request.Get(ts.URL + test.route).Set(test.req_header...).Send(test.request).End(); err != nil {
			// 1. check error
			assert.Equalf(t, nil, err, test.description)
			_ = body
		} else {

			// 2. check header
			for k, v := range test.resp_header {
				assert.Equalf(t, v, resp.Header.Get(k), test.description)
			}

			// 3. check content of body
			assert.Equalf(t, (string)(test.response), (string)(body), test.description)
		}
	}
}

func TestPost(t *testing.T) {
	test_instances := []struct {
		description string
		route       string
		req_header  []string
		resp_header map[string]string
		request     []byte
		response    []byte
	}{
		{
			description: "json request reponse",
			route:       "/json-req-resp",
			req_header:  []string{"snk", "src"},
			resp_header: map[string]string{"Content-Type": "application/json"},
			request:     []byte(`{"name": "dan"}`),
			response:    []byte(`{"age": 28}`),
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, r.Method, "POST", "should get method")

		switch r.URL.Path {
		case test_instances[0].route:
			// check header
			assert.Equalf(t, "src", r.Header.Get("snk"), test_instances[0].description)
			assert.Equalf(t, "application/json", r.Header.Get("Content-Type"), test_instances[0].description)
			body, _ := ioutil.ReadAll(r.Body)

			// check body
			assert.Equalf(t, (string)(test_instances[0].request), (string)(body), test_instances[0].description)

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(test_instances[0].response)

		}

	}))

	defer ts.Close()

	request := New()
	for _, test := range test_instances {
		if resp, body, err := request.Post(ts.URL + test.route).Set(test.req_header...).Send(test.request).End(); err != nil {
			// 1. check error
			assert.Equalf(t, nil, err, test.description)
			_ = body
		} else {

			// 2. check header
			for k, v := range test.resp_header {
				assert.Equalf(t, v, resp.Header.Get(k), test.description)
			}

			// 3. check content of body
			assert.Equalf(t, (string)(test.response), (string)(body), test.description)
		}
	}
}

func TestTimeout(t *testing.T) {
	test_instances := []struct {
		description string
		route       string
		req_header  []string
		resp_header map[string]string
		request     []byte
		response    []byte
	}{
		{
			description: "json request timeout",
			route:       "/json-req-timeout",
			req_header:  []string{"snk", "src"},
			resp_header: map[string]string{"Content-Type": "application/json"},
			request:     []byte(`{"name": "dan"}`),
			response:    []byte(`{"age": 28}`),
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, r.Method, "POST", "should get method")

		switch r.URL.Path {
		case test_instances[0].route:
			// check header
			assert.Equalf(t, "src", r.Header.Get("snk"), test_instances[0].description)
			assert.Equalf(t, "application/json", r.Header.Get("Content-Type"), test_instances[0].description)
			body, _ := ioutil.ReadAll(r.Body)

			// check body
			assert.Equalf(t, (string)(test_instances[0].request), (string)(body), test_instances[0].description)

			w.Header().Set("Content-Type", "application/json")
			time.Sleep(4 * time.Second)
			_, _ = w.Write(test_instances[0].response)
		}
	}))

	defer ts.Close()

	request := New_timeout(3, 3, 10)
	for _, test := range test_instances {
		if resp, body, err := request.Post(ts.URL + test.route).Set(test.req_header...).Send(test.request).End(); err != nil {
			// 1. check error
			assert.NotNil(t, err)
			_ = body
		} else {

			// 2. check header
			for k, v := range test.resp_header {
				assert.Equalf(t, v, resp.Header.Get(k), test.description)
			}

			// 3. check content of body
			assert.Equalf(t, (string)(test.response), (string)(body), test.description)
		}
	}
}

func TestConcurrent(t *testing.T) {
	test_instances := []struct {
		description string
		route       string
		req_header  []string
		resp_header map[string]string
		request     []byte
		response    []byte
	}{
		{
			description: "json request 1",
			route:       "/json-req1",
			req_header:  []string{"snk1", "src1"},
			resp_header: map[string]string{"Content-Type": "application/json"},
			request:     []byte(`{"name": "dan"}`),
			response:    []byte(`{"age": 28}`),
		},

		{
			description: "json request 2",
			route:       "/json-req2",
			req_header:  []string{"snk2", "src2"},
			resp_header: map[string]string{"Content-Type": "application/json"},
			request:     []byte(`{"name": "lulu"}`),
			response:    []byte(`{"age": 30}`),
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, r.Method, "POST", "should get method")

		switch r.URL.Path {
		case test_instances[0].route:
			// check header
			assert.Equalf(t, "src1", r.Header.Get("snk1"), test_instances[0].description)
			assert.Equalf(t, "application/json", r.Header.Get("Content-Type"), test_instances[0].description)
			body, _ := ioutil.ReadAll(r.Body)

			// check body
			assert.Equalf(t, (string)(test_instances[0].request), (string)(body), test_instances[0].description)

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(test_instances[0].response)

		case test_instances[1].route:
			// check header
			assert.Equalf(t, "src2", r.Header.Get("snk2"), test_instances[1].description)
			assert.Equalf(t, "application/json", r.Header.Get("Content-Type"), test_instances[1].description)
			body, _ := ioutil.ReadAll(r.Body)

			// check body
			assert.Equalf(t, (string)(test_instances[1].request), (string)(body), test_instances[1].description)

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(test_instances[1].response)

		}
	}))

	defer ts.Close()

	request := New()

	var waitForCompletion sync.WaitGroup

	for i := 0; i < 700; i++ {
		waitForCompletion.Add(2)

		go func() {
			if _, body, err := request.Post(ts.URL+"/json-req1").Set("snk1", "src1").Send(`{"name": "dan"}`).End(); err != nil {
				// 1. check error
				assert.NotNil(t, err)
				_ = body
			} else {

				// 2. check content of body
				assert.Equalf(t, `{"age": 28}`, (string)(body), "json request 1")
			}

			waitForCompletion.Done()

		}()

		go func() {
			if _, body, err := request.Post(ts.URL+"json-req2").Set("snk2", "src2").Send(`{"name": "lulu"}`).End(); err != nil {
				// 1. check error
				assert.NotNil(t, err)
				_ = body
			} else {

				// 2. check content of body
				assert.Equalf(t, `{"age": 30}`, (string)(body), "json request 2")
			}

			waitForCompletion.Done()
		}()
	}

	waitForCompletion.Wait()

}
