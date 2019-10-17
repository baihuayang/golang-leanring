package snowflake

import (
	"fmt"
	"testing"
)

func TestSnowFlake(t *testing.T) {
	s := &SnowFlake{}
	s.Init(42)
	// for i := 0; i < 100; i++ {
	// 	go func() {
	// 		println(s.UUID())
	// 	}()
	// }
	uuid := s.UUID()
	uuid = s.UUID()
	fmt.Println(uuid)
	fmt.Println(ParseUUID(uuid))
}
