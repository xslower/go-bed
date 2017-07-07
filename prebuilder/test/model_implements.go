package test

import (
	"errors"
	"greentea/orm"
	"hash/crc32"
	"strconv"
)

func NewVideoModel() *VideoModel {
	mt := &VideoModel{}
	mt.BaseModel.IModel = mt
	mt.result = &([]*Video{})
	return mt
}

type VideoModel struct {
	orm.BaseModel
	result *[]*Video
}

func (this *VideoModel) DbName() string {
	return "content"
}
func (this *VideoModel) TableName(elems ...orm.ISqlElem) string {
	if len(elems) == 0 {
		return "video_tvplay"
	}
	elem := elems[0]
	ifc, _ := elem.Get("category")
	category, _ := orm.InterfaceToString(ifc)
	elem.Del("category")
	if category == "" {
		return "video_tvplay"
	}
	return "video_" + category
}
func (this *VideoModel) PartitionKey() string {
	return "category"
}

func (this *VideoModel) Result() []*Video {
	return *(this.result)
}
func (this *VideoModel) CreateRow() orm.IRow {
	irow := &Video{}
	*(this.result) = append(*(this.result), irow)
	return irow
}
func (this *VideoModel) ResetResult() {
	*(this.result) = []*Video{}
}
func (this *VideoModel) ToIRows(rows *[]*Video) []orm.IRow {
	this.result = rows
	irows := make([]orm.IRow, len(*rows))
	for i, r := range *rows {
		irows[i] = r
	}
	return irows
}

func (this *Video) Set(key string, val []byte) error {
	var err error
	switch key {
	case "category":
		this.Category = string(val)
	case "video_id":
		this.VideoId, err = strconv.Atoi(string(val))
	case "video_title":
		this.VideoTitle = string(val)
	case "poster":
		this.Poster = string(val)
	case "area":
		this.Area = string(val)
	case "year":
		this.Year, err = strconv.Atoi(string(val))
	case "source":
		this.Source = string(val)
	case "src_video_id":
		this.SrcVideoId, err = strconv.Atoi(string(val))
	case "tags":
		this.Tags = string(val)
	case "desc":
		this.Desc = string(val)

	default:
		err = errors.New("No such column [" + key + "]")
	}
	return err
}

func (this *Video) Get(key string) interface{} {
	switch key {
	case "category":
		return this.Category
	case "video_id":
		return this.VideoId
	case "video_title":
		return this.VideoTitle
	case "poster":
		return this.Poster
	case "area":
		return this.Area
	case "year":
		return this.Year
	case "source":
		return this.Source
	case "src_video_id":
		return this.SrcVideoId
	case "tags":
		return this.Tags
	case "desc":
		return this.Desc

	default:
		return nil
	}
}

func (this *Video) Columns() []string {
	return []string{`category`, `video_id`, `video_title`, `poster`, `area`, `year`, `source`, `src_video_id`, `tags`, `desc`}
}

func NewEpisodeModel() *EpisodeModel {
	mt := &EpisodeModel{}
	mt.BaseModel.IModel = mt
	mt.result = &([]*Episode{})
	return mt
}

type EpisodeModel struct {
	orm.BaseModel
	result *[]*Episode
}

func (this *EpisodeModel) DbName() string {
	return "content"
}
func (this *EpisodeModel) TableName(elems ...orm.ISqlElem) string {
	if len(elems) == 0 {
		return "video_episode_tvplay"
	}
	elem := elems[0]
	ifc, _ := elem.Get("category")
	category, _ := orm.InterfaceToString(ifc)
	elem.Del("category")
	if category == "" {
		return "video_episode_tvplay"
	}
	return "video_episode_" + category
}
func (this *EpisodeModel) PartitionKey() string {
	return "category"
}

func (this *EpisodeModel) Result() []*Episode {
	return *(this.result)
}
func (this *EpisodeModel) CreateRow() orm.IRow {
	irow := &Episode{}
	*(this.result) = append(*(this.result), irow)
	return irow
}
func (this *EpisodeModel) ResetResult() {
	*(this.result) = []*Episode{}
}
func (this *EpisodeModel) ToIRows(rows *[]*Episode) []orm.IRow {
	this.result = rows
	irows := make([]orm.IRow, len(*rows))
	for i, r := range *rows {
		irows[i] = r
	}
	return irows
}

func (this *Episode) Set(key string, val []byte) error {
	var err error
	switch key {
	case "category":
		this.Category = string(val)
	case "episode_id":
		this.EpisodeId, err = strconv.Atoi(string(val))
	case "video_id":
		this.VideoId, err = strconv.Atoi(string(val))
	case "src_video_id":
		this.SrcVideoId, err = strconv.Atoi(string(val))
	case "name":
		this.Name = string(val)
	case "pic":
		this.Pic = string(val)
	case "url":
		this.Url = string(val)

	default:
		err = errors.New("No such column [" + key + "]")
	}
	return err
}

func (this *Episode) Get(key string) interface{} {
	switch key {
	case "category":
		return this.Category
	case "episode_id":
		return this.EpisodeId
	case "video_id":
		return this.VideoId
	case "src_video_id":
		return this.SrcVideoId
	case "name":
		return this.Name
	case "pic":
		return this.Pic
	case "url":
		return this.Url

	default:
		return nil
	}
}

func (this *Episode) Columns() []string {
	return []string{`category`, `episode_id`, `video_id`, `src_video_id`, `name`, `pic`, `url`}
}

func ormStart(dbConfig map[string]string) {
	orm.Start(dbConfig)
	orm.RegisterModel(NewVideoModel())
	orm.RegisterModel(NewEpisodeModel())

}
func packageHolder() {
	_ = crc32.ChecksumIEEE([]byte("a"))
}
