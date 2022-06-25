package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	res := in

	for _, stage := range stages {
		chain := make(Bi)

		go func(input In, output Bi) {
			defer close(chain)

			for {
				select {
				case x, ok := <-input:
					if !ok {
						return
					}
					output <- x

				case <-done:
					return
				}
			}
		}(res, chain)

		res = stage(chain)
	}

	return res
}
