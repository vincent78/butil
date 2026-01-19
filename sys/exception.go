package sys

func Try(userFn func(), catchFn func(err any)) {
	defer func() {
		if err := recover(); err != nil {
			catchFn(err)
		}
	}()
	userFn()
}
