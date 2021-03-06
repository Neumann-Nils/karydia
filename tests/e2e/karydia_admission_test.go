// Copyright 2019 Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAutomountServiceAccountToken(t *testing.T) {
	namespace, err := f.CreateTestNamespace()
	if err != nil {
		t.Fatalf("failed to create test namespace: %v", err)
	}

	namespace.ObjectMeta.Annotations = map[string]string{
		"karydia.gardener.cloud/automountServiceAccountToken": "forbidden",
	}

	namespace, err = f.KubeClientset.CoreV1().Namespaces().Update(namespace)
	if err != nil {
		t.Fatalf("failed to annotate test namespace: %v", err)
	}

	ns := namespace.ObjectMeta.Name

	automountServiceAccountToken := true

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "karydia-e2e-test-pod",
			Namespace: ns,
		},
		Spec: corev1.PodSpec{
			AutomountServiceAccountToken: &automountServiceAccountToken,
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx",
				},
			},
		},
	}

	_, err = f.KubeClientset.CoreV1().Pods(ns).Create(pod)
	if err == nil {
		t.Fatalf("expected pod creation to fail")
	}

	pod.Spec.AutomountServiceAccountToken = nil

	_, err = f.KubeClientset.CoreV1().Pods(ns).Create(pod)
	if err == nil {
		t.Fatalf("expected pod creation to fail")
	}

	automountServiceAccountToken = false
	pod.Spec.AutomountServiceAccountToken = &automountServiceAccountToken

	pod, err = f.KubeClientset.CoreV1().Pods(ns).Create(pod)
	if err != nil {
		t.Fatalf("failed to create pod: %v", err)
	}

	timeout := 2 * time.Minute
	if err := f.WaitPodRunning(pod.ObjectMeta.Namespace, pod.ObjectMeta.Name, timeout); err != nil {
		t.Fatalf("pod never reached state running")
	}
}
