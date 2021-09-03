package exiftool

var (
	JpegTags = []string{
		//Comments
		"XPComment", "Comment",
		//keywords
		"Subject", "LastKeywordIPTC", "LastKeywordXMP", "Keywords", "XPKeywords",
		//Author
		"XPAuthor", "Creator", "By-line", "Artist",

		//GPS
		"GPSLongitude", "GPSLatitude", "GPSAltitude",
	}
	Mp4Tags = []string{
		//Comments
		"Comment",
		//Tags
		"Category",
	}
	MovTags = []string{
		//Comments
		"Comment",
		//Tags
		"Category",
	}
)
