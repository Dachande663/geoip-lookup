package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldRejectEmptyIpAddress(t *testing.T) {
	info, err := lookupIp("")

	assert.Nil(t, info.Name)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "An IP address must be provided.")
}

func TestShouldRejectInvalidIpAddress(t *testing.T) {
	info, err := lookupIp("123.456.789")

	assert.Nil(t, info.Name)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "IP address is not a valid IP.")
}

func TestShouldRejectLoopbackIpAddress(t *testing.T) {
	info, err := lookupIp("127.0.0.1")

	assert.Nil(t, info.Name)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "IP address is not allowed.")
}

func TestShouldRejectPrivateIpAddress(t *testing.T) {
	info, err := lookupIp("192.168.1.89")

	assert.Nil(t, info.Name)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "IP address is not allowed.")
}
