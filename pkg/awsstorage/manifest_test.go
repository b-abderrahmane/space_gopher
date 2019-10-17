package awsstorage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateManifestEntry(t *testing.T) {

	tim := time.Now().String()
	manifestEntry := generateManifestEntry("bucket001", "file001", 11, tim)
	url := "https://bucket001.s3-us-west-1.amazonaws.com/file001"

	assert.Equal(t, manifestEntry.Filename, "file001")
	assert.Equal(t, manifestEntry.LastModified, tim)
	assert.Equal(t, manifestEntry.URL, url)
	assert.Equal(t, manifestEntry.Size, int64(11))

}
