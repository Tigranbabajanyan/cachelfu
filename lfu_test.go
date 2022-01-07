package cachelfu

import (
	"fmt"
	"testing"
)

const ballotX = "\u2717"
const checkMark = "\u2713"

func TestLfuAdd(t *testing.T) {
	t.Log("Given the need to test add item to cache.")

	lfu := New(10)
	itemKey := "a"
	itemVal := 1

	{
		t.Logf("\tWhen checking to add an element with key \"%s\"", itemKey)

		err := lfu.Add(itemKey, itemVal)
		if err != nil {
			t.Fatal("\t\tShould be able to add a value successfully", ballotX, err)
		}
		t.Log("\t\tShould be able to add a value successfully", checkMark)
	}
}

func TestLfuGet(t *testing.T) {
	t.Log("Given the need to test get item from cache by key.")

	lfu := New(10)
	itemKey := "a"
	itemVal := 1

	err := lfu.Add(itemKey, itemVal)
	if err != nil {
		t.Fatal("\tShould be able to add a value successfully for further retrieve", ballotX, err)
	}

	{
		t.Logf("\tWhen checking to get an element with key \"%s\"", itemKey)

		el, err := lfu.Get(itemKey)
		if err != nil {
			t.Fatal("\t\tShould be able to get a value successfully", ballotX, err)
		}
		t.Log("\t\tShould be able to get a value successfully", checkMark)

		if el == 1 {
			t.Log("\t\tShould return the same value that added for the test", checkMark)
		} else {
			t.Errorf("\t\tShould return the value %d that added for the test. %v %v", itemVal, ballotX, el)
		}
	}
}

func TestLfuGetNonExisting(t *testing.T) {
	t.Log("Given the need to test get non existing item from cache by key.")

	lfu := New(10)
	itemKey := "a"

	{
		t.Logf("\tWhen checking to get a non existing element with key \"%s\"", itemKey)

		_, err := lfu.Get(itemKey)
		if err == nil || err.Error() != fmt.Sprintf("Key \"%s\" is not found", itemKey) {
			t.Error("\t\tShould be return \"key is not found\" error", ballotX, err)
		} else {
			t.Log("\t\tShould be return \"key is not found\" error", checkMark)
		}
	}
}

func TestLfuRemove(t *testing.T) {
	t.Log("Given the need to test remove item from cache by key.")

	lfu := New(10)
	itemKey := "a"
	itemVal := 1

	err := lfu.Add(itemKey, itemVal)
	if err != nil {
		t.Fatal("\tShould be able to add a value successfully for further testing", ballotX, err)
	}

	{
		t.Logf("\tWhen checking to remove an element with key \"%s\"", itemKey)

		err := lfu.Remove(itemKey)
		if err != nil {
			t.Fatal("\t\tShould be able to remove a value successfully", ballotX, err)
		}
		t.Log("\t\tShould be able to remove a value successfully", checkMark)

		_, err = lfu.Get(itemKey)
		t.Log("\tWhen trying to get element after removing")
		if err == nil || err.Error() != fmt.Sprintf("Key \"%s\" is not found", itemKey) {
			t.Error("\t\tShould be return \"key is not found\" error", ballotX, err)
		} else {
			t.Log("\t\tShould be return \"key is not found\" error", checkMark)
		}
	}
}

func TestLfuRemoveNonExisting(t *testing.T) {
	t.Log("Given the need to test remove non existing item from cache by key.")

	lfu := New(10)
	itemKey := "a"

	{
		t.Logf("\tWhen checking to remove a non existing element with key \"%s\"", itemKey)

		err := lfu.Remove(itemKey)
		if err == nil || err.Error() != fmt.Sprintf("Key \"%s\" is not found", itemKey) {
			t.Error("\t\tShould be return \"key is not found\" error", ballotX, err)
		} else {
			t.Log("\t\tShould be return \"key is not found\" error", checkMark)
		}
	}
}
