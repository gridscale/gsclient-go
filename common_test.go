package gsclient

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_isValidUUID(t *testing.T) {
	validUUIDTest := make([]uuidTestCase, len(uuidCommonTestCases))
	copy(validUUIDTest, uuidCommonTestCases)
	validUUIDTest = append(validUUIDTest, uuidTestCase{
		isFailed: true,
		testUUID: "abc-123",
	})
	for _, test := range validUUIDTest {
		isValid := isValidUUID(test.testUUID)
		if test.isFailed {
			assert.False(t, isValid)
		} else {
			assert.True(t, isValid)
		}
	}
}
