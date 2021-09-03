#version 410

in vec2 position;

out vec2 textureCoords;

void main() {
	gl_Position = vec4(position, 0.0, 1.0);
	textureCoords = vec2((position.x+1.0)*0.5, 1 - (position.y+1.0)*0.5);
}

