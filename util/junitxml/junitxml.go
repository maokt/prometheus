// Copyright 2021 Marty Pauley
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package junitxml

import (
	"encoding/xml"
	"io"
)

type JUnitXML struct {
	XMLName xml.Name     `xml:"testsuites"`
	Suites  []*TestSuite `xml:"testsuite"`
}

type TestSuite struct {
	Name         string      `xml:"name,attr"`
	TestCount    int         `xml:"tests,attr"`
	FailureCount int         `xml:"failures,attr"`
	ErrorCount   int         `xml:"errors,attr"`
	Cases        []*TestCase `xml:"testcase"`
}

type TestCase struct {
	Name     string   `xml:"name,attr"`
	Failures []string `xml:"failure,omitempty"`
	Error    string   `xml:"error,omitempty"`
}

func (j *JUnitXML) Suite(name string) *TestSuite {
	ts := &TestSuite{Name: name}
	j.Suites = append(j.Suites, ts)
	return ts
}

func (j *JUnitXML) WriteXML(h io.Writer) error {
	return xml.NewEncoder(h).Encode(j)
}

func (ts *TestSuite) Case(name string) *TestSuite {
	ts.TestCount++
	tc := &TestCase{Name: name}
	ts.Cases = append(ts.Cases, tc)
	return ts
}

func (ts *TestSuite) lastCase() *TestCase {
	if len(ts.Cases) == 0 {
		ts.Case("unknown")
	}
	return ts.Cases[len(ts.Cases)-1]
}

func (ts *TestSuite) Fail(f string) {
	ts.FailureCount++ // yes, there can be more failures than test cases
	curt := ts.lastCase()
	curt.Failures = append(curt.Failures, f)
}

func (ts *TestSuite) Abort(e error) {
	ts.ErrorCount++
	curt := ts.lastCase()
	curt.Error = e.Error()
}
