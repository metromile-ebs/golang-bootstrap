package datamodel

// {
//     "id": uint,
//     "version": uint,
//     "entrypoint": uuid,
//     "nodes": {
//         [uuid]: {
//             "type": "rest"
//             "branches": [{
//                 "successNode": uuid
//             }],
//             "args" :{
//                 "url": string,
//                 "method": enum(GET, POST, PUT, PATCH, DELETE),
//                 "body": any,
//                 "output": string
//             }
//         },
//         [uuid]: {
//             "type": "set",
//             "branches": [{
//                 "successNode": uuid
//             }],
//             "args": {
//                "key": string
//                "value": any
//             }
//         }
//     }
// }

type Graph struct {
	Id         uint            `json:"id"`
	Version    uint            `json:"version"`
	Entrypoint string          `json:"entrypoint"`
	Nodes      map[string]UUID `json:"nodes"`
}

type UUID struct {
	Type     string   `json:"type"`
	Branches []Branch `json:"branches"`
	Args     Args     `json:"args"`
}

type Branch struct {
	SuccessNode string `json:"successNode"`
}

type Args struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Url    string `json:"url"`
	Method string `json:"method"`
	Body   string `json:"body"`
	Output string `json:"output"`
}
