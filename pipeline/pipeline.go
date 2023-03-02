package pipeline

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	result := make(Bi)

	// build pipeline
	for _, stage := range stages {
		in = stage(in)
	}

	go func() {
		defer close(result)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				result <- v
			}
		}
	}()

	return result
}
