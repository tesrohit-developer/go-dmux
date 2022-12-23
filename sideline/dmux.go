package sideline

import (
	"fmt"
)

//

type DmuxCustom struct {
}

func (d *DmuxCustom) DmuxStart(checkMessageSideline CheckMessageSideline) {
	fmt.Println(checkMessageSideline.SidelineMessage())
}
