package models

// Giphy API json models.
//
// The Giphy search endpoint (https://api.giphy.com/v1/gifs/search) returns:
//
//	{ "data": [ { "id": "...", "images": { ... } } ], "meta": { ... } }
//
// The Go field is named Results (tagged json:"data") so the rest of the app
// can keep referring to apiData.Results just like it did with the old Tenor
// response shape.
type Api struct {
	Results []Result `json:"data"`
	Meta    Meta     `json:"meta"`
}

type Meta struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

type Result struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Url    string `json:"url"`
	Images Images `json:"images"`
}

// Images holds the subset of Giphy renditions we use. Giphy exposes many more,
// but "original" (the full animated gif) and "original_still" (a static preview
// frame) are all term-gif needs.
type Images struct {
	Original        Rendition `json:"original"`
	OriginalStill   Rendition `json:"original_still"`
	DownsizedMedium Rendition `json:"downsized_medium"`
	PreviewGif      Rendition `json:"preview_gif"`
}

type Rendition struct {
	Url string `json:"url"`
}
