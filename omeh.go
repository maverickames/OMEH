package omeh

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type ErrorHandlerdesc func(http.ResponseWriter, *http.Request) *ErrResponse

type processError func(error, *ErrResponse, string) *ErrResponse

type handleErrorLog func(*ErrResponse)

type respondToError func(ErrorHandlerdesc) http.HandlerFunc

var StatusBadRequest = &ErrResponse{StatusCode: 400, StatusText: "Status Bad Request"}

var StatusNotFound = &ErrResponse{StatusCode: 404, StatusText: "Status Not Found"}

var StatusInternalServerError = &ErrResponse{StatusCode: 500, StatusText: "Status Internal Server Error"}

var YouDoneMessedUpAARon = &ErrResponse{StatusCode: 500, StatusText: "You done messed up, A-A-Ron"}

var NonUIError = &ErrResponse{StatusCode: 0, StatusText: "Server error!"}

type ErrResponse struct {
	// Internal Only
	Err           error   `json:"-"`
	RequestDetail string  `json:"-"`
	FuncPC        uintptr `json:"-"`
	FuncFN        string  `json:"-"`
	FuncLine      int     `json:"-"`

	// Public
	StatusCode int    `json:"-"`
	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
}

type ErrManager struct {
	debug            bool
	ErrorHandler     processError
	LogError         handleErrorLog
	HandleHTTPErrors respondToError
}

func New(d bool) *ErrManager {
	return &ErrManager{debug: d}
}

func (e *ErrManager) SetDefaultHandler(def processError) {
	e.ErrorHandler = def
}

func (e *ErrManager) SetErrorLogHandler(handler handleErrorLog) {
	e.LogError = handler
}
func (e *ErrManager) SetHTTPErrorHandler(handler respondToError) {
	e.HandleHTTPErrors = handler
}

func (e *ErrManager) SetDebug(d bool) {
	e.debug = d
}

func (e *ErrManager) ProcessErrorHTTP(h ErrorHandlerdesc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ErrResponse := h(w, r)
		if ErrResponse == nil {
			return
		}
		b, err := json.Marshal(ErrResponse)
		if err != nil {
			e.LogError(ErrResponse)
			return
		}
		if e.LogError != nil {
			e.LogError(ErrResponse)
		}

		w.WriteHeader(ErrResponse.StatusCode)
		w.Write([]byte(b))
	})
}

func returnErrorResponse(er *ErrResponse) *ErrResponse {
	if er != nil {
		return er
	}
	return &ErrResponse{}
}

func (e *ErrManager) getDebugDiagonisotics(err error) (uintptr, string, int) {
	var pc uintptr
	var fn string
	var line int
	if e.debug {
		pc, fn, line, _ = runtime.Caller(2)
	}
	return pc, fn, line
}

func (e *ErrManager) ReturnError(err error, er *ErrResponse, requestDetails string) *ErrResponse {

	errResponse := returnErrorResponse(er)
	pc, fn, line := e.getDebugDiagonisotics(err)

	errResponse.FuncPC = pc
	errResponse.FuncFN = fn
	errResponse.FuncLine = line

	errResponse.Err = err
	if errResponse.StatusCode == 0 {
		errResponse.StatusCode = 400
	}

	if requestDetails == "" {
		errResponse.RequestDetail = "Request was not passed. Most likely to protect the request data"
	} else {
		errResponse.RequestDetail = requestDetails
	}

	errResponse.AppCode = rand.New(rand.NewSource(time.Now().Unix())).Int63()

	if errResponse.StatusText == "" {
		errResponse.StatusText = "Internal Error"
	}
	return errResponse
}
