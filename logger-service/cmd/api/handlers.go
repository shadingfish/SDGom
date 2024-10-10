package main

// 这段代码展示了一个用于记录日志的HTTP处理程序，以及几个用于读取和写入JSON响应的辅助函数。代码使用了Go语言的标准库，如encoding/json 和 net/http，以及一些自定义类型和函数。下面我们分段对这两个文件进行中英文详细解析。

import (
	"log-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// 这里定义了一个JSONPayload结构体，它用来表示从HTTP请求中接收的JSON数据。
// This defines a JSONPayload struct, which represents the JSON data received from the HTTP request.

// Name 和 Data 是两个字符串字段，分别用于接收请求JSON中name和data的值。
// Name and Data are two string fields that are used to receive the values of name and data in the request JSON.

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	// 这是一个属于Config结构体的HTTP处理函数WriteLog，用于处理/log路径的POST请求。
	// This is an HTTP handler function WriteLog that belongs to the Config struct, and it handles POST requests at the /log path.

	// 使用app.readJSON()函数将请求体中的JSON数据读取到requestPayload变量中。
	// It uses the app.readJSON() function to read the JSON data from the request body into the requestPayload variable.

	// insert data
	// 将JSON数据转换为LogEntry数据，并插入数据库中
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	// 调用Models中的LogEntry的Insert方法将数据插入到数据库中
	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// event是一个data.LogEntry类型的变量，它将接收到的requestPayload数据转换为日志条目（Log Entry）。
	// event is a variable of type data.LogEntry that converts the received requestPayload data into a log entry.

	// app.Models.LogEntry.Insert(event) 是一个数据库插入操作，将日志条目插入到数据库中。
	// app.Models.LogEntry.Insert(event) is a database insertion operation that inserts the log entry into the database.

	// 如果插入数据库时发生错误，则调用app.errorJSON(w, err)，返回一个JSON格式的错误响应，并终止程序。
	// If an error occurs during database insertion, it calls app.errorJSON(w, err) to return a JSON-formatted error response and terminates the function.

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	// 如果日志条目成功插入数据库中，则创建一个jsonResponse类型的响应，表示日志记录成功，并使用app.writeJSON()将响应返回给客户端。
	// If the log entry is successfully inserted into the database, it creates a jsonResponse response indicating the log was successfully recorded, and uses app.writeJSON() to send the response back to the client.

	app.writeJSON(w, http.StatusAccepted, resp)
}
