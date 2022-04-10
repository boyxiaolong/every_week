package module

import (
	"public/message/msgtype"
	"public/message/protomsg"
)

// "public/command"

// GEpisodeTest 测试用例对象
var GEpisodeTest *EpisodeTest

func init() {
	GEpisodeTest = &EpisodeTest{}
	GEpisodeTest.InitCmds("episode")
	GEpisodeTest.RegCommand("record", recordEpisode)
}

// EpisodeTest 测试用例
type EpisodeTest struct {
	TestBase
}

// 新手引导记录
func recordEpisode() (bool, int32) {
	request := &protomsg.MsgCL2GSEpisodeRecordRequest{}
	request.GuideId = 24

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLEpisodeRecordReply).(*protomsg.MsgGS2CLEpisodeRecordReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}
