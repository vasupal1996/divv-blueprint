package internals

func CreateNewApp(doneCh chan struct{}) App {
	a := AppImpl{
		DoneCh: doneCh,
	}
	return &a
}
