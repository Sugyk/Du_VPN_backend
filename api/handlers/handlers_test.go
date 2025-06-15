package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sugyk/rest_vpn/lib/configs"
)

// Test function using sqlmock
func TestGetAccessKeysList(t *testing.T) {
	// init mock db
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	if err != nil {
		t.Fatalf("error initializing sqlmock: %s", err)
	}

	// init service with fake db
	handler := newService(sqlxDB, configs.OutlineAPIConfig{API_URL: "test"})

	// mock db queries
	test_time := time.Date(2000, 1, 1, 1, 1, 1, 1, &time.Location{})

	mock.ExpectQuery(`
	SELECT id, key_value, user_id, outline_id, created_at
	FROM "AccessKeys"
	WHERE \(\$1 = 0 OR id = \$2\)
	AND \(\$3 = 0 OR user_id = \$4\)
	AND \(\$5 = 0 OR outline_id = \$6\)
	`).WithArgs(1, 1, 1234, 1234, 3, 3).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "key_value", "user_id", "outline_id", "created_at"}).
				AddRow(1, "test", 1234, 3, test_time),
		)

	mock.ExpectQuery(`
	SELECT id, key_value, user_id, outline_id, created_at
	FROM "AccessKeys"
	WHERE \(\$1 = 0 OR id = \$2\)
	AND \(\$3 = 0 OR user_id = \$4\)
	AND \(\$5 = 0 OR outline_id = \$6\)
	`).WithArgs(1, 1, 0, 0, 0, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "key_value", "user_id", "outline_id", "created_at"}).
				AddRow(1, "test", 4321, 12, test_time),
		)

	mock.ExpectQuery(`
	SELECT id, key_value, user_id, outline_id, created_at
	FROM "AccessKeys"
	WHERE \(\$1 = 0 OR id = \$2\)
	AND \(\$3 = 0 OR user_id = \$4\)
	AND \(\$5 = 0 OR outline_id = \$6\)
	`).WithArgs(0, 0, 0, 0, 12, 12).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "key_value", "user_id", "outline_id", "created_at"}).
				AddRow(12, "test", 4321, 12, test_time),
		)

	mock.ExpectQuery(`
	SELECT id, key_value, user_id, outline_id, created_at
	FROM "AccessKeys"
	WHERE \(\$1 = 0 OR id = \$2\)
	AND \(\$3 = 0 OR user_id = \$4\)
	AND \(\$5 = 0 OR outline_id = \$6\)
	`).WithArgs(0, 0, 1234, 1234, 0, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "key_value", "user_id", "outline_id", "created_at"}).
				AddRow(444, "test", 1234, 444, test_time),
		)

	// test with all filter params
	req_all_params := httptest.NewRequest(http.MethodGet, "/find?id=1&user_id=1234&outline_id=3", nil)
	rr_all_params := httptest.NewRecorder()

	handler.GetAccessKeysList().ServeHTTP(rr_all_params, req_all_params)

	if rr_all_params.Code != http.StatusOK {
		// wrong status
		t.Errorf("Test all params\nExpected status: %d\nStatus: %d\nBody: %s\nHeaders: %#v", http.StatusOK, rr_all_params.Code, rr_all_params.Body.String(), rr_all_params.Result().Header)
	}

	expected_body_all_params := `{"access_keys":[{"id":1,"key_value":"test","user_id":1234,"outline_id":3,"created_at":"2000-01-01T01:01:01.000000001Z"}]}`
	if expected_body_all_params != rr_all_params.Body.String() {
		// wrong body
		t.Errorf("Test all params\nExpected body: %s\nStatus: %d\nBody: %s\nHeaders: %#v", expected_body_all_params, rr_all_params.Code, rr_all_params.Body.String(), rr_all_params.Result().Header)
	}

	// test only with id
	req_only_id := httptest.NewRequest(http.MethodGet, "/find?id=1", nil)
	rr_only_id := httptest.NewRecorder()

	handler.GetAccessKeysList().ServeHTTP(rr_only_id, req_only_id)
	if rr_only_id.Code != http.StatusOK {
		t.Errorf("Test only id\nCode: %d\nBody: %s\nHeaders: %#v", rr_only_id.Code, rr_only_id.Body.String(), rr_only_id.Result().Header)
	}

	// test only with outline_id
	req_only_outline_id := httptest.NewRequest(http.MethodGet, "/find?outline_id=12", nil)
	rr_only_outline_id := httptest.NewRecorder()

	handler.GetAccessKeysList().ServeHTTP(rr_only_outline_id, req_only_outline_id)
	if rr_only_outline_id.Code != http.StatusOK {
		t.Errorf("Test only outline_id\nCode: %d\nBody: %s\nHeaders: %#v", rr_only_outline_id.Code, rr_only_outline_id.Body.String(), rr_only_outline_id.Result().Header)
	}

	// test only with user_id
	req_only_user_id := httptest.NewRequest(http.MethodGet, "/find?user_id=1234", nil)
	rr_only_user_id := httptest.NewRecorder()

	handler.GetAccessKeysList().ServeHTTP(rr_only_user_id, req_only_user_id)
	if rr_only_user_id.Code != http.StatusOK {
		t.Errorf("Test only user_id\nCode: %d\nBody: %s\nHeaders: %#v", rr_only_user_id.Code, rr_only_user_id.Body.String(), rr_only_user_id.Result().Header)
	}

	// test response 500
	req_wrong_filter := httptest.NewRequest(http.MethodGet, "/find?wrong_filter=1234", nil)
	rr_wrong_filter := httptest.NewRecorder()

	handler.GetAccessKeysList().ServeHTTP(rr_wrong_filter, req_wrong_filter)
	if rr_wrong_filter.Code != http.StatusBadRequest {
		t.Errorf("Test only user_id\nCode: %d\nBody: %s\nHeaders: %#v", rr_wrong_filter.Code, rr_wrong_filter.Body.String(), rr_wrong_filter.Result().Header)
	}
}
