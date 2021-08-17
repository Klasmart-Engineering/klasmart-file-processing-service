package exiftool

var (
	JpegTags = []string{
		//Comments
		"XPComment", "Comment",
		//keywords
		"Subject", "LastKeywordIPTC", "LastKeywordXMP", "Keywords", "XPKeywords",
		//Author
		"XPAuthor", "Creator", "By-line", "Artist",
	}
)
