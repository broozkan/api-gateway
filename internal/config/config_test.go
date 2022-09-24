package config_test

import (
	"broozkan/api-gateway/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

const configPath = "../../test/testdata"
const configName = "test-config"

func TestGivenTestConfigFileWhenICallNewThenItShouldReturnConfig(t *testing.T) {
	actualConfig, _ := config.New(configPath, configName)

	expectedConfig := config.Config{
		Appname: "something-special",
		Server:  config.Server{Port: "1111"},
	}

	assert.Equal(t, expectedConfig, actualConfig)
}

func TestGivenNotExistingConfigFileWhenICallNewThenItShouldReturnError(t *testing.T) {
	fakeConfigPath := "../test/fake"
	notExistingConfigName := "nothing"

	_, err := config.New(fakeConfigPath, notExistingConfigName)

	assert.NotNil(t, err)
}

func TestGivenBadConfigurationWhenICallNewThenItShouldReturnError(t *testing.T) {
	badConfigName := "test-bad-config"

	_, err := config.New(configPath, badConfigName)

	assert.NotNil(t, err)
}
