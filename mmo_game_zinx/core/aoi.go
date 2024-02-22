package core

import "fmt"

/*
AOI区域管理模块
*/

type AOIManager struct {
	// 区域的左边界坐标
	MinX int
	// 区域右边界坐标
	MaxX int
	// X 方向格子的数量
	CntsX int
	// 区域上边界的坐标
	MinY int
	// 区域下边界的坐标
	MaxY int
	// Y 方向格子的数量
	CntsY int
	// 当前区域中有哪些格子
	// Map key = 格子ID， value = 格子对象
	grids map[int]*Grid
}

// 初始化AOI区域管理模块
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	// 给AOI初始化区域所有的格子进行编号和初始化
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// 根据x和y编号计算格子的ID
			// 格子编号 id = idY * cntX + idx
			gid := y*cntsX + x
			// 初始化gid格子
			aoiMgr.grids[gid] = NewGrid(
				gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength(),
			)
		}
	}
	return aoiMgr
}

// 得到每个格子在X轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// 得到每个格子在Y轴方向的高度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

// 打印格子的信息
func (m *AOIManager) String() string {
	// 打印AOIManager信息
	s := fmt.Sprintf("AOIManager:\n MinX: %d, MaxX: %d, cntsX: %d, minY: %d, maxY: %d, cntsY: %d\n Grids in AOIManager:\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	// 打印全部格子信息
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// 根据格子的GID得到周边九宫格格子集合
func (m *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	// 判断gID是否在AIOManager中
	if _, ok := m.grids[gID]; !ok {
		return
	}
	// 初始化grids返回值切片，将当前gid本身加入到九宫格中
	grids = append(grids, m.grids[gID])
	// 需要通过gID得到当前格子x轴的编号 nx = id % nx
	idx := gID % m.CntsX
	// gID的左边是否有格子 || 右边是否有格子 --放到gidsX中
	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}
	if idx < m.CntsX-1 {
		grids = append(grids, m.grids[gID+1])
	}
	// 遍历giddX集合中每个格子的gid
	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}
	for _, v := range gidsX {
		// 得到当前格子id的Y轴编号
		idy := v / m.CntsY
		// 上面是否有格子
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CntsX])
		}
		// 下面是否有格子
		if idy < m.CntsY-1 {
			grids = append(grids, m.grids[v+m.CntsX])
		}
	}
	return
}

// 通过横纵轴坐标得到当前格子编号
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) / m.MinY) / m.gridLength()
	return idy*m.CntsX + idx
}

// 通过横纵坐标，得到周边九宫格内全部的playerIDs
func (m *AOIManager) GetPidsByPos(x, y float32) (playerIDs []int) {
	// 得到当前玩家的GID格子
	gId := m.GetGidByPos(x, y)
	// 通过GID的得到周边九宫格信息
	grids := m.GetSurroundGridsByGid(gId)
	// 将九宫格的信息里的全部player累加到playerIDs
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
		fmt.Sprintf("====> grid ID: %d, pids: %v ====", grid.GID, grid.GetPlayerIDs())
	}
	return
}

// 添加一个playerID到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

// 移除一个格子中的playerID
func (m *AOIManager) RemovePidFromGid(pID, gID int) {
	m.grids[gID].Remove(pID)
}

// 通过GID获得全部的PlayerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = m.grids[gID].GetPlayerIDs()
	return
}

// 通过坐标将Player添加到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.Add(pID)
}

// 通过坐标将Player从一个格子中删除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.Remove(pID)
}
