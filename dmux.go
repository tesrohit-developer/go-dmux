package main

import (
	"fmt"
	"github.com/flipkart-incubator/go-dmux/sideline"
)

//

type DmuxCustom struct {
}

func (d *DmuxCustom) DmuxStart(checkMessageSideline sideline.CheckMessageSideline) {
	fmt.Println(checkMessageSideline.SidelineMessage())
}
