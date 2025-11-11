package encryption

type EncryptionResponse struct {
	Encryption string `json:"encryption"`
	Algorithm  string `json:"algorithm,omitempty"`
}
