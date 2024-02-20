package command

type StartsClass struct {
	Subject string
}

type PrepareTest struct {
	Subject string
}

type StartTest struct {
	Subject string
}

type SubmitTest struct {
	Subject string
	Name    string
}

type FinishTest struct {
	Subject string
}
