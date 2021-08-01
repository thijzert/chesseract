package notimplemented

import (
	"fmt"
	"path"
	"runtime"
)

func Error() error {
	pc, fileName, fileLine, ok := runtime.Caller(1)
	if ok {
		fun := runtime.FuncForPC(pc)
		if fun != nil {
			fileName = path.Base(fileName)
			return fmt.Errorf("function '%s' (%s:%d) is not implemented", fun.Name(), fileName, fileLine)
		} else {
			return fmt.Errorf("function at %s:%d is not implemented", fileName, fileLine)
		}
	} else {
		return fmt.Errorf("not implemented")
	}
}
