function initBuffers(gl) {
  const positionBuffer = initPositionBuffer(gl);

  const textureCoordBuffer = [initTextureBuffer(gl,0),initTextureBuffer(gl,1),initTextureBuffer(gl,2),initTextureBuffer(gl,3)];

  const indexBuffer = initIndexBuffer(gl);

  return {
    position: positionBuffer,
    textureCoord: textureCoordBuffer,
    indices: indexBuffer,
  };
}

function initPositionBuffer(gl) {
  // Create a buffer for the square's positions.
  const positionBuffer = gl.createBuffer();

  // Select the positionBuffer as the one to apply buffer
  // operations to from here out.
  gl.bindBuffer(gl.ARRAY_BUFFER, positionBuffer);

  const positions = [
    // Front face
    -1.0, -1.0, 1.0, 1.0, -1.0, 1.0, 1.0, 1.0, 1.0, -1.0, 1.0, 1.0,

    // Back face
    -1.0, -1.0, -1.0, -1.0, 1.0, -1.0, 1.0, 1.0, -1.0, 1.0, -1.0, -1.0,

    // Top face
    -1.0, 1.0, -1.0, -1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, -1.0,

    // Bottom face
    -1.0, -1.0, -1.0, 1.0, -1.0, -1.0, 1.0, -1.0, 1.0, -1.0, -1.0, 1.0,

    // Right face
    1.0, -1.0, -1.0, 1.0, 1.0, -1.0, 1.0, 1.0, 1.0, 1.0, -1.0, 1.0,

    // Left face
    -1.0, -1.0, -1.0, -1.0, -1.0, 1.0, -1.0, 1.0, 1.0, -1.0, 1.0, -1.0,
  ];

  // Now pass the list of positions into WebGL to build the
  // shape. We do this by creating a Float32Array from the
  // JavaScript array, then use it to fill the current buffer.
  gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(positions), gl.STATIC_DRAW);

  return positionBuffer;
}

function initColorBuffer(gl) {
  const faceColors = [
    [1.0, 1.0, 1.0, 1.0], // Front face: white
    [1.0, 0.0, 0.0, 1.0], // Back face: red
    [0.0, 1.0, 0.0, 1.0], // Top face: green
    [0.0, 0.0, 1.0, 1.0], // Bottom face: blue
    [1.0, 1.0, 0.0, 1.0], // Right face: yellow
    [1.0, 0.0, 1.0, 1.0], // Left face: purple
  ];

  // Convert the array of colors into a table for all the vertices.

  var colors = [];

  for (var j = 0; j < faceColors.length; ++j) {
    const c = [0.0,0.0,0.0,0.0] || faceColors[j];
    // Repeat each color four times for the four vertices of the face
    colors = colors.concat(c, c, c, c);
  }

  const colorBuffer = gl.createBuffer();
  gl.bindBuffer(gl.ARRAY_BUFFER, colorBuffer);
  gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(colors), gl.STATIC_DRAW);

  return colorBuffer;
}

function initIndexBuffer(gl) {
  const indexBuffer = gl.createBuffer();
  gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer);

  // This array defines each face as two triangles, using the
  // indices into the vertex array to specify each triangle's
  // position.

  const indices = [
    0,
    1,
    2,
    0,
    2,
    3, // front
    4,
    5,
    6,
    4,
    6,
    7, // back
    8,
    9,
    10,
    8,
    10,
    11, // top
    12,
    13,
    14,
    12,
    14,
    15, // bottom
    16,
    17,
    18,
    16,
    18,
    19, // right
    20,
    21,
    22,
    20,
    22,
    23, // left
  ];

  // Now send the element array to GL

  gl.bufferData(
    gl.ELEMENT_ARRAY_BUFFER,
    new Uint16Array(indices),
    gl.STATIC_DRAW
  );

  return indexBuffer;
}

function initTextureBuffer(gl, nr) {
  const textureCoordBuffer = gl.createBuffer();
  gl.bindBuffer(gl.ARRAY_BUFFER, textureCoordBuffer);
  if (!nr) {
     nr=0
  }
  let textureCoordinates = [];
  switch (nr) {
   case 0:
  textureCoordinates = [
    // Front (it was one of them standing right)
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Back  (it was one of them standing vertically, opposite to front)
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Top  it is invisible if rotation is about Y
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Bottom it is invisible if rotation is about Y
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Right  it is number 3, standing vertically
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Left   it is number 1, standing right
    0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0,
  ];
     break;    
   case 1:
  textureCoordinates = [
    // Front (it was one of them standing right)
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Back  (it was one of them standing vertically, opposite to front)
    1.0, 0.0, 1.0, 1.0, 0.0, 1.0, 0.0, 0.0,
    // Top  it is invisible if rotation is about Y
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Bottom it is invisible if rotation is about Y
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Right  it is number 3, standing vertically
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Left   it is number 1, standing right
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
  ];
     break;    
   case 2:
  textureCoordinates = [
    // Front (it was one of them standing right)
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Back  (it was one of them standing vertically, opposite to front)
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Top  it is invisible if rotation is about Y
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Bottom it is invisible if rotation is about Y
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Right  it is number 3, standing vertically
    1.0, 0.0, 1.0, 1.0, 0.0, 1.0, 0.0, 0.0,
    // Left   it is number 1, standing right
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
  ];
     break;                                
   case 3:
  textureCoordinates = [
    // Front (it was one of them standing right)
    0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0,
    // Back
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Top
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Bottom
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Right
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
    // Left
    0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
  ];
     break;    
   case 4:
  textureCoordinates = [
    // Front
    0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0,
    // Back
    0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0,
    // Top
    0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0,
    // Bottom
    0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0,
    // Right
    0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0,
    // Left
    0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0,
  ];
     break;    
  }
  gl.bufferData(
    gl.ARRAY_BUFFER,
    new Float32Array(textureCoordinates),
    gl.STATIC_DRAW
  );

  return textureCoordBuffer;
}
