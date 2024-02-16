package command

type ClassStarts struct {
	Subject string
}

type TestBegins struct {
	Subject string
}

type SubmitTest struct {
	Subject string
	Name    string
}

type ReceiveTest struct {
	Subject string
}

type EndTest struct {
	Subject string
}
