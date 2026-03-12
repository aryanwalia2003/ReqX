package collection

// Collection groups requests and scripts together.
type Collection struct {
	Name     string    `json:"name"`
	Auth     *Auth     `json:"auth,omitempty"`
	Requests []Request `json:"requests"`
}

//this struct will look like 
// {
// 	"name": "dev",
// 	"requests": [
// 		{
// 			"name": "request_name",
// 			"method": "GET",
// 			"url": "http://localhost:8080",
// 			"headers": {
// 				"Content-Type": "application/json"
// 			},
// 			"body": "{\"key\":\"value\"}"
// 		}
// 	]
// }
