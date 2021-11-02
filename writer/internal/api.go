package internal

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)


type ReadHandler struct {
	client *redis.Client
	log  *zap.SugaredLogger
}

func NewReadHandler(client *redis.Client, log *zap.SugaredLogger) *ReadHandler {
	return &ReadHandler{client: client, log:log}
}


func (c *ReadHandler) ReadRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/v1/values/{index}", c.Write).Methods("POST")
	return router
}


func (c *ReadHandler) Write(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	index := params["index"]
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err !=nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResponse{
			Msg:  fmt.Errorf("error while reading the value %w", err),
			Code: 100001,
		})
		return
	}
	var request WriteRequest
	err = json.Unmarshal(body, &request)
	if err !=nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResponse{
			Msg:  fmt.Errorf("error while reading the value %w", err),
			Code: 100001,
		})
		return
	}

	// sets the index
	val, err := c.client.Set(r.Context(), index, request.Val, 0).Result()
	if err != nil {
		c.log.Error(err)
		w.WriteHeader(400)

		json.NewEncoder(w).Encode(ErrorResponse{
			Msg:  fmt.Errorf("error while reading the value %w", err),
			Code: 100003,
		})
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(ReadResponse{
		Val: val,
	})
}


type ErrorResponse struct {
	Msg   error
	Code  int
}

type ReadResponse struct {
	Val   interface{}
}


type WriteRequest struct {
	Val interface{}
}





