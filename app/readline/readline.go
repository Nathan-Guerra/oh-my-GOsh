package readline

/*
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"strconv"
)

func Random() {
	fmt.Println(strconv.FormatInt(int64(C.random()), 10))
}

func Seed(i int) {
	C.srandom(C.uint(i))
}
