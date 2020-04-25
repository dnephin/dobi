package image

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestRecordPath(t *testing.T) {
	mockClient, teardown := setupMockClient(t)
	defer teardown()
	ctx, config := setupCtxAndConfig(mockClient)
	path := recordPath(ctx, config)
	assert.Equal(t, "/dir/.dobi/images/imagename tag", path)
}

func TestRecordPathEscapesSlash(t *testing.T) {
	mockClient, teardown := setupMockClient(t)
	defer teardown()
	ctx, config := setupCtxAndConfig(mockClient)
	config.Image = "repo/name"
	path := recordPath(ctx, config)
	assert.Equal(t, "/dir/.dobi/images/repo name tag", path)
}
