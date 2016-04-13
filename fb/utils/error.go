package utils


func CatchPanic() interface{}{
	if err := recover(); err != nil {
		return err
	}
	return nil
}
