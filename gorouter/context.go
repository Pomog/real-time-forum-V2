package gorouter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Context struct is used to handle http requests and responses
type Context struct {
	http.ResponseWriter
	*http.Request
	Params map[string]string
}

type Error struct {
	Error string `json:"error"`
}

func (ctx *Context) WriteString(code int, body string) { //Func to write string as response
	ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
	ctx.WriteHeader(code)

	_, err := ctx.ResponseWriter.Write([]byte(body))
	if err != nil {
		ctx.WriteError(500, "ctx.WriteString error")
	}
}

func (ctx *Context) WriteJSON(code int, data interface{}) error { //Func to serialize data to JSON format
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	ctx.WriteHeader(code)

	_, err = ctx.ResponseWriter.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Context) WriteError(code int, err string) { //Func to write error message in JSON
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	ctx.WriteHeader(code)

	jsonData, _ := json.Marshal(&Error{Error: err})
	ctx.ResponseWriter.Write(jsonData)
}

func (ctx *Context) ReadBody(data interface{}) error { //Func to read request body
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Context) setURLValues(keys, values []string) { //Rest of the functions to handle URL parameters
	for i, key := range keys {
		ctx.SetParam(key, values[i])
	}
}

func (ctx *Context) SetParam(key string, value string) {
	ctx.Params[key] = value
}

func (ctx *Context) GetStringParam(key string) (string, error) {
	value, ok := ctx.Params[key]
	if !ok {
		return "", fmt.Errorf("%s value not found", key)
	}
	return value, nil
}

func (ctx *Context) GetIntParam(key string) (int, error) {
	value, ok := ctx.Params[key]
	if !ok {
		return 0, fmt.Errorf("%s value not found", key)
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s value must be integer", key)
	}

	return n, nil
}
