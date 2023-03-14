package scripts

import "testing"

func TestGetBranches(t *testing.T) {
	gomss := new(GomssOp)
	branches, err := gomss.GetBranches()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("branches:", branches, branches[0])
	err = gomss.Checkout(branches[2])
	if err != nil {
		t.Fatal(err)
	}

	data, err := gomss.Build()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestMake(t *testing.T) {
	gomss := new(GomssOp)
	data, err := gomss.Build()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestPublish(t *testing.T) {
	gomss := new(GomssOp)
	data, err := gomss.Publish("v2.1.1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}
