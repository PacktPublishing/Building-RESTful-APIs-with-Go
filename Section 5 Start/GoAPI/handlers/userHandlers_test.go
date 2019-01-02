package handlers

import (
	"GoAPI/user"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestBodyToUser(t *testing.T) {
	valid := &user.User{
		ID:   bson.NewObjectId(),
		Name: "John",
		Role: "Tester",
	}
	valid2 := &user.User{
		ID:   valid.ID,
		Name: "John",
		Role: "Developer",
	}
	js, err := json.Marshal(valid)
	if err != nil {
		t.Errorf("Error marshalling a valid user: %s", err)
		t.FailNow()
	}
	ts := []struct {
		txt string
		r   *http.Request
		u   *user.User
		err bool
		exp *user.User
	}{
		{
			txt: "nil request",
			err: true,
		},
		{
			txt: "empty request body",
			r:   &http.Request{},
			err: true,
		},
		{
			txt: "empty user",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString("{}")),
			},
			err: true,
		},
		{
			txt: "malformed data",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"id":12}`)),
			},
			u:   &user.User{},
			err: true,
		},
		{
			txt: "valid request",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBuffer(js)),
			},
			u:   &user.User{},
			exp: valid,
		},
		{
			txt: "valid partial request",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"role":"Developer","age":22}`)),
			},
			u:   valid,
			exp: valid2,
		},
	}

	for _, tc := range ts {
		t.Log(tc.txt)
		err := bodyToUser(tc.r, tc.u)
		if tc.err {
			if err == nil {
				t.Error("Expected error, got none.")
			}
			continue
		}
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
		if !reflect.DeepEqual(tc.u, tc.exp) {
			t.Error("Unmarshalled data is different:")
			t.Error(tc.u)
			t.Error(tc.exp)
		}
	}
}
