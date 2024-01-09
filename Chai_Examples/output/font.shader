#shader vertex
#version 300 es

precision mediump float;

in vec2 coordinates;
in vec4 colors;
in vec2 uv;

out vec4 vertex_FragColor;
out vec2 vertex_UV;

uniform mat4 projection_matrix;
uniform mat4 view_matrix;

void main(void) {
	vec4 global_position = vec4(0.0);
	global_position = view_matrix * vec4(coordinates, 0.0, 1.0);
	global_position.z = 0.0;
	global_position.w = 1.0;		
	gl_Position = global_position;
	
	
	vertex_FragColor = colors;
	vertex_UV = uv;
}

#shader fragment
#version 300 es

precision mediump float;

in vec4 vertex_FragColor;
in vec2 vertex_UV;

uniform sampler2D genericSampler;

out vec4 fragColor;

void main(void) {
	vec4 thisColor = texture(genericSampler, vertex_UV);
	if(thisColor.a > 0.0){
		thisColor = vec4(1.0);
	}else {
		thisColor = vec4(0.0);
	}
	fragColor = thisColor * vertex_FragColor;
}