// Copyright 2017 Google LLC
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

package firestore

import (
	"context"
	"testing"
)

func TestCGR_TestQueryPartition_ToQuery(t *testing.T) {
	cgr := newCollectionGroupRef(testClient, testClient.path(), "collectionID")
	qp := QueryPartition{
		CollectionGroupQuery: cgr.Query,
		StartAt:              "projects/projectID/databases/(default)/documents/start/at",
		EndBefore:            "projects/projectID/databases/(default)/documents/end/before",
	}

	got := qp.ToQuery()

	want := Query{
		c:              testClient,
		path:           "projects/projectID/databases/(default)",
		parentPath:     "projects/projectID/databases/(default)/documents",
		collectionID:   "collectionID",
		startVals:      []interface{}{"projects/projectID/databases/(default)/documents/start/at"},
		endVals:        []interface{}{"projects/projectID/databases/(default)/documents/end/before"},
		startBefore:    true,
		endBefore:      true,
		allDescendants: true,
		orders:         []order{{fieldPath: []string{"__name__"}, dir: 1}},
	}

	if !testEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestCGR_TestGetPartitions(t *testing.T) {
	cgr := newCollectionGroupRef(testClient, testClient.path(), "collectionID")
	_, err := cgr.GetPartitions(context.Background(), 0)
	if err == nil {
		t.Error("Expected an error when requested partition count is < 1")
	}

	parts, err := cgr.GetPartitions(context.Background(), 1)
	if err != nil {
		t.Error("Didn't expect an error when requested partition count is 1")
	}
	if len(parts) != 1 {
		t.Fatal("Expected 1 QueryPartition")
	}
	got := parts[0]
	want := QueryPartition{CollectionGroupQuery: cgr.Query, StartAt: "", EndBefore: ""}
	if !testEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}