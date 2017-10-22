package etcd

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"io/ioutil"
	"log"
	"testing"
	"time"
	"math/rand"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {
	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func randomString(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func create(id string, data string, t *testing.T) (act activity.Activity, tc *test.TestActivityContext) {
	act = NewActivity(getActivityMetadata())
	tc = test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput("key", id)
	if data == "" {
		tc.SetInput("value", `foo`)
	} else {
		tc.SetInput("value", data)
	}

	tc.SetInput("method", `Create`)
	tc.SetInput("servers", "http://etcd:2379")

	_, insertError := act.Eval(tc)
	if insertError != nil {
		t.Error("Create error")
		t.Fail()
	}
	return
}

/*
Test a create activity
*/
func TestNewActivity(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
	log.Println("TestNewActivity successfull")
}

/*
Create test
 */
func TestCreate(t *testing.T) {
	id := randomString(5)
	create(id, "", t)
	log.Println("TestCreate successfull")
}

/*
Delete test
 */
func TestDelete(t *testing.T) {
	id := randomString(5)
	act, tc := create(id, "", t)

	tc.SetInput("method", "Delete")

	_, RemoveError := act.Eval(tc)
	if RemoveError != nil {
		t.Error("Document not removed")
		t.Fail()
	}
	log.Println("TestRemove successfull")
}

/*
Update test
*/
func TestUpdate(t *testing.T) {
	id := randomString(5)
	act, tc := create(id, "", t)
	data := "bar"

	tc.SetInput("method", "Update")
	tc.SetInput("value", data)

	_, RemoveError := act.Eval(tc)
	if RemoveError != nil {
		t.Error("Document not upserted")
		t.Fail()
	}

	tc.SetInput("method", "Update")
	act.Eval(tc)

	s := tc.GetOutput("output").(string)

	if s != data {
		t.Error("The retrieved document is not equals")
		t.Fail()
	}

	log.Println("TestUpdate successfull")
}

/*
Get test
*/
func TestGet(t *testing.T) {
	id := randomString(5)
	data := "foo"
	act, tc := create(id, data, t)

	tc.SetInput("method", "Get")

	_, GetError := act.Eval(tc)
	if GetError != nil {
		t.Error("Document not retrieved")
		t.Fail()
	}
	s := tc.GetOutput("output").(string)

	if s != data {
		t.Error("The retrieved document is not equals")
		t.Fail()
	}

	log.Println("TestGet successfull")
}
