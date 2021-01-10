package types

type CursorImage struct {
	Width  uint16
	Height uint16
	Xhot   uint16
	Yhot   uint16
	Serial uint64
	Pixels []byte
}

type ScreenSize struct {
	Width  int
	Height int
	Rate   int16
}

type ScreenConfiguration struct {
	Width  int
	Height int
	Rates  map[int]int16
}

type DesktopManager interface {
	Start()
	Shutdown() error
	OnBeforeScreenSizeChange(listener func())
	OnAfterScreenSizeChange(listener func())

	// xorg
	ChangeScreenSize(width int, height int, rate int) error
	Move(x, y int)
	Scroll(x, y int)
	ButtonDown(code int) error
	KeyDown(code uint64) error
	ButtonUp(code int) error
	KeyUp(code uint64) error
	ResetKeys()
	ScreenConfigurations() map[int]ScreenConfiguration
	GetScreenSize() *ScreenSize
	SetKeyboardLayout(layout string)
	SetKeyboardModifiers(NumLock int, CapsLock int, ScrollLock int)
	GetCursorImage() *CursorImage

	// xevent
	OnCursorChanged(listener func(serial uint64))
	OnEventError(listener func(error_code uint8, message string, request_code uint8, minor_code uint8))

	// clipboard
	ReadClipboard() string
	WriteClipboard(data string)

	// drop
	DropFiles(x int, y int, files []string)
}
