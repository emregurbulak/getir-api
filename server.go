package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	errors "github.com/mongodb-developer/docker-golang-example/Errors"
	"github.com/mongodb-developer/docker-golang-example/helper"
	models "github.com/mongodb-developer/docker-golang-example/model"
	mongodb "github.com/mongodb-developer/docker-golang-example/mongoDb"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var Cache = cache.New(5*time.Minute, 5*time.Minute)

func main() {
	relatedRecordsHandler := newRelatedRecordsHandler()
	inMemoryHandler := newInMemoryHandler()
	client, _ = mongodb.ConnectDbClient()
	http.HandleFunc("/getRelatedRecords", relatedRecordsHandler.relatedRecordMethods)
	http.HandleFunc("/in-memory", inMemoryHandler.inMemoryMethods)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func connectDbAndGetRecords() (processedRecords []models.ProcessedRecord, err error) {
	dbCollection := client.Database("getir-case-study").Collection("records")
	processedRecords, err = mongodb.GetProcessedRecords(dbCollection)
	return processedRecords, err
}

type relatedRecordsHandler struct {
	sync.Mutex
	store map[string]models.RelatedRecordsResponse
}

type inMemoryHandler struct {
	sync.Mutex
	store map[string]models.InMemoryResponse
}

func newRelatedRecordsHandler() *relatedRecordsHandler {
	return &relatedRecordsHandler{
		store: map[string]models.RelatedRecordsResponse{},
	}
}

func newInMemoryHandler() *inMemoryHandler {
	return &inMemoryHandler{
		store: map[string]models.InMemoryResponse{},
	}
}

func (h *relatedRecordsHandler) relatedRecordMethods(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleRelatedRecords(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func (h *inMemoryHandler) inMemoryMethods(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleGetInMemory(w, r)
		return
	case "POST":
		h.handlePostInMemory(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func (h *inMemoryHandler) handleGetInMemory(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	resp := models.InMemoryResponse{}
	data, found := helper.GetCache("inMemory")
	w.Header().Add("content-type", "application/json")
	if !found {
		jsonBytes, _ := json.Marshal(inMemoryFailureCondition(errors.KeyWasNotFound))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(jsonBytes))
		return
	}
	for _, v := range data {
		if v.Key == strings.Trim(key, `"`) {
			resp.Key = v.Key
			resp.Value = v.Value
		}
	}
	if resp.Key == "" {
		jsonBytes, _ := json.Marshal(inMemoryFailureCondition(errors.KeyWasNotFound))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(jsonBytes))
		return
	}

	jsonBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonBytes))
}

func (h *inMemoryHandler) handlePostInMemory(w http.ResponseWriter, r *http.Request) {
	var memo helper.InMemory
	err := json.NewDecoder(r.Body).Decode(&memo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	helper.SetCache("inMemory", memo)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("key and value posted"))
}

func (h *relatedRecordsHandler) handleRelatedRecords(w http.ResponseWriter, r *http.Request) {
	stDateStr := r.FormValue("startDate")
	endDateStr := r.FormValue("endDate")
	minCountStr := r.FormValue("minCount")
	maxCountStr := r.FormValue("maxCount")

	w.Header().Add("content-type", "application/json")

	minCount, err := strconv.Atoi(minCountStr)
	if err != nil {
		jsonBytes, _ := json.Marshal(failureCondition(errors.RequestParameterMinCountError))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(jsonBytes))
		return
	}

	maxCount, err := strconv.Atoi(maxCountStr)
	if err != nil {
		jsonBytes, _ := json.Marshal(failureCondition(errors.RequestParameterMaxCountError))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(jsonBytes))
		return
	}

	if minCount > maxCount {
		jsonBytes, _ := json.Marshal(failureCondition(errors.MinCountBiggerThanMaxCountError))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(jsonBytes))
		return
	}

	stDate, err := time.Parse("2006-01-02", strings.Trim(stDateStr, `\"`))
	if err != nil {
		jsonBytes, _ := json.Marshal(failureCondition(errors.RequestParameterDateTypeError))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(jsonBytes))
		return
	}
	endDate, err := time.Parse("2006-01-02", strings.Trim(endDateStr, `\"`))
	if err != nil {
		jsonBytes, _ := json.Marshal(failureCondition(errors.RequestParameterDateTypeError))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(jsonBytes))
		return
	}

	resp := getRelatedRecords(minCount, maxCount, stDate, endDate)
	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		jsonBytes, _ := json.Marshal(failureCondition(errors.ConnectingDbError))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(jsonBytes))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func getRelatedRecords(minCount, maxCount int, stDate, endDate time.Time) (resp models.RelatedRecordsResponse) {
	processedRecords, err := connectDbAndGetRecords()
	if err != nil {
		return failureCondition(err.Error())
	}
	//h.Lock()
	relatedRecord := models.RelatedRecord{}
	relatedRecords := []models.RelatedRecord{}
	for _, v := range processedRecords {
		resp.Code = 0
		resp.Message = "success"

		if stDate.Before(v.CreatedAt) && endDate.After(v.CreatedAt) &&
			minCount < int(v.TotalCount) && maxCount > int(v.TotalCount) {
			relatedRecord.CreatedAt = v.CreatedAt.Format("2006-01-02")
			relatedRecord.Key = v.Key
			relatedRecord.TotalCount = int(v.TotalCount)
			relatedRecords = append(relatedRecords, relatedRecord)
		}

	}
	resp.Records = []models.RelatedRecord{}
	resp.Records = append(resp.Records, relatedRecords...)
	//h.Unlock()
	return resp
}

func failureCondition(errDetails string) (resp models.RelatedRecordsResponse) {
	resp.Code = 1
	resp.Message = "failure"
	resp.Records = []models.RelatedRecord{}
	resp.ErrorDetails = errDetails
	return resp
}

func inMemoryFailureCondition(errDetails string) (resp models.InMemoryResponse) {
	resp.Message = errDetails
	return resp
}
