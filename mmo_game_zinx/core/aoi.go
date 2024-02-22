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
	grids = append()
	// 需要通过gID得到当前格子x轴的编号 nx = id % nx

	// gID的左边是否有格子 || 右边是否有格子 --放到gidsX中

	// 遍历giddX集合中每个格子的gid
	// 上面是否有格子
	// 下面是否有格子
}
