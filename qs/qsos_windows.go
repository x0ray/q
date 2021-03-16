// Package qs - q scripting language
package qs

func osStat(L *LState) int {
	L.Push(LNil)
	return 1
}

func osStatfs(L *LState) int {
	L.Push(LNil)
	return 1
}

// TODO
/* like ...
import
 	"os"
 	"runtime"
 	"syscall"

 	"golang.org/x/sys/windows"
 )

 func IsHidden(filename string) (bool, error) {

 	if runtime.GOOS == "windows" {

 		pointer, err := syscall.UTF16PtrFromString(filename)
 		if err != nil {
 			return false, err
 		}
 		attributes, err := syscall.GetFileAttributes(pointer)
 		if err != nil {
 			return false, err
 		}
 		return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0, nil
 	}
 }
*/
