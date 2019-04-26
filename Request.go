package log

import (
	"github.com/ssgo/standard"
)

func (logger *Logger) Request(app, node, clientIp, fromApp, fromNode, clientId, sessionId, requestId, host, scheme, proto string, authLevel, priority int, method, path string, requestHeaders map[string]string, requestData map[string]interface{}, usedTime float32, responseCode int, responseHeaders map[string]string, responseDataLength uint, responseData interface{}, extra ...interface{}) {
	if !logger.checkLevel(INFO) {
		return
	}
	logger.log(standard.RequestLog{
		BaseLog:            logger.getBaseLog(standard.LogTypeRequest, extra...),
		App:                app,
		Node:               node,
		ClientIp:           clientIp,
		FromApp:            fromApp,
		FromNode:           fromNode,
		ClientId:           clientId,
		SessionId:          sessionId,
		RequestId:          requestId,
		Host:               host,
		Scheme:             scheme,
		Proto:              proto,
		AuthLevel:          authLevel,
		Priority:           priority,
		Method:             method,
		Path:               path,
		RequestHeaders:     requestHeaders,
		RequestData:        requestData,
		UsedTime:           usedTime,
		ResponseCode:       responseCode,
		ResponseHeaders:    responseHeaders,
		ResponseDataLength: responseDataLength,
		ResponseData:       responseData,
	})
}
