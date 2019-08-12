package catgl

// 摄像机类
//   实现了摄像机
// ? 日志
// !  2019-8-3 重构

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Camera 相机类
type Camera struct {
	// 窗口
	ShowGl *ShowGl
	// 坐标
	Up         mgl32.Vec3
	Eye        mgl32.Vec3
	Center     mgl32.Vec3
	Projection mgl32.Mat4
}

// New 创建相机
// ? 		摄像机位置 x,y,z
func (C *Camera) New(x float32, y float32, z float32) *Camera {
	C.Up = mgl32.Vec3{0, 1, 0}
	C.Eye = mgl32.Vec3{x, y, z}
	C.Center = mgl32.Vec3{0, 0, 0}
	return C
}

// Set 绑定到窗口
func (C *Camera) Set(S *ShowGl) *Camera {
	// 设置变量
	C.ShowGl = S
	// 窗口宽高比
	C.Projection = mgl32.Perspective(mgl32.DegToRad(45.0), C.ShowGl.AspectRatio, 0.1, 10.0)
	return C
}

// Update 更新渲染器相机
func (C *Camera) Update() {
	// 循环设置着色器值
	for _, Shader := range C.ShowGl.QueueShader {
		// 设置矩阵
		gl.UseProgram(Shader.Program)
		// ? 投影矩阵
		projectionUniform := gl.GetUniformLocation(Shader.Program, gl.Str("vP_Projection\x00"))
		gl.UniformMatrix4fv(projectionUniform, 1, false, &C.Projection[0])
		// ? 摄像机位置
		cameraUniform := gl.GetUniformLocation(Shader.Program, gl.Str("vP_CameraPos\x00"))
		look := mgl32.LookAtV(C.Eye, C.Center, C.Up)
		gl.UniformMatrix4fv(cameraUniform, 1, false, &look[0])
		// ? 设置灯光信息
		// 参数测试
		VfModelColor := mgl32.Vec3{1.0, 0.5, 0.31}
		VfLightColor := mgl32.Vec3{1.0, 1.0, 1.0}
		VflightPos := mgl32.Vec3{2.0, 2.0, 0.0}
		// 设置灯光参数
		UniformobjectColor := gl.GetUniformLocation(Shader.Program, gl.Str("fP_ModelColor\x00"))
		UniformlightColor := gl.GetUniformLocation(Shader.Program, gl.Str("fP_LightColor\x00"))
		UniformlightPos := gl.GetUniformLocation(Shader.Program, gl.Str("fP_LightPos\x00"))
		gl.Uniform3fv(UniformobjectColor, 1, &VfModelColor[0]) // 物体颜色
		gl.Uniform3fv(UniformlightColor, 1, &VfLightColor[0])  // 光源颜色
		gl.Uniform3fv(UniformlightPos, 1, &VflightPos[0])      // 灯光位置
		// ? 更新着色器
		Shader.Update()
	}
}
