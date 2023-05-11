package cmd

import (
	"os"
	"testing"
)

func TestCompatUTF8_when_read_utf8_bytes(t *testing.T) {
	f, err := os.ReadFile("test4utf8.txt")
	if err != nil {
		t.Error(err)
	}
	bo := compatUTF8(f)

	if !bo {
		t.Error("test4utf8.txt is compat with UTF8")
	}
}

func TestCompatUTF8_when_read_gbk_bytes(t *testing.T) {
	f, err := os.ReadFile("test4gbk.txt")
	if err != nil {
		t.Error(err)
	}
	bo := compatUTF8(f)

	if bo {
		t.Error("test4gbk.txt isn't compat with UTF8")
	}
}
