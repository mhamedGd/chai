#shader vertex
#version 300 es

precision mediump float;

in vec2 coordinates;
in vec4 colors;

out vec4 vertex_FragColor;

uniform mat4 projection_matrix;
uniform mat4 view_matrix;

void main(void) {
	vec4 global_position = vec4(0.0);
	global_position = view_matrix * vec4(coordinates, 0.0, 1.0);
	global_position.z = 0.0;
	global_position.w = 1.0;		
	gl_Position = global_position;
	
	vertex_FragColor = colors;
}

#shader fragment
#version 300 es

precision mediump float;

in vec4 vertex_FragColor;

out vec4 fragColor;
void main(void) {
    fragColor = vertex_FragColor;
}