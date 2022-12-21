package main

import (
	"fmt"
	"github.com/flipkart-incubator/go-dmux/sideline"
)

//

// **************** Bootstrap ***********

func dmuxStart(checkMessageSideline sideline.CheckMessageSideline) {
	fmt.Println(checkMessageSideline.SidelineMessage())
}
