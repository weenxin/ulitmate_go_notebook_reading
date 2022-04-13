package others_test

import (
	"os"
	"testing"
)

func TestEnvironment(t *testing.T) {
	info, err := os.Stat("./testdata/test.file")
	if err != nil {
		t.Fatalf("can not find file")
	}
	t.Logf("file name is : %s \n", info.Name())
	pwd, _ := os.Getwd()
	t.Logf("Testing environment:%s \n", pwd)
}
