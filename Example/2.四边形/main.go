package main

import (
	"fmt"

	"gitee.com/LittleRuicat/catgl"
	// "github.com/go-gl/gl/v3.3-core/gl"
	// "github.com/go-gl/glfw/v3.1/glfw"
)

func main() {
	//* 创建窗口
	Gw, _ := catgl.ShowGlNew(900, 600, "创建四边形示例")
	//*　创建四边形
	Triangle(Gw)
	//* 主循环
	catgl.ShowGlLoop()
}

// Triangle 顶点创建四边形
func Triangle(Gw *catgl.ShowGl) {
	//? 创建顶点
	shader := NewShader(Gw)
	vertex := shader.NewVertex()
	//?　设置顶点
	vertex.SetVertex([]float32{
		//* 顶点位置
		0.5, 0.5, 0.0, // 右上角
		0.5, -0.5, 0.0, // 右下角
		-0.5, -0.5, 0.0, // 左下角
		-0.5, 0.5, 0.0, // 左上角
	}, nil, nil)
	//? 设置顶点索引
	vertex.SetIndex([]uint32{
		0, 1, 3, // 第一个三角形
		1, 2, 3, // 第二个三角形
	})
	//? 主渲染
	Gc := (&catgl.Camera{}).New(2, 2, 0).Set(Gw) //? 创建相机
	Gw.AddRender("四边形", func() {
		Gc.Update()
	})
}

// NewShader 默认着色器
func NewShader(Gw *catgl.ShowGl) *catgl.Shader {
	//? 设置当前上下文
	Gw.SetContext()
	s, err := Gw.NewShader(`
	#version 330 core
	//? 默认数据顶点
	layout (location = 0) in vec3 apositions;	//* 位置
	layout (location = 1) in vec3 anormals;		//* 法线
	layout (location = 2) in vec2 auv;			//* uv
	//? 引擎传递参数
	uniform mat4 vP_Projection;    //* 投影矩阵
	uniform mat4 vP_CameraPos;     //* 相机位置
	uniform mat4 vP_ModelPos;      //* 模型位置(Vertex类)
	//? 顶点着色器输出结构
	struct vP {
		vec2 ModelUv;	  //* 模型 Uv
		vec3 ModelNormal; //* 模型 法线
		vec3 FragPos;     //* 摄像机视点(顶点位置)
		vec3 CameraPos;   //* 摄像机位置
	};
	//? 传递到几何着色器
	out vP vP_out; 
	//? 主处理
	void main(){
		gl_Position = vP_Projection * vP_CameraPos * vP_ModelPos * vec4(apositions, 1);
		// 处理传值
		vP_out.ModelUv = auv;
		vP_out.ModelNormal = mat3(transpose(inverse(vP_ModelPos))) * anormals;
		vP_out.FragPos = vec3(vP_ModelPos * vec4(apositions, 1.0));
		vP_out.CameraPos = vec3(vP_CameraPos);
	}
	`, `
	#version 330 core
	layout(triangles) in ;
	layout(triangle_strip, max_vertices = 3) out;
	//? 顶点着色器输出结构
	struct vP {
		vec2 ModelUv;	  //* 模型 Uv
		vec3 ModelNormal; //* 模型 法线
		vec3 FragPos;     //* 摄像机视点(顶点位置)
		vec3 CameraPos;   //* 摄像机位置
	};
	//? 得到顶点着色器传值
	in vP[] vP_out; 
	//? 输出到目标片面着色器
  	out vP gP_out; 
	void main()
	{
		gP_out = vP_out[0];
		gl_Position = gl_in[0].gl_Position;
		EmitVertex();
		gP_out = vP_out[1];
		gl_Position = gl_in[1].gl_Position;
		EmitVertex();
		gP_out = vP_out[2];
		gl_Position = gl_in[2].gl_Position;
		EmitVertex();
		//* 完成绘制
		EndPrimitive();
	}
	`, `
	#version 330 core
	//? 引擎传递参数
	uniform vec3 fP_ModelColor; //* 物体颜色
	uniform vec3 fP_LightColor; //* 光源颜色
	uniform vec3 fP_LightPos;   //* 光源位置
	//? 顶点着色器输出结构
	struct vP {
		vec2 ModelUv;	  //* 模型 Uv
		vec3 ModelNormal; //* 模型 法线
		vec3 FragPos;     //* 摄像机视点(顶点位置)
		vec3 CameraPos;   //* 摄像机位置
	};
	//? 得到顶点着色器传值
	in vP gP_out; 
	//? 片面着色器输出
	out vec4 fP_Color;
	void main() {
		fP_Color =vec4(fP_ModelColor,1);
	}
	`)
	fmt.Println(err)
	return s
}
