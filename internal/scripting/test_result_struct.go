package scripting

// TestResult tracks the outcome of a pm.test execution block
type TestResult struct {
	Name    string
	Passed  bool
	Error   string
}

// Ensure the slice of results can easily be shared across the system
type TestResults []TestResult 

//this looks like a slice of test results : [TestResult{Name: "Test 1", Passed: true, Error: ""}, TestResult{Name: "Test 2", Passed: false, Error: "Test 2 failed"}]
