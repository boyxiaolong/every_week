package module

import (
	"math/rand"
	"public/message/msgtype"
	"time"

	//"fmt"
	"consoleclient/data"
	"public/command"
	"public/common"
	"public/message/protomsg"
	"strconv"
)

// Init 初始化
func Init() {
}

// Release 释放资源
func Release() {

}

// ClearGuildInfo comment
func ClearGuildInfo() {
	data.GDataCenter.SetData(data.ClearGuildInfo, 0)
}

// GetGuildID comment
func GetGuildID() uint64 {
	return data.GetGuildID()
}

// GetGuildName 获取公会名字
func GetGuildName() string {
	return data.GetGuildName()
}

// SetGuildName comment
func SetGuildName(name string) {
	data.GDataCenter.SetData(data.SetGuildName, name)
}

// SetGuildID comment
func SetGuildID(id uint64) {
	data.GDataCenter.SetData(data.SetGuildID, id)
}

// SetGuildPosition comment
func SetGuildPosition(position protomsg.GuildPosition) {
	data.GDataCenter.SetData(data.SetGuildPosition, position)
}

// CreateGuild comment
func CreateGuild(name string) bool {
	command.GCommand.ExecuteCommand("pm rich")

	rand.Seed(time.Now().Unix())
	command.GCommand.ExecuteCommand("pm createguild " + name + "t" + strconv.FormatUint(rand.Uint64()%100, 10) + strconv.FormatUint(GGameInfo.Account%100, 10))

	waitmsg := GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLGuildCreateReply)
	if waitmsg == nil {
		return false
	}
	reply := waitmsg.(*protomsg.MsgGS2CLGuildCreateReply)
	if reply.ErrorCode != 0 {
		common.GStdout.Error("CreateGuild error %v ", reply.ErrorCode)
		return false
	}

	data.GDataCenter.SetData(data.SetGuildID, reply.GuildId)

	return true
}

// JoinGuild comment
func JoinGuild(guildID uint64) bool {
	request := &protomsg.MsgCL2GSGuildJoinRequest{}
	request.GuildId = guildID

	waitmsg := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildJoinReply)
	if waitmsg == nil {
		return false
	}
	reply := waitmsg.(*protomsg.MsgGS2CLGuildJoinReply)
	if reply.ErrorCode != 0 {
		common.GStdout.Error("JoinGuild error %v ", reply.ErrorCode)
		return false
	}

	return true
}

// ForcePlayerJoinGuild 强制拉其他人进指定公会
func ForcePlayerJoinGuild(guildID uint64, playerID uint64) bool {

	ForcePlayerQuitGuild(playerID)

	cmd := "pm joinguild " + strconv.FormatUint(guildID, 10) + " " + strconv.FormatUint(playerID, 10)
	command.GCommand.ExecuteCommand(cmd)
	waitmsg := GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLGuildJoinReply)
	if waitmsg == nil {
		return false
	}

	reply := waitmsg.(*protomsg.MsgGS2CLGuildJoinReply)
	if reply.ErrorCode != 0 {
		common.GStdout.Error("Force player %v JoinGuild error %v ", playerID, reply.ErrorCode)
		return false
	}

	return true
}

// ForcePlayerQuitGuild 强制他人退出公会
func ForcePlayerQuitGuild(playerID uint64) bool {

	cmd := "pm quitguild " + strconv.FormatUint(playerID, 10)
	command.GCommand.ExecuteCommand(cmd)

	return true
}

// QuitGuild comment
func QuitGuild() {

	if GetGuildID() == 0 {
		return
	}

	command.GCommand.ExecuteCommand("pm quitguild")
	GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLGuildQuitReply)
}

// DestroyGuild comment
func DestroyGuild() {
	command.GCommand.ExecuteCommand("pm dismiss")
	GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLGuildDismissReply)
}

// ChangePos comment
func ChangePos(playerID uint64, pos protomsg.GuildPosition) {
	request := &protomsg.MsgCL2GSGuildChangePositionRequest{}
	request.PlayerId = playerID
	request.Position = pos

	GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildChangePositionReply)
}
