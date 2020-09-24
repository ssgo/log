package log

import (
	"github.com/ssgo/standard"
)

func (logger *Logger) Request(serverId, app, node, clientIp, fromApp, fromNode, clientId, deviceId, clientAppName, clientAppVersion, sessionId, requestId, host, scheme, proto string, authLevel, priority int, method, path string, requestHeaders map[string]string, requestData map[string]interface{}, usedTime float32, responseCode int, responseHeaders map[string]string, responseDataLength uint, responseData string, extra ...interface{}) {
	if !logger.CheckLevel(INFO) {
		return
	}
	logger.Log(logger.MakeRequestLog(standard.LogTypeRequest, serverId, app, node, clientIp, fromApp, fromNode, clientId, deviceId, clientAppName, clientAppVersion, sessionId, requestId, host, scheme, proto, authLevel, priority, method, path, requestHeaders, requestData, usedTime, responseCode, responseHeaders, responseDataLength, responseData, extra...))
}

func (logger *Logger) MakeRequestLog(logType, serverId, app, node, clientIp, fromApp, fromNode, clientId, deviceId, clientAppName, clientAppVersion, sessionId, requestId, host, scheme, proto string, authLevel, priority int, method, path string, requestHeaders map[string]string, requestData map[string]interface{}, usedTime float32, responseCode int, responseHeaders map[string]string, responseDataLength uint, responseData string, extra ...interface{}) standard.RequestLog {
	return standard.RequestLog{
		BaseLog:            logger.MakeBaseLog(logType, extra...),
		ServerId:           serverId,
		App:                app,
		Node:               node,
		ClientIp:           clientIp,
		FromApp:            fromApp,
		FromNode:           fromNode,
		ClientId:           clientId,
		DeviceId:           deviceId,
		ClientAppName:      clientAppName,
		ClientAppVersion:   clientAppVersion,
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
	}
}
