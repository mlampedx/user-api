package main_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"."
)

var a main.API

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS users
(
	id SERIAL,
	username TEXT NOT NULL,
	password BYTEA NOT NULL,
	data JSON,
	adminGroupId INTEGER NULL,
	dateCreated TIMESTAMPTZ NOT NULL,
	dateModified TIMESTAMPTZ NOT NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id)
)`

func TestMain(m *testing.M) {
	a = main.API{}
	a.Initialize(
		os.Getenv("TEST_PSQL_DB_USERNAME"),
		os.Getenv("TEST_PSQL_DB_PASSWORD"),
		os.Getenv("TEST_USER_API_PSQL_DB_NAME"))
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}
func TestEmptyTable(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/users", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s.\n", body)
	}
}

func TestGetNonExistentUser(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/user/123", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(res.Body.Bytes(), &m)
	if m["error"] != "User not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'User not found'. Got '%s'", m["error"])
	}
}
func TestCreateUser(t *testing.T) {
	clearTable()
	testUser := []byte(`{"id":123, "username":"jdoe", "password":$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK, "test":["data"], "adminGroupId":1, "dateCreated": 2017-05-13 15:50:13.793654 +0000 UTC, "dateModified": 2017-05-13 15:50:13.793654 +0000 UTC}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(testUser))
	res := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, res.Code)

	var m map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &m)

	// json.Unmarshal converts numbers to floats when the target is type map[string]interface{}
	if m["id"] != 123.0 {
		t.Errorf("Expected id to be '123'. Got '%v'", m["id"])
	}

	if m["username"] != "jdoe" {
		t.Errorf("Expected username to be 'jdoe'. Got '%v'", m["username"])
	}

	if m["password"] != $2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK {
		t.Errorf("Expected hashed password to be '$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK'. Got '%v'", m["password"])
	}

	if m["data"] != []string{"data"} {
		t.Errorf("Expected data to be '["data"]'. Got '%v'", m["data"])
	}

	if m["adminGroupId"] != 1.0 {
		t.Errorf("Expected adminGroupId to be '1'. Got '%v'", m["adminGroupId"])
	}

	if m["dateCreated"] != 2017-05-13 15:50:13.793654 +0000 UTC {
		t.Errorf("Expected dateCreated to be '2017-05-13 15:50:13.793654 +0000 UTC'. Got '%v'", m["dateCreated"])
	}

	if m["dateModified"] != 2017-05-13 15:50:13.793654 +0000 UTC {
		t.Errorf("Expected dateModified to be '2017-05-13 15:50:13.793654 +0000 UTC'. Got '%v'", m["dateModified"])
	}

}

func TestRetrieveUser(t *testing.T) {
	clearTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/user/1", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)
}

func TestUpdateUser(t *testing.T) {
	clearTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/user/1", nil)
	res := executeRequest(req)
	var originalUser map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &originalUser)
	payload := []byte(`{"data":["newdata"]}`)
	req, _ = http.NewRequest("PUT", "/user/1", bytes.NewBuffer(payload))
	res = executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)
	var m map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &m)

	if m["data"] == originalUser["data"] {
		t.Errorf("Expected the data to change from '%v' to '%v'. Got '%v'", originalUser["data"], m["data"], m["data"])
	}
}

func TestDeleteUser(t *testing.T) {
	clearTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/user/1", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)

	req, _ = http.NewRequest("DELETE", "/user/1", nil)
	res = executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)

	req, _ = http.NewRequest("GET", "/user/1", nil)
	res = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, res.Code)
}

func addUsers(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO users(id, username, password, data, adminGroupId, dateCreated, dateModified) VALUES($1, $2, $3, $4, $5, $6, $7)", "User "(i+1.0)*10, "jdoe", []byte{$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK}, []string{"data"}, 1, time.Now(), time.Now())
	}
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM users")
	a.DB.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d.\n", expected, actual)
	}
}
