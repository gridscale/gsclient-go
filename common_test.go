package gsclient

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_isValidUUID(t *testing.T) {
	validationUUIDTestCases := make([]uuidTestCase, len(uuidCommonTestCases))
	copy(validationUUIDTestCases, uuidCommonTestCases)
	validationUUIDTestCases = append(validationUUIDTestCases,
		uuidTestCase{
			isFailed: true,
			testUUID: "abc-123",
		},
		uuidTestCase{
			isFailed: false,
			testUUID: "690de890-13c0-4e76-8a01-e10ba8786e53",
		},
	)
	for _, test := range validationUUIDTestCases {
		isValid := isValidUUID(test.testUUID)
		if test.isFailed {
			assert.False(t, isValid)
		} else {
			assert.True(t, isValid)
		}
	}
}

func Test_retryWithTimeout(t *testing.T) {
	type testCase struct {
		isContinue     bool
		err            error
		timeout, delay time.Duration
	}
	testCases := []testCase{
		{
			true,
			nil,
			time.Duration(1) * time.Second,
			time.Duration(100) * time.Millisecond,
		},
		{
			false,
			nil,
			time.Duration(1) * time.Second,
			time.Duration(100) * time.Millisecond,
		},
		{
			false,
			errors.New("just test"),
			time.Duration(1) * time.Second,
			time.Duration(100) * time.Millisecond,
		},
		{
			true,
			errors.New("just test"),
			time.Duration(1) * time.Second,
			time.Duration(100) * time.Millisecond,
		},
	}
	for _, test := range testCases {
		err := retryWithTimeout(func() (bool, error) {
			return test.isContinue, test.err
		}, test.timeout, test.delay)
		if test.err != nil || test.isContinue {
			assert.NotNil(t, err, fmt.Sprintf("%v %v", err, test.err))
		} else {
			assert.Nil(t, err, err)
		}
	}
}

func Test_retryWithLimitedNumOfRetries(t *testing.T) {
	type testCase struct {
		isContinue   bool
		err          error
		delay        time.Duration
		numOfRetries int
	}
	testCases := []testCase{
		{
			true,
			nil,
			time.Duration(500) * time.Millisecond,
			10,
		},
		{
			false,
			nil,
			time.Duration(500) * time.Millisecond,
			10,
		},
		{
			false,
			errors.New("just test"),
			time.Duration(500) * time.Millisecond,
			10,
		},
	}
	for _, test := range testCases {
		err := retryWithLimitedNumOfRetries(func() (bool, error) {
			return test.isContinue, test.err
		}, test.numOfRetries, test.delay)
		if test.err != nil || test.isContinue {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}
