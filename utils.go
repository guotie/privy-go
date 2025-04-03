package privy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func unmarshalResp[T any](resp *http.Response, v *T) error {
	// if resp.StatusCode != 200 {
	// 	return fmt.Errorf("invalid response status: %d", resp.StatusCode)
	// }

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// fmt.Println("response: " + string(body))
	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid response status: %d, message: %v", resp.StatusCode, string(body))
	}
	return json.Unmarshal(body, v)
}
