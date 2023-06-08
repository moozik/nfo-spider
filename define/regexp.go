package define

import "regexp"

var PattrenAvName *regexp.Regexp
var PattrenAvFile *regexp.Regexp

func init() {
	PattrenAvName, _ = regexp.Compile(`^[a-zA-Z]{3,5}-\d{3,5}`)
	PattrenAvFile, _ = regexp.Compile(`(.*)\.(mp4|avi|mkv)$`)
}
