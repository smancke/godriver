package exec

func N(ntimes int, f func()) {
	for i := 0; i < ntimes; i++ {
		f()
	}
}
