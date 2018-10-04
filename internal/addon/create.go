package addon

import (
	"github.com/kubernauts/tk8/internal"
)

func PrepareExample(addonName string) {
	cloneExample(addonName)
	common.ReplaceGit("./addons/" + addonName)
}
