package frame

type Post struct {
	UntrustedData UntrustedData `json:"untrustedData"`
}

type UntrustedData struct {
	ButtonIndex int             `json:"buttonIndex"`
	CastID      UntrustedCastID `json:"castId"`
	FID         int             `json:"fid"`
}

type UntrustedCastID struct {
	FID  int    `json:"fid"`
	Hash string `json:"hash"`
}
