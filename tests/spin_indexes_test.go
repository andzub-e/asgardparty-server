package tests

import (
	"net/http"
	"testing"
)

func TestSpinIndexesNotFound(t *testing.T) {
	manager.GenerateNewUserID()
	req := manager.DefaultSpinIndexesRequest()
	expectCode := http.StatusConflict

	code, content := SendRequest("POST", "core/spin_indexes/update", req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestSpinIndexesSuccess(t *testing.T) {
	SendSuccessWagerRequest()

	req := manager.DefaultSpinIndexesRequest()
	expectCode := http.StatusOK

	code, content := SendRequest("POST", "core/spin_indexes/update", req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestSpinIndexesNoSession(t *testing.T) {
	req := manager.DefaultSpinIndexesRequest()
	expectCode := http.StatusUnprocessableEntity

	req.SessionToken = ""

	code, content := SendRequest("POST", "core/spin_indexes/update", req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestSpinIndexesWrongSession(t *testing.T) {
	req := manager.DefaultSpinIndexesRequest()
	expectCode := http.StatusUnauthorized

	req.SessionToken = manager.GetWrongSessionToken()

	code, content := SendRequest("POST", "core/spin_indexes/update", req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestSpinIndexesNoBaseIndex(t *testing.T) {
	req := manager.DefaultSpinIndexesRequest()
	expectCode := http.StatusUnprocessableEntity

	req.BaseSpinIndex = nil

	code, content := SendRequest("POST", "core/spin_indexes/update", req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestSpinIndexesNegativeBaseIndex(t *testing.T) {
	req := manager.DefaultSpinIndexesRequest()
	expectCode := http.StatusUnprocessableEntity

	minusOne := -1
	req.BaseSpinIndex = &minusOne

	code, content := SendRequest("POST", "core/spin_indexes/update", req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestSpinIndexesNoBonusIndex(t *testing.T) {
	req := manager.DefaultSpinIndexesRequest()
	expectCode := http.StatusUnprocessableEntity

	req.BonusSpinIndex = nil

	code, content := SendRequest("POST", "core/spin_indexes/update", req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestSpinIndexesNegativeBonusIndex(t *testing.T) {
	req := manager.DefaultSpinIndexesRequest()
	expectCode := http.StatusUnprocessableEntity

	minusOne := -1
	req.BonusSpinIndex = &minusOne

	code, content := SendRequest("POST", "core/spin_indexes/update", req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}
