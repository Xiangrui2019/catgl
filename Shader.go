package catgl

// 着色器类
//   实现着色器相关操作
// ? 日志
// !  2019-8-3 重构
import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Shader 着色器类
type Shader struct {
	Vertex   string // 顶点着色器
	Geometry string // 几何着色器
	Fragment string // 片面着色器
	Program  uint32 // 着色器
	// 顶点组
	QueueVertex []*Vertex
	// 标记
	ifCreate bool
}

// New 创建着色器
func (S *Shader) New() error {
	if S.Vertex == "" || S.Fragment == "" || S.Geometry == "" {
		return fmt.Errorf("无法创建:\n\r 顶点着色器: %v\n\r 几何着色器: %v\n\r 片面着色器: %v\n\r", S.Vertex, S.Geometry, S.Fragment)
	}
	// 保证释放
	S.Delete()
	S.ifCreate = true
	// 创建着色器
	vertex, err := NewShader(S.Vertex, gl.VERTEX_SHADER)
	if err != nil {
		gl.DeleteShader(vertex)
		return err
	}
	geometry, err := NewShader(S.Geometry, gl.GEOMETRY_SHADER)
	if err != nil {
		gl.DeleteShader(vertex)
		gl.DeleteShader(geometry)
		return err
	}
	fragment, err := NewShader(S.Fragment, gl.FRAGMENT_SHADER)
	if err != nil {
		gl.DeleteShader(vertex)
		gl.DeleteShader(geometry)
		gl.DeleteShader(fragment)
		return err
	}
	// 编译着色器 -> 着色器程序
	Program, err := NewProgram(vertex, geometry, fragment)
	// 销毁着色器代码
	gl.DeleteShader(vertex)
	gl.DeleteShader(geometry)
	gl.DeleteShader(fragment)
	// 处理编译错误
	if err != nil {
		gl.DeleteProgram(S.Program)
		return err
	}
	S.Program = Program
	S.ifCreate = true
	return nil
}

// NewVertex 创建顶点组
func (S *Shader) NewVertex() (V *Vertex) {
	V = &Vertex{
		Position: mgl32.Ident4(),
	}
	S.QueueVertex = append(S.QueueVertex, V)
	return
}

// Delete 销毁着色器
func (S *Shader) Delete() error {
	if S.ifCreate {
		// 删除着色器对象
		gl.DeleteProgram(S.Program)
		// 初始化
		S.ifCreate = false
		S.Program = 0
	}
	return nil
}

// Update 更新着色器
func (S *Shader) Update() {
	if S.ifCreate {
		//? 更新顶点列表
		for _, Vertex := range S.QueueVertex {
			//? 更新顶点
			Vertex.Update(S.Program)
		}
	}
}

// NewShader 创建着色器
func NewShader(source string, shaderType uint32) (uint32, error) {
	// 创建着色器
	shader := gl.CreateShader(shaderType)
	// 获得指针
	csource, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csource, nil)
	// 销毁缓存
	free()
	// 编译
	gl.CompileShader(shader)
	// 获得错误
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("编译着色器失败 %v: %v", source, log)
	}
	return shader, nil
}

// DeleteShader 删除着色器
func DeleteShader(shader uint32) {
	gl.DeleteShader(shader)
}

// NewProgram 编译着色器程序
func NewProgram(vertexShader, geometryShader, fragmentShader uint32) (uint32, error) {
	// 着色器程序
	shaderProgram := gl.CreateProgram()
	// 设置
	if vertexShader != 0 {
		gl.AttachShader(shaderProgram, vertexShader) // 顶点着色器
	}
	if geometryShader != 0 {
		gl.AttachShader(shaderProgram, geometryShader) // 几何着色器
	}
	if fragmentShader != 0 {
		gl.AttachShader(shaderProgram, fragmentShader) // 片段着色器
	}
	// 链接
	gl.LinkProgram(shaderProgram)
	// Go - 错误获取
	var status int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		gl.DeleteShader(shaderProgram)
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("着色器链接失败: %v \n\r", log)
	}
	return shaderProgram, nil
}

// NewTexture 创建材质
// *   file 材质文件名
// *   Target 纹理类型
func NewTexture(file string, Target uint32) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	defer imgFile.Close() // 退出关闭文件
	// 解码图片
	img, err := png.Decode(imgFile)
	if err != nil {
		return 0, err
	}
	// 得到图片通道信息
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	// 转换格式
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(Target, texture)
	// 纹理参数
	gl.TexParameteri(Target, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(Target, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(Target, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(Target, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	// 添加纹理
	gl.TexImage2D(
		Target,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	// 解除纹理
	gl.BindTexture(Target, 0)
	return texture, nil
}
