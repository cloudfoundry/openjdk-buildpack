/*
 * Copyright 2018 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jdk_test

import (
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/cloudfoundry/openjdk-buildpack/jdk"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestJDK(t *testing.T) {
	spec.Run(t, "JDK", testJDK, spec.Report(report.Terminal{}))
}

func testJDK(t *testing.T, when spec.G, it spec.S) {

	it("returns true if build plan exists", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, jdk.Dependency, "stub-openjdk-jdk.tar.gz")
		f.AddBuildPlan(t, jdk.Dependency, buildplan.Dependency{})

		_, ok, err := jdk.NewJDK(f.Build)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("NewJDK = %t, expected true", ok)
		}
	})

	it("returns false if build plan does not exist", func() {
		f := test.NewBuildFactory(t)

		_, ok, err := jdk.NewJDK(f.Build)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Errorf("NewJDK = %t, expected false", ok)
		}
	})

	it("contributes JDK", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, jdk.Dependency, "stub-openjdk-jdk.tar.gz")
		f.AddBuildPlan(t, jdk.Dependency, buildplan.Dependency{})

		j, _, err := jdk.NewJDK(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		if err := j.Contribute(); err != nil {
			t.Fatal(err)
		}

		layer := f.Build.Layers.Layer("openjdk-jdk")
		test.BeLayerLike(t, layer, true, true, false)
		test.BeFileLike(t, filepath.Join(layer.Root, "fixture-marker"), 0644, "")
		test.BeFileLike(t, filepath.Join(layer.Root, "env.build", "JAVA_HOME.override"), 0644, layer.Root)
		test.BeFileLike(t, filepath.Join(layer.Root, "env.build", "JDK_HOME.override"), 0644, layer.Root)
	})
}