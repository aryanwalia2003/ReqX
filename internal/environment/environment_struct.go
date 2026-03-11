package environment

// Environment represents a collection of variables to be used during request execution.
type Environment struct {
	Name      string            `json:"name"`
	Variables map[string]string `json:"variables"`
}

//this struct will look like 
// {
// 	"name": "dev",
// 	"variables": {
// 		"base_url": "http://localhost:8080",
// 		"api_key": "1234567890"
// 	}
// }
//this is a json object that will be stored in a file and will be used to store the environment variables
