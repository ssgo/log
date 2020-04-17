package log

import (
	"bytes"
	"fmt"
	"github.com/ssgo/u"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type esWriter struct {
	config   *Config
	url      string
	user     string
	password string
	group    string
	lock     sync.Mutex
	queue    []string
	last     int64
	client   *http.Client
	prefix   string
}

func esWriterMaker(conf *Config) Writer {
	w := new(esWriter)
	w.config = conf
	w.url = conf.File
	w.queue = make([]string, 0)
	w.client = new(http.Client)
	esUrl, err := url.Parse(conf.File)
	if err != nil {
		DefaultLogger.Error(err.Error(), "conf", conf)
		return nil
	}

	if esUrl.User != nil {
		w.user = esUrl.User.Username()
		w.password, _ = esUrl.User.Password()
		esUrl.User = nil
	}
	if esUrl.Scheme == "ess" {
		esUrl.Scheme = "https"
	} else {
		esUrl.Scheme = "http"
	}

	timeout := esUrl.Query().Get("timeout")
	if timeout != "" {
		w.client.Timeout = u.Duration(timeout)
	}

	if len(esUrl.Path) > 1 {
		w.group = strings.ReplaceAll(esUrl.Path[1:], "/", ".")
	}

	esUrl.Path = "_bulk"
	esUrl.RawQuery = ""

	w.url = esUrl.String()

	if w.group != "" {
		w.prefix = fmt.Sprintf("{\"index\":{\"_index\":\"%s.%s\"}}", w.group, w.config.Name)
	} else {
		w.prefix = fmt.Sprintf("{\"index\":{\"_index\":\"%s\"}}", w.config.Name)
	}

	return w
}

func (w *esWriter) Log(data []byte) {
	l := len(data)
	if data == nil || l == 0 {
		return
	}
	dataString := string(data)

	// 将数据加入队列
	w.lock.Lock()
	w.queue = append(w.queue, w.prefix, dataString)
	w.lock.Unlock()
}

var responseOkString = []byte("\"errors\":false")

func (w *esWriter) Run() {
	now := time.Now().Unix()
	// 超过100条数据 或 过了1秒 发送数据（1秒内不超过100条不发送）
	if len(w.queue) > 100 || (len(w.queue) > 0 && (now > w.last || !writerRunning)) {
		var sendings []string
		w.lock.Lock()
		sendings = w.queue
		w.queue = make([]string, 0)
		w.lock.Unlock()

		data := strings.Join(sendings, "\n") + "\n"
		req, err := http.NewRequest("POST", w.url, bytes.NewReader([]byte(data)))
		if err != nil {
			log.Println("es sent failed", err.Error(), data)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		if w.user != "" {
			req.SetBasicAuth(w.user, w.password)
		}

		res, err := w.client.Do(req)
		if err != nil {
			log.Println("es sent failed", err.Error(), data)
			return
		}

		result, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("es sent failed", err.Error(), data)
			return
		}
		_ = res.Body.Close()

		if bytes.Index(result, responseOkString) == -1 {
			log.Println("es sent failed", string(result), data)
		}
		w.last = now
	}
}
