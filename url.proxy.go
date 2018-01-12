package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func writeJsonp(callback string, output []byte) []byte {
	if callback != "" {
		var buf bytes.Buffer
		buf.WriteString(callback)
		buf.Write(output)
		buf.WriteByte(')')
		return buf.Bytes()
	}
	return output
}

func UrlProxyHandler(c *gin.Context) {
	var (
		callback    = c.Query("callback")
		contentType string
		urlResponse = struct{}{}
	)

	if callback == "" {
		contentType = "application/json; charset=utf-8"
	} else {
		contentType = "application/x-javascript; charset=utf-8"
	}

	if callback != "" {
		output, err := json.Marshal(urlResponse)
		if err != nil {
			logger.Error(err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, contentType, writeJsonp(callback, output))
	} else {
		c.JSON(http.StatusOK, urlResponse)
	}
	return
}
