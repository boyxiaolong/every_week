package module

import (
    "public/message/msgtype"
    "public/message/protomsg"
)

// KingdomListTest comment
type KingdomListTest struct {
    TestBase
}

//GKingdomListTest comment
var GKingdomListTest *KingdomListTest

func init() {
    GKingdomListTest = &KingdomListTest{}
    GKingdomListTest.InitCmds("kingdom_list")
    GKingdomListTest.RegCommand("info", testInfo)
}






//testInfo comment
func testInfo() (bool, int32) {
    request := &protomsg.MsgCL2GSKingdomListInfoRequest{}

    
    response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLKingdomListInfoReply).(*protomsg.MsgGS2CLKingdomListInfoReply)
    
    return response.ErrorCode == 0, response.ErrorCode
}




