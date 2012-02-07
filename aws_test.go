package aws

import (
	"encoding/xml"
	"github.com/bmizerany/assert"
	"strings"
	"testing"
)

func TestDoError(t *testing.T) {
	v, err := DescribeInstances()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v", v)
}

func TestUnmarshalError(t *testing.T) {
	body := strings.NewReader(`
		<?xml version="1.0" encoding="UTF-8"?>
		<Response>
			<Errors>
				<Error>
					<Code>AuthFailure</Code>
					<Message>AWS was not able to validate the provided access credentials</Message>
				</Error>
			</Errors>
			<RequestID>afc00dc9-0c19-46db-a987-f7de2a12a361</RequestID>
		</Response>
	`)

	type Error struct {
		Code    string
		Message string
	}

	type Response struct {
		RequestId string
		Errors    []Error `xml:"Errors>Error"`
	}

	got := new(Response)
	err := xml.Unmarshal(body, got)
	if err != nil {
		t.Fatal(err)
	}

	exp := &Response{
		RequestId: "afc00dc9-0c19-46db-a987-f7de2a12a361",
		Errors: []Error{
			{Code: "AuthFailure", Message: "AWS was not able to validate the provided access credentials"},
		},
	}

	assert.Equal(t, exp, got)
}
