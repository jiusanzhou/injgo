package main

import (
	"fmt"
	"syscall"

	"github.com/tadvi/winc"
	"github.com/tadvi/winc/w32"
)

var (
	libuxtheme     uintptr
	setWindowTheme uintptr
)

func main() {
	mainWindow := winc.NewForm(nil)

	mainWindow.SetAndClearStyleBits(0, w32.WS_EX_CONTROLPARENT|w32.WS_EX_APPWINDOW|w32.WS_OVERLAPPEDWINDOW) //
	mainWindow.SetSize(400, 300)                                                                            // (width, height)
	mainWindow.SetText("效能平台")
	mainWindow.Center()

	// go func() {
	// 	for {
	// 		time.Sleep(5 * time.Second)
	// 		w32.SetForegroundWindow(mainWindow.Handle())
	// 	}
	// }()

	btn := winc.NewPushButton(mainWindow)
	btn.SetText("退出")
	btn.SetPos(0, 0)
	btn.SetSize(100, 40)
	btn.OnClick().Bind(func(e *winc.Event) {
		mainWindow.Close()
		winc.Exit()
	})

	mainWindow.Show()

	mainWindow.OnLBDown().Bind(func(e *winc.Event) {
		w32.ReleaseCapture()
		w32.SendMessage(mainWindow.Handle(), uint32(WM_SYSCOMMAND), uintptr(SC_MOVE|HTCAPTION), 0)
	})

	// find the first window
	w, _ := FindWindow("WeChat") // WeChat
	fmt.Println("target window", w.thread)

	// w.Enable()
	// w.Show()

	// 设置为前景
	w.SetForeground()

	fn := func(hook syscall.Handle, event uint32, hwnd syscall.Handle, idObject int32, idChild int32, dwEventThread uint32, dwmsEventTime uint32) syscall.Handle {
		fmt.Println("=-===>")
		return 1
	}

	// 位置和大小变动
	_, err := w.SetWinEventHook(func(evt *Event) error {
		fmt.Println("位置调整")
		return nil
	}, EVENT_SYSTEM_MOVESIZESTART)
	if err != nil {
		fmt.Println("===> set win event hook error:", err)
	}

	// 前景变动
	_, err = w.SetWinEventHook(func(evt *Event) error {
		fmt.Println("窗口激活")
		return nil
	}, EVENT_SYSTEM_FOREGROUND)
	if err != nil {
		fmt.Println("===> set win event hook error:", err)
	}

	defer w.UnhookWinEvent(fn)

	// add the hook resizer and postion chagne hook

	// ws, err := ListWindows()
	// if err != nil {
	// 	fmt.Println("list windows error", err)
	// }
	// for _, w := range ws {
	// 	fmt.Println(">===", w.title, w.hwnd)
	// }

	// mainWindow.OnSize().Bind(func(arg *winc.Event) {
	// 	fmt.Println("====> size", arg)
	// })

	// mainWindow.OnMouseMove().Bind(func(arg *winc.Event) {
	// 	fmt.Println("====> mouse move", arg)
	// })
	// mainWindow.OnPaint().Bind(func(arg *winc.Event) {
	// 	fmt.Println("====> paint", arg)
	// })

	mainWindow.OnClose().Bind(func(arg *winc.Event) {
		winc.Exit()
	})

	winc.RunMainLoop() // Must call to start event loop.
}
