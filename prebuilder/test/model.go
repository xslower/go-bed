package test

//@db content
//@table T1:video_+tvplay{category}
type Video struct {
	Category   string
	VideoId    int
	VideoTitle string
	Poster     string
	Area       string
	Year       int
	Source     string
	SrcVideoId int
	Tags       string
	Desc       string
}

// //@db content
// //@table T1:video_detail_+tvplay{category}
// type VideoDetail struct {
// 	Category string
// 	VideoId  int
// }

//Episode, the type inner Episodes field.
//@db content
//@table T1:video_episode_+tvplay{category}
type Episode struct {
	Category   string
	EpisodeId  int
	VideoId    int
	SrcVideoId int
	Name       string
	Pic        string
	Url        string
}
