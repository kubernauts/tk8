package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type restServer interface {
	Routes() []*Route
	String() string
	logRequest(request string, id string) *logrus.Entry
	sendError(request string, id string, w http.ResponseWriter, msg string, code int)
}

type restBase struct {
	restServer
	version string
	name    string
}

func StartTK8API(tk8ApiBase string, tk8Port uint16) error {
	tk8Api := newTk8Api()

	// start server as before
	_, p, err := startServer("tk8", tk8ApiBase, tk8Port, tk8Api)
	if err != nil {
		return err
	}
	logrus.Println(p.Addr)
	return nil
}

func startServer(name string, sockBase string, port uint16, rs restServer) (*http.Server, *http.Server, error) {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFound)
	for _, v := range rs.Routes() {
		logrus.Println(v.path)
		router.Methods(v.verb).Path(v.path).HandlerFunc(v.fn)
	}
	return startServerCommon(name, sockBase, port, rs, router)
}

func startServerCommon(name string, sockBase string, port uint16, rs restServer, router *mux.Router) (*http.Server, *http.Server, error) {
	// var (
	// 	listener net.Listener
	// 	err      error
	// )
	// socket := path.Join(sockBase, name+".sock")
	// os.Remove(socket)
	// os.MkdirAll(path.Dir(socket), 0755)

	// logrus.Printf("Starting REST service on socket : %+v", socket)
	// listener, err = net.Listen("unix", socket)
	// if err != nil {
	// 	logrus.Warnln("Cannot listen on UNIX socket: ", err)
	// 	return nil, nil, err
	// }
	// unixServer := &http.Server{Handler: router}
	// go unixServer.Serve(listener)

	//if port != 0 {
	logrus.Printf("Starting REST service on port : %v", port)
	portServer := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: router}
	go portServer.ListenAndServe()
	return nil, portServer, nil
	//	}
	//return unixServer, nil, nil
}

func (rest *restBase) logRequest(request string, id string) *logrus.Entry {
	return logrus.WithFields(map[string]interface{}{
		"Driver":  rest.name,
		"Request": request,
		"ID":      id,
	})
}
func (rest *restBase) sendError(request string, id string, w http.ResponseWriter, msg string, code int) {
	rest.logRequest(request, id).Warnln(code, " ", msg)
	http.Error(w, msg, code)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	logrus.Warnf("Not found: %+v ", r.URL)
	http.NotFound(w, r)
}
