package logo

import "github.com/wsrf16/swiss/sugar/runtimekit"

type DefaultFieldsHook struct {
	Info runtimekit.StackTraceInfo
}

func (h *DefaultFieldsHook) Fire(entry *Entry) error {
	//info := runtimekit.GetStackTraceInfos()[1]
	//if entry.Level < InfoLevel {
	//    entry.Data["stackoverflow"] = runtimekit.GetStackTraceLines()[1:]
	//    entry.Data["file"] = info.File
	//    entry.Data["line"] = info.Line
	//    entry.Data["func"] = info.Function
	//    entry.Data["pc"] = info.ProgramCounter
	//}
	return nil
}

func (h *DefaultFieldsHook) Levels() []Level {
	return AllLevels
}
