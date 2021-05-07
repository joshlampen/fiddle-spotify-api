package json

import (
	"encoding/json"
	"net/http"

	"github.com/JoshLampen/fiddle/spotify-api/internal/utils/logger"
)

func WriteResponse(w http.ResponseWriter, resp interface{}) {
    logger := logger.NewLogger()

    w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "content-type")
    jsonBody, err := json.Marshal(resp)
	if err != nil {
        logger.Error().Err(err).Msg("json.WriteResponse - failed to marshal response body")
		return
	}
    w.Write(jsonBody)
}

func WriteError(w http.ResponseWriter, err error, code int) {
    // w.Header().Set("Content-Type", "application/json; charset=utf-8")
    // w.Header().Set("X-Content-Type-Options", "nosniff")
    jsonErr := struct {
        Text string `json:"text"`
    }{
        Text: err.Error(),
    }
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(jsonErr)
}
