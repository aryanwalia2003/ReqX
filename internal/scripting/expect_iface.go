package scripting

// ExpectBuilder defines the assertion chain available in JS (e.g. pm.expect(x).To.Eql(y))
type ExpectBuilder interface {
	ToEql(expected interface{}) //yeh check karega ki value equal hai ya nahi
	ToBe(expected interface{}) //yeh check karega ki value same hai ya nahi
	ToExist() //yeh check karega ki value exist karti hai ya nahi
}
