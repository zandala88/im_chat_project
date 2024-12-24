package test

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func Test_Time(t *testing.T) {
	t1 := timestamppb.New(time.Now())
	t.Log(t1)
	t2 := t1.AsTime()
	t.Log(t2)
}
