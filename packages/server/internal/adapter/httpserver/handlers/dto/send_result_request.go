package dto

type ResultRequest struct {
	Song1ID    string `json:"song1Id"`
	ComparedTo []struct {
		Song2ID string `json:"song2Id"`
		Result  int    `json:"result"`
	} `json:"comparedTo"`
}

type SendResultRequest struct {
	Results []ResultRequest `json:"results"`
}
