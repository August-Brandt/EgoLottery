package config

type Config struct {
	GroupType   string   `fig:"grouptype" default:"days"`
	TimeAgo     int      `fig:"timeago" default:"150"`
	SearchDepth int      `fig:"searchdepth" default:"0"`
	Emails      []string `fig:"emails" validate:"required"`
	Directories []string `fig:"directories" default:"."`
}