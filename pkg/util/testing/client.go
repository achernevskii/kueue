/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testing

import (
	"context"
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	kueue "sigs.k8s.io/kueue/apis/kueue/v1beta1"
	"sigs.k8s.io/kueue/pkg/controller/core/indexer"
)

func NewFakeClient(objs ...client.Object) client.Client {
	return NewClientBuilder().WithObjects(objs...).WithStatusSubresource(objs...).Build()
}

func NewClientBuilder(addToSchemes ...func(s *runtime.Scheme) error) *fake.ClientBuilder {
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		panic(err)
	}
	if err := kueue.AddToScheme(scheme); err != nil {
		panic(err)
	}
	for i := range addToSchemes {
		if err := addToSchemes[i](scheme); err != nil {
			panic(err)
		}
	}

	return fake.NewClientBuilder().WithScheme(scheme).
		WithIndex(&kueue.LocalQueue{}, indexer.QueueClusterQueueKey, indexer.IndexQueueClusterQueue).
		WithIndex(&kueue.Workload{}, indexer.WorkloadQueueKey, indexer.IndexWorkloadQueue).
		WithIndex(&kueue.Workload{}, indexer.WorkloadClusterQueueKey, indexer.IndexWorkloadClusterQueue)
}

type builderIndexer struct {
	*fake.ClientBuilder
}

func (b *builderIndexer) IndexField(ctx context.Context, obj client.Object, field string, extractValue client.IndexerFunc) error {
	b.ClientBuilder = b.ClientBuilder.WithIndex(obj, field, extractValue)
	return nil
}

func AsIndexer(builder *fake.ClientBuilder) client.FieldIndexer {
	return &builderIndexer{ClientBuilder: builder}
}

type EventRecord struct {
	Regarding types.NamespacedName
	Related   types.NamespacedName
	EventType string
	Reason    string
	Action    string
	Message   string
}

type EventRecorder struct {
	lock           sync.Mutex
	RecordedEvents []EventRecord
}

func (tr *EventRecorder) Eventf(regarding, related runtime.Object, eventtype, reason, action, note string, args ...interface{}) {
	tr.AnnotatedEventf(regarding, related, eventtype, reason, action, note, args...)
}

func (tr *EventRecorder) AnnotatedEventf(targetObject runtime.Object, relatedObject runtime.Object, eventtype, reason, action, note string, args ...interface{}) {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	regardingKey := types.NamespacedName{}
	if cobj, iscobj := targetObject.(client.Object); targetObject != nil && iscobj {
		regardingKey = client.ObjectKeyFromObject(cobj)
	}
	relatedKey := types.NamespacedName{}
	if cobj, iscobj := relatedObject.(client.Object); relatedObject != nil && iscobj {
		relatedKey = client.ObjectKeyFromObject(cobj)
	}
	tr.RecordedEvents = append(tr.RecordedEvents, EventRecord{
		Regarding: regardingKey,
		Related:   relatedKey,
		EventType: eventtype,
		Reason:    reason,
		Action:    action,
		Message:   fmt.Sprintf(note, args...),
	})
}
