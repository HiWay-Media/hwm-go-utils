package redis

import (
	"reflect"
	"testing"
)

func TestInitSkipsUnsetAddresses(t *testing.T) {
	for _, env := range []string{"REDIS_1", "REDIS_2", "REDIS_3", "REDIS_4", "REDIS_5", "REDIS_6"} {
		t.Setenv(env, "")
	}
	t.Setenv("REDIS_1", "10.0.0.1:6379")
	t.Setenv("REDIS_3", "10.0.0.3:6379")

	Init()
	if want := []string{"10.0.0.1:6379", "10.0.0.3:6379"}; !reflect.DeepEqual(Brokers, want) {
		t.Errorf("Brokers = %v, want %v", Brokers, want)
	}
}
