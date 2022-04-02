package rplace

var start = StartMessage{
	ID:   "1",
	Type: "start",
	Payload: StartMessagePayload{
		Extensions:    struct{}{},
		OperationName: "replace",
		Query: `subscription replace($input: SubscribeInput!) {
	subscribe(input: $input) {
		id
		... on BasicMessage {
			data {
				__typename
				... on FullFrameMessageData {
					__typename
					name
					timestamp
				}
				... on DiffFrameMessageData {
					__typename
					name
					currentTimestamp
					previousTimestamp
				}
			}
			__typename
		}
		__typename
	}
}`,
		Variables: StartMessagePayloadVariables{
			Input: StartMessagePayloadVariablesInput{
				Channel: StartMessagePayloadVariablesInputChannel{
					TeamOwner: "AFD2022",
					Category:  "CANVAS",
					Tag:       "0",
				},
			},
		},
	},
}

type ConnectionInitMessage struct {
	Type    string                       `json:"type"`
	Payload ConnectionInitMessagePayload `json:"payload"`
}

type ConnectionInitMessagePayload struct {
	Authorization string
}

type StartMessage struct {
	ID      string              `json:"id"`
	Type    string              `json:"type"`
	Payload StartMessagePayload `json:"payload"`
}

type StartMessagePayload struct {
	Extensions    struct{}                     `json:"extensions"`
	OperationName string                       `json:"operationName"`
	Query         string                       `json:"query"`
	Variables     StartMessagePayloadVariables `json:"variables"`
}

type StartMessagePayloadVariables struct {
	Input StartMessagePayloadVariablesInput `json:"input"`
}

type StartMessagePayloadVariablesInput struct {
	Channel StartMessagePayloadVariablesInputChannel `json:"channel"`
}

type StartMessagePayloadVariablesInputChannel struct {
	TeamOwner string `json:"teamOwner"`
	Category  string `json:"category"`
	Tag       string `json:"tag"`
}
