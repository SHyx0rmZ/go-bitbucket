package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bitbucket Server Suite")
}

//
//var (
//	testMux    *http.ServeMux
//	testServer *httptest.Server
//)
//
//var _ = BeforeEach(func() {
//	testMux = http.NewServeMux()
//	testServer = httptest.NewServer(testMux)
//})
//
//var _ = AfterEach(func() {
//	testServer.Close()
//})
//
//type ResettableServeMux struct {
//	mux   *http.ServeMux
//	mutex sync.RWMutex
//}
//
//func (mux *ResettableServeMux) Handle(pattern string, handler http.Handler) {
//	mux.mutex.RLock()
//
//	mux.mux.Handle(pattern, handler)
//}
//
//func (mux *ResettableServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
//	mux.Handle(pattern, http.HandlerFunc(handler))
//}
//
//func (mux *ResettableServeMux) Handler(r *http.Request) (h http.Handler, pattern string) {
//	mux.mutex.RLock()
//	defer mux.mutex.RUnlock()
//
//	return mux.mux.Handler(r)
//}
//
//func (mux *ResettableServeMux) Reset() {
//	mux.mutex.Lock()
//	defer mux.mutex.Unlock()
//
//	mux.mux = http.NewServeMux()
//}
//
//func (mux *ResettableServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	mux.mutex.RLock()
//	defer mux.mutex.RUnlock()
//
//	mux.mux.ServeHTTP(w, r)
//}
