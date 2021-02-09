package kudu

import (
	"fmt"
	"testing"
)

func TestConnection(t *testing.T) {
	b := NewClientBuilder()
	defer b.Free()
	b.AddMasterServerAddr("172.16.71.6:7051")
	client, err := b.Build()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Make a table if it doesn't exist
	exists, err := client.TableExists("cpp_test")
	if err != nil {
		panic(err)
	}
	fmt.Println(exists)

	table, err := client.OpenTable("cpp_test")
	if err != nil {
		panic(err)
	}
	defer table.Close()

	session := client.NewSession()
	defer session.Close()

	session.SetFlushMode(AutoFlushBackground)
	for i := 0; i < 10; i++ {
		ins := table.NewInsert()
		if err := ins.SetInt32("c1", int32(i)); err != nil {
			panic(err)
		}
		if err := ins.SetString("c2", "test"); err != nil {
			panic(err)
		}
		if err := session.Apply(ins); err != nil {
			panic(err)
		}
	}
	if err := session.Flush(); err != nil {
		panic(err)
	}

}
