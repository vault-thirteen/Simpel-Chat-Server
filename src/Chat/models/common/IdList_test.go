package common

import (
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_AddId(t *testing.T) {
	aTest := tester.New(t)

	var il *IdList
	var err error

	x := IdList([]ObjectId{1, 2, 3, 4, 5})
	il = &x

	// Test #1.
	err = il.AddId(6)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(([]ObjectId)(*il), []ObjectId{1, 2, 3, 4, 5, 6})

	// Test #2.
	err = il.AddId(3)
	aTest.MustBeAnError(err)
}

func Test_RemoveId(t *testing.T) {
	aTest := tester.New(t)

	var il *IdList
	var err error

	x := IdList([]ObjectId{1, 2, 3, 4, 5})
	il = &x

	// Test #1.
	err = il.RemoveId(4)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(([]ObjectId)(*il), []ObjectId{1, 2, 3, 5})

	// Test #2.
	err = il.RemoveId(10)
	aTest.MustBeAnError(err)
}
