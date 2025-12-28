package job

type (
	Job interface {
		Run()
	}
	job struct {
		Spec string
	}
)
