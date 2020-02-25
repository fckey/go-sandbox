package retry

import (
	"testing"
)

func TestRunWithHTTPRequestRetry(t *testing.T) {
	//ctx := context.Background()
	//// Without a context deadline, retry will run until the function
	//// says not to retry any more.
	//n := 0
	//
	//server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	w.WriteHeader(http.StatusRequestTimeout)
	//}))
	//defer server.Close()
	//
	//http.Client{}
	//
	//retriableFunc := func() error {
	//	n++
	//	if n < 5 {
	//		return middleRetry
	//	}
	//	return endRetry
	//}
	//checkerFunc := func(err error) bool {
	//	if strings.Compare(err.Error(), retriableMsg) == 0 {
	//		return true
	//	}
	//	return false
	//}
	//backoff := NewBackoff(1*time.Second, 32*time.Second, 2)
	//err := RunWithHTTPRetry(ctx, backoff, retriableFunc, checkerFunc)
	//if n != 5 {
	//	t.Errorf("n: got %d, want %d", n, 10)
	//}
	//
	//if err == nil {
	//	t.Error("got nil, want error")
	//}
}