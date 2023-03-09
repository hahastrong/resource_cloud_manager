package rcm

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)



func handler( w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}


func TestRegistry(t *testing.T) {
	http.HandleFunc("/api/v1/resource/manage", handler)
	_ = http.ListenAndServe(":8080", nil)


}
