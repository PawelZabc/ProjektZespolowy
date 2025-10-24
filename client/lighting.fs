#version 330

in vec3 fragNormal;
in vec3 fragPosition;

out vec4 finalColor;

uniform vec3 lightDir;
uniform vec4 baseColor;
uniform vec4 ambientColor;

void main() {
    float diff = max(dot(normalize(fragNormal), -lightDir), 0.0);
    vec4 diffuse = diff * baseColor;
    vec4 ambient = ambientColor;
    finalColor = diffuse + ambient;
}
