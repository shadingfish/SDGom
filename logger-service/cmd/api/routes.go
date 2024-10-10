package main

// 这段代码是用于配置HTTP路由的，用到了go-chi/chi框架，这是一个轻量级的Go语言路由库，非常适合用于创建RESTful API。在这个函数中，定义了允许的跨域（CORS）策略、加入了中间件（middleware），并指定了一个用于日志记录的POST请求路径（/log）

import (
	"github.com/go-chi/chi/v5"            // 引入go-chi路由库，用于管理路由
	"github.com/go-chi/chi/v5/middleware" // 引入go-chi的中间件库
	"github.com/go-chi/cors"              // 引入go-chi的CORS处理库，用于处理跨域请求
	"net/http"                            // 引入用于处理HTTP请求和响应的标准库
)

// 这里导入了net/http标准库，以及chi框架和它的中间件库。chi用于路由管理，middleware用于增加路由中间件，cors用于处理跨域资源共享（CORS）。
// This imports the net/http standard library, as well as the chi framework and its middleware library. chi is used for route management, middleware for adding route middleware, and cors for handling Cross-Origin Resource Sharing (CORS).

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// 定义了一个接收者函数routes，该函数属于Config结构体。返回值类型是http.Handler，表示可以处理HTTP请求的对象。
	// This defines a receiver function routes that belongs to the Config struct. The return type is http.Handler, which indicates that it can handle HTTP requests.

	// mux := chi.NewRouter()：创建了一个新的chi路由器，用于定义和管理HTTP路由。
	// mux := chi.NewRouter() creates a new chi router for defining and managing HTTP routes.

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},                                   // 允许所有HTTP和HTTPS协议的跨域请求
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                 // 允许的HTTP方法
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // 允许的请求头
		ExposedHeaders:   []string{"Link"},                                                    // 允许暴露的响应头
		AllowCredentials: true,                                                                // 允许客户端发送包含身份验证信息的请求（如cookie）
		MaxAge:           300,                                                                 // 设置预检请求的缓存时间（单位为秒）
	}))

	// mux.Use(cors.Handler(cors.Options{...}))：为mux路由器添加了CORS中间件。
	// This line adds a CORS middleware to the mux router.

	// AllowedOrigins: []string{"https://*", "http://*"}：允许所有HTTP和HTTPS协议的跨域请求。
	// Allows all HTTP and HTTPS protocol cross-origin requests.

	// AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}：允许的HTTP方法包括GET、POST、PUT、DELETE、OPTIONS。
	// Permits the HTTP methods GET, POST, PUT, DELETE, and OPTIONS.

	// AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}：允许的请求头字段有Accept、Authorization、Content-Type、X-CSRF-Token。
	// Permits request headers such as Accept, Authorization, Content-Type, X-CSRF-Token.

	// ExposedHeaders: []string{"Link"}：指定哪些响应头可以暴露给客户端（允许客户端访问）。
	// Specifies which response headers can be exposed to the client (allowed to be accessed by the client).

	// AllowCredentials: true：表示是否允许客户端发送带有身份验证信息的请求（如cookie）。
	// Indicates whether the client is allowed to send requests with authentication information (e.g., cookies).

	// MaxAge: 300：设置预检请求的缓存时间（以秒为单位）。浏览器会在这段时间内缓存该CORS策略，不会再次发送预检请求。
	// Sets the cache time for preflight requests (in seconds). The browser will cache this CORS policy during this time and will not send a preflight request again.

	mux.Use(middleware.Heartbeat("/ping")) // 添加一个心跳中间件，用于检测服务是否正常运行
	// mux.Use(middleware.Heartbeat("/ping"))：添加了一个心跳检测中间件，它会在/ping路径上返回一个200状态码的响应，表示服务正常。
	// Adds a heartbeat middleware, which returns a 200 status code response at the /ping path, indicating that the service is running normally.

	mux.Post("/log", app.WriteLog) // 定义一个POST请求，路径是`/log`，请求处理函数是`app.WriteLog`
	// mux.Post("/log", app.WriteLog)：定义了一个POST请求，路径是/log，处理函数是 app.WriteLog。这意味着当客户端发送一个POST请求到 /log 时，会调用 app.WriteLog 函数处理该请求。
	// Defines a POST request with the path /log and the handler function app.WriteLog. This means when a client sends a POST request to /log, the app.WriteLog function will handle it.

	return mux
	// return mux：返回配置好的mux路由器实例，作为HTTP处理程序供外部使用。
	// Returns the configured mux router instance as an HTTP handler for external use.
}

// 重点解释（Key Explanations）
// Chi路由库（Chi Router Library）:

// chi 是一个轻量级但功能强大的Go语言路由库，它提供了类似 net/http 标准库的接口，并且可以方便地与中间件集成。
// chi is a lightweight yet powerful Go router library that provides interfaces similar to the net/http standard library and can be easily integrated with middleware.
// CORS（跨域资源共享，Cross-Origin Resource Sharing）:

// CORS是一种安全机制，用于允许或拒绝跨域请求。mux.Use(cors.Handler(cors.Options{...})) 定义了允许的源、方法、头和其他选项，以确保客户端能够跨域访问资源。
// CORS is a security mechanism used to allow or deny cross-origin requests. mux.Use(cors.Handler(cors.Options{...})) defines the allowed origins, methods, headers, and other options to ensure clients can access resources across origins.

// 中间件（Middleware）:

// 中间件是处理请求的拦截器，可以在请求被路由到具体的处理器之前或响应返回客户端之前对请求或响应进行处理。middleware.Heartbeat("/ping") 是一个示例，它用来检测服务的健康状态。
// Middleware are interceptors that process requests before they are routed to specific handlers or before responses are returned to clients. middleware.Heartbeat("/ping") is an example used to check the health status of the service.

// 总结（Summary）
// 这段代码使用了 chi 路由库定义了一个HTTP路由器，并配置了CORS策略、中间件，以及一个处理 /log POST请求的路由。
// This code uses the chi router library to define an HTTP router, configure CORS policies, middleware, and a route to handle /log POST requests.

// cors 处理跨域策略，使得客户端能够根据定义好的策略访问服务器资源。
// cors handles cross-origin policies, allowing clients to access server resources based on the defined policy.

// 中间件 middleware.Heartbeat("/ping") 添加了一个健康检查功能，使得客户端能够轻松检测服务器的运行状态。
// The middleware middleware.Heartbeat("/ping") adds a health check feature, enabling clients to easily check the server's running status.

// 跨域请求（Cross-Origin Request）是指在一个域名（origin）下的网页向不同域名的服务器发起的请求。简单来说，跨域请求发生在以下场景中：

// 协议不同（Protocol Difference）：例如，http://example.com 和 https://example.com 属于不同的域。
// 域名不同（Domain Difference）：例如，http://example.com 和 http://anotherdomain.com。
// 端口不同（Port Difference）：例如，http://example.com:8080 和 http://example.com:3000。
// 举例说明：
// 假设你有一个网页，托管在 http://example.com，当该网页尝试向 http://api.anotherdomain.com 发起请求时，这个请求就属于跨域请求。由于出于安全原因，浏览器默认会阻止跨域请求，从而防止恶意网站窃取用户数据。

// 为什么会有跨域问题？
// 跨域问题源自浏览器的同源策略（Same-Origin Policy）。同源策略是一种重要的安全机制，用于防止恶意脚本在没有授权的情况下读取另一个网站的敏感数据。根据同源策略，浏览器允许在同源（协议、域名、端口都相同）的情况下进行数据交互，而对于不同源的请求，则默认会被阻止。

// 同源策略的规则（Same-Origin Policy Rules）
// 要满足同源策略，以下三个条件必须都相同：

// 协议（Protocol）：例如，http 和 https 不同。
// 域名（Domain）：例如，example.com 和 sub.example.com 不同。
// 端口（Port）：例如，80 和 8080 不同。
// 如果任何一个条件不满足，则认为是跨域请求，浏览器会根据CORS策略决定是否允许该请求。

// CORS（跨域资源共享，Cross-Origin Resource Sharing）
// 为了安全地处理跨域请求，浏览器和服务器可以使用CORS（Cross-Origin Resource Sharing）协议来进行跨域资源共享。CORS是一种机制，它允许服务器通过设置HTTP头，告知浏览器该资源是否可以被其他域名访问。

// 例如，一个跨域请求的CORS响应头可能如下：

// http
// Access-Control-Allow-Origin: https://anotherdomain.com
// Access-Control-Allow-Methods: GET, POST, PUT, DELETE
// Access-Control-Allow-Headers: Content-Type, Authorization
// 这些头信息告诉浏览器：

// Access-Control-Allow-Origin：允许来自 https://anotherdomain.com 的请求。
// Access-Control-Allow-Methods：允许的HTTP方法包括 GET、POST、PUT 和 DELETE。
// Access-Control-Allow-Headers：允许的请求头包括 Content-Type 和 Authorization。
// CORS请求的类型
// 简单请求（Simple Request）：浏览器认为是简单请求的跨域请求，会自动处理并发送给目标服务器。一般满足以下条件的请求被认为是简单请求：
// 使用的HTTP方法是 GET、POST 或 HEAD。
// 请求头中包含的字段是：Accept、Accept-Language、Content-Language、Content-Type（值仅限于 application/x-www-form-urlencoded、multipart/form-data、text/plain）。
// 预检请求（Preflight Request）：当跨域请求不满足简单请求的条件时，浏览器会在实际请求之前，先发送一个 OPTIONS 请求到服务器。这就是所谓的预检请求。服务器通过这个请求来告知浏览器是否允许该跨域请求。这种机制可以避免某些潜在的安全风险。
// 解决跨域请求的方法
// 使用CORS：

// 服务器端设置 Access-Control-Allow-Origin 等响应头信息，以允许特定域名的请求。
// 例如，在Node.js中使用 cors 中间件：
// javascript
// Copy code
// const cors = require('cors');
// app.use(cors({ origin: 'https://anotherdomain.com' }));
// JSONP（JSON with Padding）：

// JSONP是一种历史上用于跨域请求的解决方案，通过动态创建 script 标签来实现跨域。虽然可以解决GET请求的跨域问题，但存在安全风险，且不支持POST请求，现已较少使用。
// 反向代理（Reverse Proxy）：

// 可以在前端服务器上配置反向代理，将跨域请求转发到同源地址，从而实现跨域数据请求。比如在Nginx中配置反向代理，将所有跨域请求指向目标服务器。
// nginx
// Copy code
// location /api/ {
//     proxy_pass http://targetserver.com/api/;
// }
// 在同一个域名下部署：

// 如果可能，将前端和后端应用部署在同一个域名或子域下，从而避免跨域问题。
// 总结
// 跨域请求是一种在不同源之间进行数据传输的方式，但由于安全策略（同源策略），通常会被浏览器阻止。为了解决这个问题，可以使用CORS协议、反向代理或其他技术来实现安全的跨域数据共享。
