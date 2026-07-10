package models

// GifResult is the provider-agnostic shape the handler consumes. Each provider
// converts its own response into a slice of these.
type GifResult struct {
	Id         string
	GifUrl     string // animated GIF URL to fetch and render
	PreviewUrl string // static image URL used for ?img=true previews
}

// ---------------------------------------------------------------------------
// GIPHY  —  GET https://api.giphy.com/v1/gifs/search
// Response: { "data": [ { "id", "images": { "original": {"url"}, ... } } ], "meta": {...} }
// ---------------------------------------------------------------------------

type GiphyApi struct {
	Results []GiphyResult `json:"data"`
	Meta    GiphyMeta     `json:"meta"`
}

type GiphyMeta struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

type GiphyResult struct {
	Id     string      `json:"id"`
	Images GiphyImages `json:"images"`
}

type GiphyImages struct {
	Original      GiphyRendition `json:"original"`
	OriginalStill GiphyRendition `json:"original_still"`
}

type GiphyRendition struct {
	Url string `json:"url"`
}

// ToResults normalizes a GIPHY search response.
func (a GiphyApi) ToResults() []GifResult {
	out := make([]GifResult, 0, len(a.Results))
	for _, r := range a.Results {
		if r.Images.Original.Url == "" {
			continue
		}
		out = append(out, GifResult{
			Id:         r.Id,
			GifUrl:     r.Images.Original.Url,
			PreviewUrl: r.Images.OriginalStill.Url,
		})
	}
	return out
}

// ---------------------------------------------------------------------------
// KLIPY  —  GET https://api.klipy.com/api/v1/{key}/gifs/search
// Response: { "result": true, "data": { "data": [
//              { "slug", "title", "file": { "<size>": { "<format>": {"url"} } } } ] } }
// sizes: hd|md|sm|xs ; formats: gif|webp|mp4|jpg
// ---------------------------------------------------------------------------

type KlipyApi struct {
	Result bool          `json:"result"`
	Data   KlipyDataWrap `json:"data"`
}

type KlipyDataWrap struct {
	Data []KlipyResult `json:"data"`
}

type KlipyResult struct {
	Slug  string                               `json:"slug"`
	Title string                               `json:"title"`
	File  map[string]map[string]KlipyRendition `json:"file"`
}

type KlipyRendition struct {
	Url string `json:"url"`
}

// pickKlipy returns the first non-empty URL for the given format, checking the
// size tiers in the provided priority order. It is tolerant of missing tiers.
func pickKlipy(file map[string]map[string]KlipyRendition, sizes []string, format string) string {
	for _, s := range sizes {
		if formats, ok := file[s]; ok {
			if r, ok := formats[format]; ok && r.Url != "" {
				return r.Url
			}
		}
	}
	return ""
}

// ToResults normalizes a KLIPY search response.
func (a KlipyApi) ToResults() []GifResult {
	out := make([]GifResult, 0, len(a.Data.Data))
	for _, r := range a.Data.Data {
		// Prefer a mid-size animated gif for fast terminal rendering.
		gifUrl := pickKlipy(r.File, []string{"md", "sm", "hd", "xs"}, "gif")
		if gifUrl == "" {
			continue
		}
		// Static preview: a small jpg still, falling back to a small gif.
		prev := pickKlipy(r.File, []string{"xs", "sm", "md"}, "jpg")
		if prev == "" {
			prev = pickKlipy(r.File, []string{"xs", "sm", "md"}, "gif")
		}
		out = append(out, GifResult{Id: r.Slug, GifUrl: gifUrl, PreviewUrl: prev})
	}
	return out
}
