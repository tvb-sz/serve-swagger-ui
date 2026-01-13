package route

import (
	"bytes"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tvb-sz/serve-swagger-ui/app/service"
	"github.com/tvb-sz/serve-swagger-ui/client"
	"github.com/tvb-sz/serve-swagger-ui/conf"
)

const (
	XRequestID          = "x-request-id"       // 请求ID名称
	XRequestIDPrefix    = "R"                  // 当使用纳秒时间戳作为请求ID时拼接的前缀字符串
	TextGinRouteInit    = "gin.route.init"     // gin 路由注册日志标记
	TextGinPanic        = "gin.panic.recovery" // gin panic日志标记
	TextGinRequest      = "gin.request"        // gin request请求日志标记
	TextGinResponseFail = "gin.response.fail"  // gin 业务层面失败响应日志标记
	TextGinPreflight    = "gin.preflight"      // gin preflight 预检options请求类型日志
)

// router 包内路由变量，请勿覆盖
//   - 一般扩展路由是基于该变量链式添加，为了识别可将固定前缀的路由拆分文件
var router *gin.Engine

// iniRoute 路由init-logger、recovery、cors 等
func iniRoute() {
	router = gin.New()

	// set base middleware
	router.Use(ginLogger(appendEmailIfExist), ginRecovery)
	if conf.Config.Server.Cors {
		router.Use(ginCors)
	}

	// 请求找不到路由时输出错误
	router.NoRoute(notRoute)
}

// appendEmailIfExist append email field to logger if exist
func appendEmailIfExist(ctx *gin.Context) []slog.Attr {
	if tokenInter, exist := ctx.Get("token"); exist {
		if token, ok := tokenInter.(service.Token); ok && token.Authenticated {
			filed := make([]slog.Attr, 0)
			return append(filed, slog.String("email", token.Email))
		}
	}

	return nil
}

// Bootstrap 引导初始化路由route
func Bootstrap() *gin.Engine {
	iniRoute()
	routeSetting()
	return router
}

// setRequestID 内部方法设置请求ID
func setRequestID(ctx *gin.Context) string {
	requestID := ctx.GetHeader(XRequestID)
	if requestID == "" {
		requestID = XRequestIDPrefix + strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	ctx.Set(XRequestID, requestID)
	return requestID
}

// GetRequestID 暴露方法：读取当前请求ID
func GetRequestID(ctx *gin.Context) string {
	if reqId, exist := ctx.Get(XRequestID); exist {
		return reqId.(string)
	}
	return ""
}

func ginLogger(appendHandle func(ctx *gin.Context) []slog.Attr) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		start := time.Now()

		// set XRequestID
		requestID := setRequestID(ctx)

		// +++++++++++++++++++++++++
		// 记录请求 body 体
		// Notice: http包里对*http.Request.Body这个Io是一次性读取，此处读取完需再次设置Body以便其他位置能顺利读取到参数内容
		// +++++++++++++++++++++++++
		bodyData := getRequestBody(ctx, true)

		// executes at end
		ctx.Next()

		latencyTime := time.Now().Sub(start)
		fields := []slog.Attr{
			slog.String("module", TextGinRequest),
			slog.String("ua", ctx.GetHeader("User-Agent")),
			slog.String("method", ctx.Request.Method),
			slog.String("req_id", requestID),
			slog.String("req_body", bodyData),
			slog.String("client_ip", ctx.ClientIP()),
			slog.String("url_path", ctx.Request.URL.Path),
			slog.String("url_query", ctx.Request.URL.RawQuery),
			slog.String("url", ctx.Request.URL.String()),
			slog.Int("http_status", ctx.Writer.Status()),
			slog.Duration("duration", latencyTime),
		}

		// 额外自定义补充字段
		if appendHandle != nil {
			fields = append(fields, appendHandle(ctx)...)
		}

		if latencyTime.Seconds() > 0.5 {
			slog.LogAttrs(ctx, slog.LevelWarn, ctx.Request.URL.Path, fields...)
		} else {
			slog.LogAttrs(ctx, slog.LevelInfo, ctx.Request.URL.Path, fields...)
		}
	}
}

func ginRecovery(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// dump出http请求相关信息
			httpRequest, _ := httputil.DumpRequest(ctx.Request, false)

			// 检查是否为tcp管道中断错误：这样就没办法给客户端通知消息
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}

			// record log
			client.Logger.Error(
				"http request failed",
				"url", ctx.Request.URL.Path,
				"request", string(httpRequest),
				"error", err,
			)

			if brokenPipe {
				// tcp中断导致panic，终止无输出
				_ = ctx.Error(err.(error))
				ctx.Abort()
			} else {
				// 非tcp中断导致panic，响应500错误
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
		}
	}()
	ctx.Next()
}

// ginCors 为gin开启跨域功能<尽量通过nginx反代处理>
func ginCors(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")
	ctx.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,X-App-Client,X-Requested-With,Authorization")
	ctx.Header("Access-Control-Allow-Methods", "GET,OPTIONS,POST,PUT,DELETE,PATCH")
	// REF https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Access-Control-Max-Age
	ctx.Header("Access-Control-Max-Age", "7200")
	if ctx.Request.Method == http.MethodOptions {
		client.Logger.Debug(
			TextGinPreflight,
			"module", TextGinPreflight,
			"ua", ctx.GetHeader("User-Agent"),
			"method", ctx.Request.Method,
			"req_id", GetRequestID(ctx),
			"client_ip", ctx.ClientIP(),
			"url_path", ctx.Request.URL.Path,
			"url_query", ctx.Request.URL.RawQuery,
			"url", ctx.Request.URL.String(),
		)
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}

// getRequestBody 获取请求body体
//   - strip 是否要将JSON类型的body体去除反斜杠和大括号，以便于Es等不做深层字段解析而当做一个字符串
func getRequestBody(ctx *gin.Context, strip bool) string {
	bodyData := ""

	// post、put、patch、delete四种类型请求记录body体
	if isModifyMethod(ctx.Request.Method) {
		// 判断是否为JSON实体类型<application/json>，仅需要判断content-type包含/json字符串即可
		if strings.Contains(ctx.ContentType(), "/json") {
			buf, _ := io.ReadAll(ctx.Request.Body)
			bodyData = string(buf)
			_ = ctx.Request.Body.Close()
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(buf)) // 重要

			// strip json `\{}` to ignore transfer JSON struct
			if strip {
				bodyData = strings.Replace(bodyData, "\\", "", -1)
				bodyData = strings.Replace(bodyData, "{", "", -1)
				bodyData = strings.Replace(bodyData, "}", "", -1)
			}
		} else {
			_ = ctx.Request.ParseForm() // 尝试解析表单, 文件表单忽略
			bodyData = ctx.Request.PostForm.Encode()
		}
	}

	return bodyData
}

// isModifyMethod 检查当前请求方式否为修改类型
//   - 即判断请求是否为post、put、patch、delete
func isModifyMethod(method string) bool {
	return method == http.MethodPost ||
		method == http.MethodPut ||
		method == http.MethodPatch ||
		method == http.MethodDelete
}
