package kiwi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func RestHandler(c Client) http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/", methodSwitch(c))
	return r
}

func methodSwitch(c Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGet(c)(w, r)
		case http.MethodPut, http.MethodPost:
			handlePut(c)(w, r)
		case http.MethodPatch:
			handlePatch(c)(w, r)
		case http.MethodDelete:
			handleDelete(c)(w, r)
		default:
			handleError(w, http.StatusBadRequest, errors.Errorf("method %q not allowed", r.Method))
		}
	}
}

func handleGet(c Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucket, key := parseBucketAndKey(w, r)
		if bucket == "" {
			handleListBuckets(c)(w, r)
			return
		}
		if key == "" {
			handleListBucketKeys(c, bucket)(w, r)
			return
		}

		var value ByteValue
		if err := c.Read(bucket, key, &value); err != nil {
			if IsNotFound(err) {
				handleError(w, http.StatusNotFound, err)
				return
			}
			handleError(w, http.StatusInternalServerError, err)
			return
		}

		data, err := value.MarshalBinary()
		if err != nil {
			handleEncodingResponseError(w, err)
			return
		}

		w.Write(data)
	}
}

func handlePut(c Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucket, key := parseBucketAndKey(w, r)
		if key == "" {
			handleRestStatus(w, restKeyMissingError)
			return
		}

		data, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			handleError(w, http.StatusBadRequest, errors.Wrap(err, "failed to read request body"))
			return
		}

		if err := c.Create(bucket, key, ByteValue(data)); err != nil {
			if IsNotFound(err) {
				handleError(w, http.StatusBadRequest, err)
			}
			handleError(w, http.StatusNotModified, err)
			return
		}

		handleRestStatus(w, restStatus{Status: http.StatusCreated, Msg: "success"})
	}
}

func handlePatch(c Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucket, key := parseBucketAndKey(w, r)
		if key == "" {
			handleRestStatus(w, restKeyMissingError)
			return
		}

		data, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			handleError(w, http.StatusBadRequest, errors.Wrap(err, "unable to read request body"))
			return
		}

		if err := c.Update(bucket, key, ByteValue(data)); err != nil {
			if IsNotFound(err) {
				handleError(w, http.StatusNotFound, err)
			}
			return
		}

		handleRestStatus(w, restStatus{Status: http.StatusAccepted, Msg: "success"})
	}
}

func handleDelete(c Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucket, key := parseBucketAndKey(w, r)
		if key == "" {
			handleRestStatus(w, restKeyMissingError)
			return
		}

		if err := c.Destroy(bucket, key); err != nil {
			if IsNotFound(err) {
				handleError(w, http.StatusNotFound, err)
				return
			}
			return
		}

		handleRestStatus(w, restStatus{Status: http.StatusOK, Msg: "success"})
	}
}

func handleListBuckets(c Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buckets, err := c.ListBuckets()
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}

		j := json.NewEncoder(w)
		if err := j.Encode(buckets); err != nil {
			return
		}
	}
}

func handleListBucketKeys(c Client, bucket string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, err := c.ListKeys(bucket)
		if err != nil {
			handleError(w, http.StatusNotFound, err)
			return
		}

		// response := map[string][]string{bucket: keys}
		j := json.NewEncoder(w)
		if err := j.Encode(keys); err != nil {
			return
		}
	}
}

func parseBucketAndKey(w http.ResponseWriter, r *http.Request) (b string, k string) {
	// url := r.URL.Path
	// idx := strings.LastIndex(url, "/")
	// k = url[idx+1:]
	// url = url[:idx]
	// idx = strings.LastIndex(url, "/")
	// b = url[idx+1:]
	// return

	b = r.URL.Query().Get("bucket")
	k = r.URL.Query().Get("key")

	return
}

type restStatus struct {
	Status int    `json:"status"`
	Msg    string `json:"message"`
}

var (
	restKeyMissingError = restStatus{
		Status: http.StatusBadRequest,
		Msg:    "unable to parse bucket and key from URL: missing key",
	}
)

func handleRestStatus(w http.ResponseWriter, s restStatus) {
	data, err := json.Marshal(s)
	if err != nil {
		http.Error(w, "internal encoding error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(s.Status)
	w.Write(data)
}

func handleError(w http.ResponseWriter, code int, e error) {
	handleRestStatus(w, restStatus{
		Status: code,
		Msg:    e.Error(),
	})
}

func handleEncodingResponseError(w http.ResponseWriter, e error) {
	handleRestStatus(w, restStatus{
		Status: http.StatusInternalServerError,
		Msg:    errors.Wrap(e, "failed to encode response").Error(),
	})
}
