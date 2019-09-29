package clouds

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"satellity/internal/configs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImage(t *testing.T) {
	setupTestContext()
	assert := assert.New(t)

	file, err := ioutil.ReadFile("./testdata/city.jpg")
	assert.Nil(err)
	data := base64.StdEncoding.EncodeToString(file)

	str, err := UploadImage(nil, "tests/images/city", data)
	assert.Nil(err)
	assert.Equal("tests/images/city.jpeg", str)
}

const (
	testEnvironment = "test"
)

func setupTestContext() {
	if err := configs.Init("./../configs", testEnvironment); err != nil {
		log.Panicln(err)
	}
}
