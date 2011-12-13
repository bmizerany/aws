# aws.go - A simple AWS client for Go.

# Install

	$ goinstall github.com/bmizerany/aws.go

## Use

	// NOTE: This project is still a work in progress.  For those of you who
	// see where I'm going with it and want to help, please do!

	package main

	import (
		"github.com/bmizerany/aws.go"
		"log"
	)

	// Assumes AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are set in ENV.
	func main() {
		v, err := aws.DescribeInstances()
		if err != nil {
			log.Fatal(err)
		}

		for _, r := range v.Reservations {
			for _, i := range r.Instances {
				log.Println("%s -> %s", i.DnsName, i.StateName)
			}
		}

		////
		// Custom (i.e. No sugar yet)

		r := aws.TemplateRequest // Make a copy of the template request

		// Set the Action parameter
		r.Add("Action", "DescribeKeyPairs")

		// Describe the response you expect using
		// http://weekly.golang.org/pkg/encoding/xml/#Unmarshal as a guide
		type DescribeKeyPairsResponse struct {
			aws.Header       // Holds the RequestId
			Keys       []Key `xml:"keySet>item"`
		}

		type Key struct {
			KeyName        string
			KeyFingerprint string
		}

		v := new(DescribeKeyPairsResponse)
		aws.Do(r, v)

		log.Printf("%v", v)
	}

## Change API/Creds

You can change the API and Credentials, used by the Sugar commands, by updating `aws.TemplateRequest`.

## Running tests

	$ cd /path/to/repo
	$ git config add aws.key <aws-key>
	$ git config add aws.secret <aws-secret>
	$ ./development.sh gotest // Does NOT create resources

or

	# Most *nix systems will not log a command with spaces in front in your history
	# BEWARE if yours does not
	$ <space><space> AWS_SECRET_ACCESS_KEY=<secret> AWS_ACCESS_KEY_ID=<key> gotest

## LICENCES

Copyright (C) 2011 by Blake Mizerany (@bmizerany)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE. 

