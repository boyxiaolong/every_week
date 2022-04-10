package player

import (
	//"public/command"

	"public/common"
	"runtime"
	"sync"
	"time"
)

// GPlayerMgr comment
var GPlayerMgr *Mgr

func init() {
	GPlayerMgr = &Mgr{
		Players:       make(map[uint64]*Player, 0),
		RemovePlayers: make(map[uint64]bool, 0),
	}

	go updateRemovePlayer()
}

func updateRemovePlayer() {
	for {
		//common.GStdout.Console("remove ========================= 1")
		GPlayerMgr.mutexRemovePlayer.Lock()
		if len(GPlayerMgr.RemovePlayers) > 0 {
			for k := range GPlayerMgr.RemovePlayers {
				//common.GStdout.Console("remove ========================= 2 %v", k)
				GPlayerMgr.DelPlayer(k)
			}
			GPlayerMgr.RemovePlayers = make(map[uint64]bool, 0)
		}
		GPlayerMgr.mutexRemovePlayer.Unlock()
		time.Sleep(5 * time.Second)
	}
}

// Mgr comment
type Mgr struct {
	Players           map[uint64]*Player
	mutexPlayer       sync.RWMutex
	RemovePlayers     map[uint64]bool
	mutexRemovePlayer sync.RWMutex
}

func PlayerRelease(player *Player) {
	common.GStdout.Debug("player[%v] release", player.PlayerID)
}

// CreatePlayers comment
func (m *Mgr) CreatePlayers(accountStart uint64, accountEnd uint64) {
	for i := accountStart; i <= accountEnd; i++ {
		player := &Player{}
		runtime.SetFinalizer(player, PlayerRelease)
		player.Init(uint64(i), uint64(i))
		m.AddPlayer(player)
	}
	return
}

// CreatePlayer comment
func (m *Mgr) CreatePlayer(playerid uint64) *Player {
	if player, ok := m.GetPlayer(playerid); ok {
		return player
	}

	player := &Player{}
	runtime.SetFinalizer(player, PlayerRelease)
	player.Init(playerid, uint64(playerid))
	m.AddPlayer(player)
	return player
}

// GetPlayer comment
func (m *Mgr) GetPlayer(playerid uint64) (*Player, bool) {
	m.mutexPlayer.Lock()
	defer m.mutexPlayer.Unlock()

	player, ok := m.Players[playerid]
	return player, ok
}

// AddPlayer comment
func (m *Mgr) AddPlayer(player *Player) {
	m.mutexPlayer.Lock()
	defer m.mutexPlayer.Unlock()

	m.Players[player.PlayerID] = player
}

// DelPlayer comment
func (m *Mgr) DelPlayer(playerID uint64) {
	m.mutexPlayer.Lock()

	//common.GStdout.Console("DelPlayer============================1 %v", playerID)
	player, ok := m.Players[playerID]
	if ok {
		delete(m.Players, playerID)
		//common.GStdout.Console("DelPlayer============================2 %v", playerID)
	}
	m.mutexPlayer.Unlock()
	player.Release()
	//common.GStdout.Console("DelPlayer============================3 %v", playerID)
}

// RemovePlayer comment
func (m *Mgr) RemovePlayer(player *Player) {
	m.mutexRemovePlayer.Lock()
	defer m.mutexRemovePlayer.Unlock()
	//common.GStdout.Console("remove player ========================= 1 %v", player.PlayerID)
	m.RemovePlayers[player.PlayerID] = true
}

func (m *Mgr) GetPlayerCount() int {
	m.mutexPlayer.Lock()
	defer m.mutexPlayer.Unlock()
	return len(m.Players)
}
