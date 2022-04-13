package others_test

import "testing"

func BenchmarkChildTest(b *testing.B) {
	sum := 0
	for i := 0; i < b.N; i++ {
		sum += i
	}
	b.Logf("BenchmarkChildTest run : %d times", b.N)

}

func TestChildTest(t *testing.T) {

	t.Logf("TestChildTest running")

	t.Run("hello", func(t *testing.T) {
		t.Logf("TestChildTest/hello running")
	})
	t.Run("hello1", func(t *testing.T) {
		t.Logf("TestChildTest/hello1 running")
	})
	t.Run("hello2", func(t *testing.T) {
		t.Logf("TestChildTest/hello2 running")
	})
}
