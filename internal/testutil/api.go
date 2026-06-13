package testutil

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
)

func NewTestAPI(t *testing.T) humatest.TestAPI {
	t.Helper()
	_, api := humatest.New(t, huma.DefaultConfig("Private Markets Test API", "0.0.1"))
	return api
}
