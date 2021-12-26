package payload_parser

import (
	"fmt"
	"net/http"
	"strings"
)

type RequestShape struct {
	QueryString    string `query:"q,required"`
	RequiredString string `query:"rs,required"`
	Referrer       string `header:"referer,-"`
}

func (rs *RequestShape) ValidateQueryString() error {
	if rs.QueryString == "<script>lol</script>" {
		rs.QueryString = ""
	}
	return nil
}

func (rs *RequestShape) ValidateRequiredString() error {
	rs.RequiredString = strings.Trim(rs.RequiredString, " ")
	if len(rs.RequiredString) == 0 {
		return fmt.Errorf("error!")
	}
	return nil
}

func ExampleParsePayload() {
	requestPayload := &RequestShape{}
	request, err := http.NewRequest("GET", "/example?q=lol&rs=+x", nil)
	if err != nil {
		fmt.Println("Invalid")
	}
	request.Header.Add("referer", "google.com")

	err = ParsePayload(requestPayload, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(requestPayload.QueryString)
	fmt.Println(requestPayload.Referrer)
	fmt.Println(requestPayload.RequiredString)

	// Output:
	// lol
	// google.com
	// x
}
