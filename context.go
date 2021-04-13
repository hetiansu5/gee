package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type H map[string]interface{}

type Context struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	method   string
	path     string
	params   map[string]string
	handlers []HandlerFunc
	index    int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: req,
		method:  strings.ToUpper(req.Method),
		path:    req.URL.Path,
		params:  make(map[string]string),
		index:   -1,
	}
}

func (c *Context) Next() {
	c.index++
	if c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.Next()
	}
}

func (c *Context) PostForm(key string) string {
	return c.Request.PostFormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Param(key string) string {
	return c.params[key]
}

func (c *Context) Page403() {
	c.String(http.StatusForbidden, "Access Forbidden")
}

func (c *Context) Page404() {
	c.String(http.StatusNotFound, "Not Found")
}

func (c *Context) Page405() {
	c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
}

func (c *Context) Page500() {
	c.String(http.StatusInternalServerError, "Internal Server Error")
}

func (c *Context) JSON(statusCode int, mp H) {
	bytes, err := json.Marshal(mp)
	if err != nil {
		c.String(500, "Internal Error")
		return
	}
	c.SetHeader("Content-Type", "application/json")
	c.Data(statusCode, bytes)
}

func (c *Context) HTML(statusCode int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Data(statusCode, []byte(html))
}

func (c *Context) String(statusCode int, format string, a ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Data(statusCode, []byte(fmt.Sprintf(format, a...)))
}

func (c *Context) Chunks(statusCode int, chunks []string) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetHeader("Transfer-Encoding", "chunked")
	c.Status(statusCode)
	for _, chunk := range chunks {
		if len(chunk) == 0 {
			continue
		}
		//c.writeString(fmt.Sprintf("%x\r\n", len(chunk)))
		c.writeString(fmt.Sprintf("%s\r\n", chunk))
		c.Flush()
		time.Sleep(time.Second)
	}
	//c.writeString("0\r\n\r\n")
	//c.Flush()
}

func (c *Context) Data(statusCode int, data []byte) {
	c.Status(statusCode)
	_, _ = c.Writer.Write(data)
}

func (c *Context) writeString(data string) {
	_, _ = c.Writer.Write([]byte(data))
}

func (c *Context) Status(statusCode int) {
	c.Writer.WriteHeader(statusCode)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) Flush() {
	c.Writer.(http.Flusher).Flush()
}
