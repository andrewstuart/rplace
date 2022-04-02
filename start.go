package rplace

var start = startMessage{
	ID:   "1",
	Type: "start",
	Payload: startMessagePayload{
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
		Variables: startMessagePayloadVariables{
			Input: startMessagePayloadVariablesInput{
				Channel: startMessagePayloadVariablesInputChannel{
					TeamOwner: "AFD2022",
					Category:  "CANVAS",
					Tag:       "0",
				},
			},
		},
	},
}

type connectionInitMessage struct {
	Type    string                       `json:"type"`
	Payload connectionInitMessagePayload `json:"payload"`
}

type connectionInitMessagePayload struct {
	Authorization string
}

type startMessage struct {
	ID      string              `json:"id"`
	Type    string              `json:"type"`
	Payload startMessagePayload `json:"payload"`
}

type startMessagePayload struct {
	Extensions    struct{}                     `json:"extensions"`
	OperationName string                       `json:"operationName"`
	Query         string                       `json:"query"`
	Variables     startMessagePayloadVariables `json:"variables"`
}

type startMessagePayloadVariables struct {
	Input startMessagePayloadVariablesInput `json:"input"`
}

type startMessagePayloadVariablesInput struct {
	Channel startMessagePayloadVariablesInputChannel `json:"channel"`
}

type startMessagePayloadVariablesInputChannel struct {
	TeamOwner string `json:"teamOwner"`
	Category  string `json:"category"`
	Tag       string `json:"tag"`
}
