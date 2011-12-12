package aws

import (
	"bytes"
	"crypto/hmac"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	EC2Host    = "ec2.amazonaws.com"
	EC2Version = "2011-11-01"
)

var TemplateRequest = Request{
	Host:    EC2Host,
	Version: EC2Version,
	Key:     os.Getenv("AWS_ACCESS_KEY_ID"),
	Secret:  os.Getenv("AWS_SECRET_ACCESS_KEY"),
}

type Param struct {
	Key string
	Val string
}

func (p *Param) Encode() string {
	return p.Key + "=" + url.QueryEscape(p.Val)
}

type Params []*Param

func (p *Params) Add(key, val string) {
	*p = append(*p, &Param{key, val})
}

func (p *Params) Len() int {
	return len(*p)
}

func (p *Params) Less(i, j int) bool {
	a := *p
	return a[i].Key < a[j].Key
}

func (p *Params) Swap(i, j int) {
	a := *p
	a[i], a[j] = a[j], a[i]
}

func (p *Params) Encode() (s string) {
	parts := make([]string, len(*p))
	for i, param := range *p {
		parts[i] = param.Encode()
	}
	return strings.Join(parts, "&")
}

type Request struct {
	Host    string
	Key     string
	Secret  string
	Version string
	Params
}

func (r *Request) Encode() string {
	r.Add("AWSAccessKeyId", r.Key)
	r.Add("SignatureMethod", "HmacSHA256")
	r.Add("SignatureVersion", "2")
	r.Add("Version", r.Version)
	r.Add("Timestamp", time.Now().UTC().Format(time.RFC3339))

	sort.Sort(r)

	data := strings.Join([]string{
		"POST",
		r.Host,
		"/",
		r.Params.Encode(),
	}, "\n")

	h := hmac.NewSHA256([]byte(r.Secret))
	h.Write([]byte(data))

	sig := base64.StdEncoding.EncodeToString(h.Sum([]byte{}))

	r.Add("Signature", sig)

	fmt.Printf("%v\n", r.Params.Encode())

	return r.Params.Encode()
}

type Header struct {
	RequestId string
}

type Error struct {
	Header
	Errors []struct {
		Code string
		Message string
	} `xml:"Errors>Error"`
}

// Example:
//  aws: ->
//    AuthFailure: "There is a problem with your secret"
//    OMG: "You're servers are all gone!"
func (err *Error) Error() string {
	var s string
	for _, e := range err.Errors {
		s += fmt.Sprintf("\t%s: %q\n", e.Code, e.Message)
	}

	return fmt.Sprintf("aws: ->\n%s", s)
}

func Do(r *Request, v interface{}) error {
	res, err := http.Post(
		"https://"+r.Host,
		"application/x-www-form-urlencoded",
		bytes.NewBufferString(r.Encode()),
	)
	if err != nil {
		return err
	}

	return unmarshal(res, v)
}

func unmarshal(res *http.Response, v interface{}) error {
	if res.StatusCode != http.StatusOK {
		e := new(Error)
		err := xml.Unmarshal(&logReader{res.Body}, e)
		if err != nil {
			return err
		}
		return e
	}

	return xml.Unmarshal(&logReader{res.Body}, v)
}

// Utils

// Used for debugging
type logReader struct {
	r io.Reader
}

func (lr *logReader) Read(b []byte) (n int, err error) {
	n, err = lr.r.Read(b)
	fmt.Print(string(b))
	return
}

// Sugar
type DescribeInstancesResponse struct {
	Header
	Reservations []Reservation `xml:"reservationSet>item"`
}

type Reservation struct {
	ReservationId string
	Instances []Instance `xml:"instancesSet>item"`
}

type Instance struct {
	InstanceId string
	StateName string `xml:"instanceState>name"`
	DnsName string
	IpAddress string
}

func DescribeInstances() (*DescribeInstancesResponse, error) {
	r := TemplateRequest
	r.Add("Action", "DescribeInstances")

	v := new(DescribeInstancesResponse)
	return v, Do(&r, v)
}
