package core

import (
	"fmt"
	"sync"
)

/*
一个AOI地图中的格子类型
*/
type Grid struct {
	// 格子的ID
	GID int
	// 格子左边坐标
	MinX int
	// 格子右边坐标
	MaxX int
	// 格子上面坐标
	MinY int
	// 格子下面坐标
	MaxY int
	// 当前玩家或物品的集合
	playerIDs map[int]bool
	// 保护当前集合的锁
	pIDLock sync.RWMutex
}

// 初始化当前格子的方法
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// 给格子添加一个玩家
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	g.playerIDs[playerID] = true
}

// 给格子删除玩家
func (g *Grid) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	delete(g.playerIDs, playerID)
}

// 得到当前格子初始玩家的ID
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}

	return
}

// 打印当前格子全部信息
func (g *Grid) String() string {
	return fmt.Sprintf("Gird id: %d, minX: %d, maxX: %d, minY: %d, maxY: %d, playerIDs: %v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
