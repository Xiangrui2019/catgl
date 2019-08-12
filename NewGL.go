package catgl

// 主类
//   实现了gl窗口创建
// ! 注:
// *   坐标系是右手坐标系
// ? 日志
// !  2019-8-3 重构
// !  2019-8-6 重写完成多窗口创建
import (
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

// ShowGl Gl显示类型
type ShowGl struct {
	QueueRender map[string]func() // 渲染队列
	QueueShader []*Shader         // 绑定的着色器
	// 大小
	Width       int
	Height      int
	AspectRatio float32 // 屏幕高宽比
	// 内部变量
	window *glfw.Window
}

// SetContext 设置上下文
func (G *ShowGl) SetContext() {
	glfw.DetachCurrentContext()   //? 关闭上下文
	G.window.MakeContextCurrent() //? 设置当前窗口上下文
}

// AddRender 添加渲染
func (G *ShowGl) AddRender(Name string, Render func()) {
	G.QueueRender[Name] = Render
}

// NewShader 创建着色器
func (G *ShowGl) NewShader(
	Vertex string, // 顶点着色器
	Geometry string, // 几何着色器
	Fragment string, // 片面着色器
) (S *Shader, err error) {
	S = &Shader{
		Vertex:   Vertex,
		Fragment: Fragment,
		Geometry: Geometry,
	}
	G.QueueShader = append(G.QueueShader, S)
	err = S.New()
	return
}

// ShowGlList 窗口列表
var ShowGlList map[*glfw.Window]*ShowGl

// init 初始化
func init() {
	//? 绑定到线程
	runtime.LockOSThread()
	//? 窗口列表
	ShowGlList = make(map[*glfw.Window]*ShowGl)
	//? 初始化
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	//? 参数
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
}

// ShowGlNew 创建窗口
func ShowGlNew(Width, Height int, Title string) (*ShowGl, error) {
	//? 创建窗口
	window, err := glfw.CreateWindow(Width, Height, Title, nil, nil)
	if err != nil {
		return nil, err
	}
	//? 上下文生效
	window.MakeContextCurrent()
	//? 初始化 gl
	if err := gl.Init(); err != nil {
		return nil, err
	}
	//? 设置参数
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	//? 分离上下文
	glfw.DetachCurrentContext()
	//? 添加
	Gl := &ShowGl{
		QueueRender: make(map[string]func()),
		Width:       Width,
		Height:      Height,
		AspectRatio: float32(Width / Height),
		window:      window,
	}
	//? 返回
	ShowGlList[window] = Gl
	return Gl, err
}

// ShowGlLoop 创建循环
func ShowGlLoop() {
	//? 主循环
	for len(ShowGlList) > 0 {
		glfw.PollEvents()
		// 循环窗口列表
		for window, Gl := range ShowGlList {
			if !window.ShouldClose() {
				//? 上下文生效
				window.MakeContextCurrent()
				//? 背景颜色
				gl.ClearColor(0.1, 0.3, 0.3, 1.0)
				gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
				//? 渲染队列
				for key := range Gl.QueueRender {
					Gl.QueueRender[key]()
				}
				//? 更新
				window.SwapBuffers()
				//? 分离上下文
				glfw.DetachCurrentContext()
			} else {
				//! 销毁窗口
				window.Destroy()
				//! 删除窗口
				delete(ShowGlList, window)
			}
		}
	}
}

// 主类
//   实现了gl窗口创建
// ! 注:
// *   坐标系是右手坐标系
// ? 日志
// !  2019-8-3 重构

/*
// ShowGlQuantity 窗口个数
var ShowGlQuantity chan uint

// 默认
func init() {
	runtime.LockOSThread()
	ShowGlQuantity = make(chan uint, 0)
	// 初始化
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	// 参数
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
}

// ShowGl Gl显示类型
type ShowGl struct {
	AspectRatio float32
	QueueRender map[string]func() // 渲染队列
	QueueShader []*Shader         // 绑定的着色器
	// 内部变量
	window *glfw.Window
	Width  int
	Height int
}

// New 创建
func (G *ShowGl) New(Width, Height int) (*ShowGl, error) {
	var err error
	G.QueueRender = make(map[string]func())
	G.Width = Width
	G.Height = Height
	// 创建窗口
	G.window, err = glfw.CreateWindow(G.Width, G.Height, "GL-Go", nil, nil)
	if err != nil {
		return nil, err
	}
	glfw.WindowHint(glfw.Visible, glfw.False)
	G.AspectRatio = float32(G.Width / G.Height)
	// 上下文生效
	G.window.MakeContextCurrent()
	// 初始化 gl
	if err := gl.Init(); err != nil {
		return nil, err
	}
	// 设置参数
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	return G, nil
}

// Loop 循环
func (G *ShowGl) Loop() {
	// 主循环
	for !G.window.ShouldClose() {
		//背景颜色
		gl.ClearColor(0.1, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// 渲染队列
		for key := range G.QueueRender {
			G.QueueRender[key]()
		}
		// 交换缓冲区
		G.window.SwapBuffers()
		// 同步
		// glfw.PollEvents()
	}
	// ! 销毁窗口
	G.window.Destroy()
}

// AddRender 添加渲染
func (G *ShowGl) AddRender(Name string, Render func()) {
	G.QueueRender[Name] = Render
}

// !-------------------------------------------------! \\

// NewShader 创建着色器
func (G *ShowGl) NewShader(
	Vertex string, // 顶点着色器
	Fragment string, // 片面着色器
) (S *Shader, err error) {
	S = &Shader{
		Vertex:   Vertex,
		Fragment: Fragment,
	}
	G.QueueShader = append(G.QueueShader, S)
	err = S.New()
	return
}
*/
