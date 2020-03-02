package injgo

// Inject is the function inject dynamic library to a process
//
// In windows, name is a file with dll extion.If the file
// name exits, we will return error.
// The workflow of injection in windows is:
// 0. load kernel32.dll in current process.
// 1. open target process T.
// 2. malloc memory in T to store the name of the library.
// 3. get address of function LoadLibraryA from kernel32.dll
//    in T.
// 4. call CreateRemoteThread method in kernel32.dll to execute
//    LoadLibraryA in T.
func Inject(pid int, name string) error {

	// TODO:

	return nil
}
