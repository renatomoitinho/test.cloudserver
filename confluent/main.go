type ErrorResponse struct {
	Errors []ErroDetails `json:"errors"`
}

type Message struct {
	Data []struct {
		Metadata struct {
			Self string `json:"self"`
		} `json:"metadata"`
	} `json:"data"`
}

type SingleMessage struct {
	Data []struct {
		Spec struct {
			DisplayName           string `json:"display_name"`
			KafkaBoostrapEndpoint string `json:"kafka_boostrap_endpoint"`
			Region                string `json:"region"`
		} `json:"spec"`
	} `json:"data"`
}

func main() {
	var (
		err         error
		msg         *Message
		smsg        *SingleMessage
		errRes      ErrorResponse
		ApiKey      = "your"
		AccessToken = "your"
		ApiUrl      = "http://localhost:40097/envrionments"
	)

	msg, err = confluent.Get[Message](ApiUrl, confluent.WithBasicAuthorization(ApiKey, AccessToken))
	if err := confluent.HasError(err); err != nil {
		fmt.Printf("Error: %v\n%v\n", err.Status, err.Err)
		err.Unmarshal(&errRes)
		fmt.Printf("Error: %v\n", errRes)
		return
	}
	fmt.Println("Success 1")
	fmt.Println(msg)

	smsg, err = confluent.Get[SingleMessage](msg.Data[0].Metadata.Self, confluent.WithBasicAuthorization(ApiKey, AccessToken))
	if err := confluent.HasError(err); err != nil {
		fmt.Printf("Error: %v\n%v\n", err.Status, err.Err)
		err.Unmarshal(&errRes)
		fmt.Printf("Error: %v\n", errRes)
		return
	}

	fmt.Println("Success 2")
	fmt.Println(smsg)
}
