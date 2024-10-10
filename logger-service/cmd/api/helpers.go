package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// 定义了一个jsonResponse结构体，用于表示标准的JSON响应格式。
type jsonResponse struct {
	Error   bool   `json:"error"`          // 是否发生错误
	Message string `json:"message"`        // 响应消息
	Data    any    `json:"data,omitempty"` // 可选的数据字段
}

// jsonResponse 是一个标准的JSON响应格式的结构体，用于在返回给客户端时提供一致的格式。
// jsonResponse is a struct for a standard JSON response format, used to provide a consistent format when returning to the client.

// Error 表示是否发生错误，Message 是响应的消息，Data 字段是可选的，可以包含任意类型的数据。
// Error indicates whether an error occurred, Message is the response message, and the Data field is optional and can contain any type of data.

// readJSON tries to read the body of a request and converts it into JSON
// readJSON 尝试读取请求体并将其转换为JSON格式
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 设置最大允许读取的请求体大小为1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body) // 创建JSON解码器
	err := dec.Decode(data)        // 将请求体解码为data指向的结构体
	if err != nil {
		return err
	}

	// 检查是否只有一个JSON对象
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value") // 如果解码第二次没有遇到EOF，说明存在多余的JSON数据
	}

	return nil
}

// readJSON() 函数用于读取HTTP请求体，并将其解码为JSON格式的数据。
// The readJSON() function reads the HTTP request body and decodes it into JSON format.

// maxBytes := 1048576 设置了读取请求体的最大字节数为1MB。
// maxBytes := 1048576 sets the maximum number of bytes to read from the request body to 1MB.

// dec := json.NewDecoder(r.Body) 创建了一个JSON解码器来读取请求体内容。
// dec := json.NewDecoder(r.Body) creates a JSON decoder to read the request body content.

// err = dec.Decode(data) 将请求体解码为data结构体。如果发生错误，则返回错误信息。
// err = dec.Decode(data) decodes the request body into the data struct. If an error occurs, it returns the error message.

// err = dec.Decode(&struct{}{}) 检查请求体是否仅包含一个JSON对象。如果解码第二次时没有遇到EOF（表示文件结束），则说明存在多余的JSON数据。
// err = dec.Decode(&struct{}{}) checks if the request body only contains a single JSON object. If it does not encounter EOF on the second decode, it indicates that there is extra JSON data.

// writeJSON takes a response status code and arbitrary data and writes a json response to the client
// writeJSON 接收响应状态码和任意数据，并将其写入JSON格式响应中返回给客户端
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data) // 将数据编码为JSON格式
	if err != nil {
		return err
	}

	// 如果提供了额外的响应头，则将其添加到响应中
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json") // 设置响应头为JSON类型
	w.WriteHeader(status)                              // 设置响应状态码
	_, err = w.Write(out)                              // 写入响应体
	if err != nil {
		return err
	}

	return nil
}

// writeJSON() 函数接收一个响应状态码和任意数据，并将其转换为JSON格式返回给客户端。
// The writeJSON() function takes a response status code and arbitrary data, converts it to JSON format, and returns it to the client.

// out, err := json.Marshal(data) 将数据编码为JSON格式。如果编码失败，则返回错误信息。
// out, err := json.Marshal(data) encodes the data into JSON format. If encoding fails, it returns an error.

// w.Header().Set("Content-Type", "application/json") 设置响应头类型为JSON。
// w.Header().Set("Content-Type", "application/json") sets the response header type to JSON.

// w.WriteHeader(status) 设置响应状态码。
// w.WriteHeader(status) sets the response status code.

// errorJSON takes an error, and optionally a response status code, and generates and sends
// a json error response
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}
