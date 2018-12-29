package actions

import (
	"fmt"
	"net/http"
	"strings"

	"golang-fave/engine/sessions"
)

type hRun func(e *Action)

type Action struct {
	W         *http.ResponseWriter
	R         *http.Request
	VHost     string
	VHostHome string
	RemoteIp  string
	Session   *sessions.Session
	list      map[string]hRun
}

func (e *Action) register(name string, handle hRun) {
	e.list[name] = handle
}

func (e *Action) write(data string) {
	(*e.W).Write([]byte(data))
}

func (e *Action) msg_show(title string, msg string) {
	e.write(fmt.Sprintf(
		`ModalShowMsg('%s', '%s');`,
		strings.Replace(strings.Replace(title, `'`, `&rsquo;`, -1), `"`, `&rdquo;`, -1),
		strings.Replace(strings.Replace(msg, `'`, `&rsquo;`, -1), `"`, `&rdquo;`, -1)))
}

func (e *Action) msg_success(msg string) {
	e.msg_show("Success", msg)
}

func (e *Action) msg_error(msg string) {
	e.msg_show("Error", msg)
}

func New(w *http.ResponseWriter, r *http.Request, vhost string, vhosthome string, remoteip string, session *sessions.Session) *Action {
	act := Action{w, r, vhost, vhosthome, remoteip, session, make(map[string]hRun)}

	// Register all action here
	act.register("mysql", action_mysql)
	act.register("signin", action_signin)

	return &act
}

func (e *Action) Call() bool {
	if e.R.Method != "POST" {
		return false
	}
	if err := e.R.ParseForm(); err == nil {
		action := e.R.FormValue("action")
		if action != "" {
			fn, ok := e.list[action]
			if ok {
				(*e.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				(*e.W).Header().Set("Content-Type", "text/html; charset=utf-8")
				fn(e)
				return true
			}
		}
	}
	return false
}
