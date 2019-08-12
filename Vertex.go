package catgl

// 顶点类
//   实现顶点相关操作
// ? 日志
// !  2019-8-3 重构
import (
	"errors"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// * 纹理变量
const (
	// 纹理单元
	TEXTURE   = 0x1702
	TEXTURE0  = 0x84C0
	TEXTURE1  = 0x84C1
	TEXTURE10 = 0x84CA
	TEXTURE11 = 0x84CB
	TEXTURE12 = 0x84CC
	TEXTURE13 = 0x84CD
	TEXTURE14 = 0x84CE
	TEXTURE15 = 0x84CF
	TEXTURE16 = 0x84D0
	TEXTURE17 = 0x84D1
	TEXTURE18 = 0x84D2
	TEXTURE19 = 0x84D3
	TEXTURE2  = 0x84C2
	TEXTURE20 = 0x84D4
	TEXTURE21 = 0x84D5
	TEXTURE22 = 0x84D6
	TEXTURE23 = 0x84D7
	TEXTURE24 = 0x84D8
	TEXTURE25 = 0x84D9
	TEXTURE26 = 0x84DA
	TEXTURE27 = 0x84DB
	TEXTURE28 = 0x84DC
	TEXTURE29 = 0x84DD
	TEXTURE3  = 0x84C3
	TEXTURE30 = 0x84DE
	TEXTURE31 = 0x84DF
	TEXTURE4  = 0x84C4
	TEXTURE5  = 0x84C5
	TEXTURE6  = 0x84C6
	TEXTURE7  = 0x84C7
	TEXTURE8  = 0x84C8
	TEXTURE9  = 0x84C9
	// 纹理类型
	TEXTURE1D                 = 0x0DE0
	TEXTURE1DARRAY            = 0x8C18
	TEXTURE2D                 = 0x0DE1
	TEXTURE2DARRAY            = 0x8C1A
	TEXTURE2DMULTISAMPLE      = 0x9100
	TEXTURE2DMULTISAMPLEARRAY = 0x9102
	TEXTURE3D                 = 0x806F
)

// Vertex 顶点类
type Vertex struct {
	VAO    uint32
	Buffer uint32
	// 显示模式
	DisplayMode uint32
	// 坐标
	Position mgl32.Mat4
	// 标记
	ifCreate bool
	ifIndex  bool
	// 索引信息
	indexN   int32
	indexIbo uint32
}

// SetVertex 设置顶点
func (V *Vertex) SetVertex(
	vertices []float32, // 位置
	normals []float32, // 法线
	uv []float32, // 纹理
) error {
	// 保证参数
	if vertices == nil {
		return errors.New("顶点不能为空")
	}
	// 销毁
	V.Delete()
	V.ifCreate = true
	// 设置显示模式
	V.DisplayMode = gl.TRIANGLES
	// 创建 VAO
	gl.GenVertexArrays(1, &(V.VAO))
	// 绑定VAO
	gl.BindVertexArray(V.VAO)
	//!..............创建缓存..............!\\
	// 设置顶点缓存
	gl.GenBuffers(1, &(V.Buffer))
	gl.BindBuffer(gl.ARRAY_BUFFER, V.Buffer)
	// 获得数据大小
	var p, n, t int
	p = 4 * len(vertices)
	if normals != nil {
		n = 4 * len(normals)
	}
	if uv != nil {
		t = 4 * len(uv)
	}
	// 预分配空间
	gl.BufferData(gl.ARRAY_BUFFER, p+n+t, nil, gl.STATIC_DRAW)
	// 设置数据结构
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, p, gl.Ptr(vertices))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, gl.PtrOffset(0))
	//? 添加顶点法线
	if normals != nil {
		gl.BufferSubData(gl.ARRAY_BUFFER, p, n, gl.Ptr(normals))
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 12, gl.PtrOffset(p))
	}
	//? 设置订顶点 UV
	if uv != nil {
		gl.BufferSubData(gl.ARRAY_BUFFER, p+n, t, gl.Ptr(uv))
		gl.EnableVertexAttribArray(2)
		gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8, gl.PtrOffset(p+n))
	}
	// 设置数量
	V.indexN = int32(len(vertices))
	// 完成设置
	gl.BindVertexArray(0)
	return nil
}

// SetIndex 设置索引缓冲
func (V *Vertex) SetIndex(
	indices []uint32, // 索引
) {
	if V.ifCreate {
		if V.ifCreate {
			gl.DeleteBuffers(0, &(V.indexIbo))
		}
		V.ifIndex = true
		// 设置顶点
		gl.BindVertexArray(V.VAO)
		gl.GenBuffers(1, &(V.indexIbo))
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, V.indexIbo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
		gl.BindVertexArray(0)
		// 设置数量
		V.indexN = int32(len(indices))
	}
}

// SetTexture 设置材质
func (V *Vertex) SetTexture(
	file string, // 材质文件名
	unit uint32, // 纹理单元
	Target uint32, // 纹理类型
) error {
	texture, err := NewTexture(file, Target)
	if err != nil {
		return err
	}
	gl.BindVertexArray(V.VAO) // 绘画
	// 激活纹理单元
	gl.ActiveTexture(unit)
	// 绑定纹理
	gl.BindTexture(Target, texture)
	gl.BindVertexArray(0) // 结束
	return nil
}

// Update 更新顶点
func (V *Vertex) Update(Program uint32) {
	gl.BindVertexArray(V.VAO) // 绘画
	//? 设置模型位置
	cameraUniform := gl.GetUniformLocation(Program, gl.Str("vP_ModelPos\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &(V.Position[0]))
	//? 设置材质
	// gl.Uniform1i(gl.GetUniformLocation(Program, gl.Str("ourTexture\x00")), 0)
	// ? 设置灯光信息
	// 参数测试
	VfModelColor := mgl32.Vec3{1.0, 0.5, 0.31}
	VfLightColor := mgl32.Vec3{1.0, 1.0, 1.0}
	VflightPos := mgl32.Vec3{2.0, 2.0, 0.0}
	// 设置灯光参数
	UniformobjectColor := gl.GetUniformLocation(Program, gl.Str("fP_ModelColor\x00"))
	UniformlightColor := gl.GetUniformLocation(Program, gl.Str("fP_LightColor\x00"))
	UniformlightPos := gl.GetUniformLocation(Program, gl.Str("fP_LightPos\x00"))
	gl.Uniform3fv(UniformobjectColor, 1, &VfModelColor[0]) // 物体颜色
	gl.Uniform3fv(UniformlightColor, 1, &VfLightColor[0])  // 光源颜色
	gl.Uniform3fv(UniformlightPos, 1, &VflightPos[0])      // 灯光位置

	//? 判断是否为索引
	if V.ifIndex {
		gl.DrawElements(V.DisplayMode, V.indexN, gl.UNSIGNED_INT, gl.PtrOffset(0))
	} else {
		gl.DrawArrays(V.DisplayMode, 0, V.indexN)
	}
	gl.BindVertexArray(0) // 结束
}

// Delete 销毁
func (V *Vertex) Delete() error {
	if V.ifCreate {
		gl.DeleteVertexArrays(0, &(V.VAO))
		gl.DeleteBuffers(0, &(V.Buffer))
		V.VAO = 0
		V.Buffer = 0
		V.ifCreate = false
		if V.ifCreate {
			gl.DeleteBuffers(0, &(V.indexIbo))
		}
	}
	return nil
}
