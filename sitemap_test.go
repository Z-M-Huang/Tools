package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildSitemap(t *testing.T) {
	assert.NotEmpty(t, BuildSitemap())
}
