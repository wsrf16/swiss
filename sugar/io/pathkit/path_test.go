package pathkit

import "testing"

func TestPath(t *testing.T) {
	t.Log(GetWorkDirectory())
	t.Log(GetRecursiveFileList("./"))
	t.Log(GetPWD("./"))
	t.Log(Join("//root///", "//a", "///b", "c", "d"))
	t.Log(JoinBy("/", "//root///", "//a", "///b", "c", "d"))
}
