package models

type SpotifyImage struct {
	Url    string `json:"url,omitempty"`
	Height int    `json:"height,omitempty"`
	Width  int    `json:"width,omitempty"`
}

type SpotifyFollowers struct {
	Href  string `json:"href,omitempty"`
	Total int    `json:"total,omitempty"`
}

type SpotifyExternalUrls struct {
	Spotify string `json:"spotify,omitempty"`
}

type SpotifyArtist struct {
	ExternalUrls SpotifyExternalUrls `json:"external_urls,omitempty"`
	Name         string              `json:"name,omitempty"`
}

type SpotifyAlbum struct {
	Image []SpotifyImage `json:"images,omitempty"`
	Name  string         `json:"name,omitempty"`
}

type CurrentPlayingItem struct {
	Album        SpotifyAlbum        `json:"album,omitempty"`
	Artists      []SpotifyArtist     `json:"artists,omitempty"`
	Name         string              `json:"name,omitempty"`
	ExternalUrls SpotifyExternalUrls `json:"external_urls,omitempty"`
}

type CurrentPlayingResponse struct {
	Item CurrentPlayingItem `json:"item,omitempty"`
}

type SpotifyTokenResponseBody struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
