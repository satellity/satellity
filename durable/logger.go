package durable

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/logging"
)

type LoggerClient struct {
	impl *logging.Client
}

type Logger struct {
	impl    *logging.Logger
	request *logging.HTTPRequest
}

func NewLoggerClient(project string, syslog bool) (*LoggerClient, error) {
	if syslog {
		return &LoggerClient{}, nil
	}
	client, err := logging.NewClient(context.Background(), project)
	if err != nil {
		return nil, err
	}
	return &LoggerClient{impl: client}, nil
}

func (client *LoggerClient) Close() error {
	if client.impl == nil {
		return nil
	}
	return client.impl.Close()
}

func BuildLogger(client *LoggerClient, name string, r *http.Request) *Logger {
	if client.impl == nil {
		return &Logger{}
	}
	if r != nil {
		labels := logging.CommonLabels(map[string]string{
			"request-id": r.Header.Get("X-Request-Id"),
		})
		return &Logger{
			impl: client.impl.Logger(name, labels),
			request: &logging.HTTPRequest{
				Request:  r,
				RemoteIP: r.RemoteAddr,
			},
		}
	}
	return &Logger{
		impl: client.impl.Logger(name),
	}
}

func (logger *Logger) FillResponse(status int, responseSize int64, latency time.Duration) {
	if logger.request != nil {
		logger.request.Status = status
		logger.request.ResponseSize = responseSize
		logger.request.Latency = latency
	}
}

func (logger *Logger) Debug(v ...interface{}) {
	if logger.impl == nil {
		log.Println(v...)
		return
	}
	logger.impl.Log(logging.Entry{
		Severity:    logging.Debug,
		HTTPRequest: logger.request,
		Payload:     fmt.Sprint(v),
	})
}

func (logger *Logger) Debugf(format string, v ...interface{}) {
	if logger.impl == nil {
		log.Printf(format, v...)
		return
	}
	logger.impl.Log(logging.Entry{
		Severity:    logging.Debug,
		HTTPRequest: logger.request,
		Payload:     fmt.Sprintf(format, v...),
	})
}

func (logger *Logger) Info(v ...interface{}) {
	if logger.impl == nil {
		log.Println(v...)
		return
	}
	logger.impl.Log(logging.Entry{
		Severity:    logging.Info,
		HTTPRequest: logger.request,
		Payload:     fmt.Sprint(v),
	})
}

func (logger *Logger) Infof(format string, v ...interface{}) {
	if logger.impl == nil {
		log.Printf(format, v...)
		return
	}
	logger.impl.Log(logging.Entry{
		Severity:    logging.Info,
		HTTPRequest: logger.request,
		Payload:     fmt.Sprintf(format, v...),
	})
}

func (logger *Logger) Error(v ...interface{}) {
	if logger.impl == nil {
		log.Println(v...)
		return
	}
	logger.impl.Log(logging.Entry{
		Severity:    logging.Error,
		HTTPRequest: logger.request,
		Payload:     fmt.Sprint(v),
	})
}

func (logger *Logger) Errorf(format string, v ...interface{}) {
	if logger.impl == nil {
		log.Printf(format, v...)
		return
	}
	logger.impl.Log(logging.Entry{
		Severity:    logging.Error,
		HTTPRequest: logger.request,
		Payload:     fmt.Sprintf(format, v...),
	})
}

func (logger *Logger) Panicln(v ...interface{}) {
	if logger.impl == nil {
		log.Panicln(v)
		return
	}
	logger.impl.Log(logging.Entry{
		Severity:    logging.Critical,
		HTTPRequest: logger.request,
		Payload:     fmt.Sprint(v),
	})
	log.Panicln(v)
}
