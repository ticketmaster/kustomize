// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package loader_test

import (
	"testing"

	"sigs.k8s.io/kustomize/v3/api/internal/loadertest"
	"sigs.k8s.io/kustomize/v3/api/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/v3/api/plugins/config"
	. "sigs.k8s.io/kustomize/v3/api/plugins/loader"
	"sigs.k8s.io/kustomize/v3/api/resmap"
	"sigs.k8s.io/kustomize/v3/api/resource"
	"sigs.k8s.io/kustomize/v3/api/testutils/kusttest"
	"sigs.k8s.io/kustomize/v3/api/testutils/valtest"
)

const (
	secretGenerator = `
apiVersion: builtin
kind: SecretGenerator
metadata:
  name: secretGenerator
name: mySecret
behavior: merge
envFiles:
- a.env
- b.env
valueFiles:
- longsecret.txt
literals:
- FRUIT=apple
- VEGETABLE=carrot
`
	someServiceGenerator = `
apiVersion: someteam.example.com/v1
kind: SomeServiceGenerator
metadata:
  name: myServiceGenerator
service: my-service
port: "12345"
`
)

func TestLoader(t *testing.T) {
	tc := kusttest_test.NewPluginTestEnv(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"builtin", "", "SecretGenerator")
	tc.BuildGoPlugin(
		"someteam.example.com", "v1", "SomeServiceGenerator")

	rmF := resmap.NewFactory(resource.NewFactory(
		kunstruct.NewKunstructuredFactoryImpl()), nil)

	ldr := loadertest.NewFakeLoader("/foo")

	pLdr := NewLoader(config.ActivePluginConfig(), rmF)
	if pLdr == nil {
		t.Fatal("expect non-nil loader")
	}

	m, err := rmF.NewResMapFromBytes([]byte(
		someServiceGenerator + "---\n" + secretGenerator))
	if err != nil {
		t.Fatal(err)
	}

	_, err = pLdr.LoadGenerators(ldr, valtest_test.MakeFakeValidator(), m)
	if err != nil {
		t.Fatal(err)
	}
}
