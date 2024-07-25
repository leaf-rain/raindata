package retcd

import (
	"testing"
	"time"
)

func TestLock_Lock(t *testing.T) {
	lock, err := NewLock(ctx, cli, "test", 10)
	if err != nil {
		t.Fatal(err)
	}
	err = lock.Lock()
	t.Logf("1,%v", err)
	//go func() {
	//	time.Sleep(time.Second * 10)
	//	lock.Close()
	//}()
	time.Sleep(time.Second * 100)
}

func TestLock_Lock2(t *testing.T) {
	lock, err := NewLock(ctx, cli, "test", 10)
	if err != nil {
		t.Fatal(err)
	}
	err = lock.Lock()
	t.Logf("2,%v", err)
	time.Sleep(time.Second * 100)
}
