package model

type PlayerRequest struct {
    URIs []string `json:"uris"`
}

func MapCreatePlayerRequest(uri string) (PlayerRequest) {
    var pr PlayerRequest
    uris := []string{uri}
    pr.URIs = uris

    return pr
}
