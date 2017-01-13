// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"archive/zip"
	"enjarify-go/jvm"
	"enjarify-go/util"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

// load stub files lazily
var STUB_FILES [][2]string

func executeTest(name string, opts jvm.Options) {
	fmt.Printf("Running test %s\n", name)
	dir := "../tests/" + name
	data := Read(dir + "/classes.dex")

	results := translate(opts, data)
	classes := make([][2]string, len(results))
	for i := range results {
		classes[i][0] = results[i][0]
		classes[i][1] = results[i][1]
		util.Assert(results[i][1] != "")
	}
	writeToJar("out.jar", classes)

	out, err := exec.Command("java", "-Xss515m", "-jar", "out.jar", "a.a").CombinedOutput()
	check(err)
	result := string(out)

	expected := Read(dir + "/expected.txt")
	expected = strings.Replace(expected+"\n", "\r\n", "\n", -1)

	if result != expected {
		panic(util.Unreachable)
	}
}

func runTests() {
	stubs := [][2]string{}
	r, err := zip.OpenReader("../tests/stubs/stubs.zip")
	check(err)
	for _, f := range r.File {
		rc, err := f.Open()
		check(err)
		data, err := ioutil.ReadAll(rc)
		check(err)
		stubs = append(stubs, [2]string{f.Name, string(data)})
		rc.Close()
	}
	r.Close()
	STUB_FILES = stubs

	// Now do the tests
	for _, opts := range [...]jvm.Options{jvm.NONE, jvm.PRETTY, jvm.ALL} {
		for i := 1; i < 7; i++ {
			executeTest(fmt.Sprintf("test%d", i), opts)
		}
	}
}
