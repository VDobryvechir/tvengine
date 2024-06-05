
function drawScene(gl, programInfo, buffers, textures, cubeRotation) {
  gl.clearColor(0.0, 0.0, 0.0, 1.0); // Clear to black, fully opaque
  gl.clearDepth(1.0); // Clear everything
  gl.enable(gl.DEPTH_TEST); // Enable depth testing
  gl.depthFunc(gl.LEQUAL); // Near things obscure far things

  // Clear the canvas before we start drawing on it.

  gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);

  // Create a perspective matrix, a special matrix that is
  // used to simulate the distortion of perspective in a camera.
  // Our field of view is 45 degrees, with a width/height
  // ratio that matches the display size of the canvas
  // and we only want to see objects between 0.1 units
  // and 100 units away from the camera.

  const fieldOfView = (23 * Math.PI) / 180; // in radians
  const aspect = gl.canvas.clientWidth / gl.canvas.clientHeight;
  const zNear = 0.1;
  const zFar = 100.0;
  const projectionMatrix = mat4.create();

  // note: glmatrix.js always has the first argument
  // as the destination to receive the result.
  mat4.perspective(projectionMatrix, fieldOfView, aspect, zNear, zFar);

  // Set the drawing position to the "identity" point, which is
  // the center of the scene.
  const modelViewMatrix = mat4.create();

  // Now move the drawing position a bit to where we want to
  // start drawing the square.
  mat4.translate(
    modelViewMatrix, // destination matrix
    modelViewMatrix, // matrix to translate
    [-0.0, 0.0, -6.0]
  ); // amount to translate
/*************
  mat4.rotate(
    modelViewMatrix, // destination matrix
    modelViewMatrix, // matrix to rotate
    cubeRotation, // amount to rotate in radians
    [0, 0, 1]
  ); // axis to rotate around (Z)
***************/
  mat4.rotate(
    modelViewMatrix, // destination matrix
    modelViewMatrix, // matrix to rotate
    cubeRotation * Math.PI * 0.1, // amount to rotate in radians
    [0, 1, 0]
  ); // axis to rotate around (Y)
/***************************
  mat4.rotate(
    modelViewMatrix, // destination matrix
    modelViewMatrix, // matrix to rotate
    cubeRotation * 0.3, // amount to rotate in radians
    [1, 0, 0]
  ); // axis to rotate around (X)
****************************/
  // Tell WebGL how to pull out the positions from the position
  // buffer into the vertexPosition attribute.
  setPositionAttribute(gl, buffers, programInfo);

  setTextureAttribute(gl, buffers, programInfo);

  // Tell WebGL which indices to use to index the vertices
  gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffers.indices);

  // Tell WebGL to use our program when drawing
  gl.useProgram(programInfo.program);

  // Set the shader uniforms
  gl.uniformMatrix4fv(
    programInfo.uniformLocations.projectionMatrix,
    false,
    projectionMatrix
  );
  gl.uniformMatrix4fv(
    programInfo.uniformLocations.modelViewMatrix,
    false,
    modelViewMatrix
  );

  for(let i=0;i<4;i++) {
  	// Tell WebGL we want to affect texture unit 0
  	gl.activeTexture(gl.TEXTURE0 + i);

  	// Bind the texture to texture unit 0
  	gl.bindTexture(gl.TEXTURE_2D, textures[logicTexture[i]].getTexture(gl));

  	// Tell the shader we bound the texture to texture unit 0
  	gl.uniform1i(programInfo.uniformLocations.uSampler[i], i);
  }
  const cubPart = (Math.floor(cubeRotation / 5)+2) % 4;

  if (logicBasis==3 || logicBasis>4) {
     logicTexture[cubPart]=(logicTexture[(cubPart+3)%4]+1)  % logicBasis;
  }
  {
    const vertexCount = 36;
    const type = gl.UNSIGNED_SHORT;
    const offset = 0;
    gl.drawElements(gl.TRIANGLES, vertexCount, type, offset);
  }
  for(let i=0;i<textures.length;i++) {
        const fn = textures[i]; 
        if (fn.prepareTexture) {
            fn.prepareTexture(gl,i);
        }
  }
}

// Tell WebGL how to pull out the positions from the position
// buffer into the vertexPosition attribute.
function setPositionAttribute(gl, buffers, programInfo) {
  const numComponents = 3;
  const type = gl.FLOAT; // the data in the buffer is 32bit floats
  const normalize = false; // don't normalize
  const stride = 0; // how many bytes to get from one set of values to the next
  // 0 = use type and numComponents above
  const offset = 0; // how many bytes inside the buffer to start from
  gl.bindBuffer(gl.ARRAY_BUFFER, buffers.position);
  gl.vertexAttribPointer(
    programInfo.attribLocations.vertexPosition,
    numComponents,
    type,
    normalize,
    stride,
    offset
  );
  gl.enableVertexAttribArray(programInfo.attribLocations.vertexPosition);
}

// Tell WebGL how to pull out the colors from the color buffer
// into the vertexColor attribute.
function setColorAttribute(gl, buffers, programInfo) {
  const numComponents = 4;
  const type = gl.FLOAT;
  const normalize = false;
  const stride = 0;
  const offset = 0;
  gl.bindBuffer(gl.ARRAY_BUFFER, buffers.color);
  gl.vertexAttribPointer(
    programInfo.attribLocations.vertexColor,
    numComponents,
    type,
    normalize,
    stride,
    offset
  );
  gl.enableVertexAttribArray(programInfo.attribLocations.vertexColor);
}

// tell webgl how to pull out the texture coordinates from buffer
function setTextureAttribute(gl, buffers, programInfo) {
 for(let i=0;i<4;i++) {
  const num = 2; // every coordinate composed of 2 values
  const type = gl.FLOAT; // the data in the buffer is 32-bit float
  const normalize = false; // don't normalize
  const stride = 0; // how many bytes to get from one set to the next
  const offset = 0; // how many bytes inside the buffer to start from
  gl.bindBuffer(gl.ARRAY_BUFFER, buffers.textureCoord[i]);
  gl.vertexAttribPointer(
    programInfo.attribLocations.textureCoord[i],
    num,
    type,
    normalize,
    stride,
    offset
  );
  gl.enableVertexAttribArray(programInfo.attribLocations.textureCoord[i]);
 }
}


let cubeRotation = 0.0;
let deltaTime = 0;
const logicTexture=[];
let logicBasis = 0;
let drawFinished = null;
let drawResolve = null;

async function stopPictureVideo(){
    if (logicTexture.length===0) {
       return Promise.resolve();
    }
    if (drawFinished) {
       return drawFinished;
    }
    drawFinished=new Promise((resolve,reject)=>{
       drawResolve = resolve;
    });
    return drawFinished; 
}


function presentPictureVideosAndDuration(files, durations, canvas) {
  drawFinished = null;
  cubeRotation = 0.0;
  deltaTime = 0;
  logicBasis = files.length;
  // Initialize the GL context
  const gl = canvas.getContext("webgl");

  // Only continue if WebGL is available and working
  if (gl === null) {
    alert(
      "Unable to initialize WebGL. Your browser or machine may not support it."
    );
    return;
  }
  for(let i=0;i<4;i++) {
      logicTexture[i] = i % logicBasis;
  }
  if (logicBasis===3) {
      logicTexture[3]=2;
  }
  // Set clear color to black, fully transparent
  gl.clearColor(0.0, 0.0, 0.0, 0.0);
  // Clear the color buffer with specified clear color
  gl.clear(gl.COLOR_BUFFER_BIT);

  // Vertex shader program

  const vsSource = `
  attribute vec4 aVertexPosition;
  attribute vec2 aTextureCoord0;
  attribute vec2 aTextureCoord1;  
  attribute vec2 aTextureCoord2;
  attribute vec2 aTextureCoord3;

  uniform mat4 uModelViewMatrix;
  uniform mat4 uProjectionMatrix;

  varying highp vec2 vTextureCoord0;
  varying highp vec2 vTextureCoord1;
  varying highp vec2 vTextureCoord2;
  varying highp vec2 vTextureCoord3;
  
  void main(void) {
    gl_Position = uProjectionMatrix * uModelViewMatrix * aVertexPosition;
    vTextureCoord0 = aTextureCoord0;
    vTextureCoord1 = aTextureCoord1;
    vTextureCoord2 = aTextureCoord2;
    vTextureCoord3 = aTextureCoord3;
  }
`;

  // Fragment shader program

  const fsSource = `
  varying highp vec2 vTextureCoord0;
  varying highp vec2 vTextureCoord1;
  varying highp vec2 vTextureCoord2;
  varying highp vec2 vTextureCoord3;

  uniform sampler2D uSampler0;
  uniform sampler2D uSampler1;
  uniform sampler2D uSampler2;
  uniform sampler2D uSampler3;

  void main(void) {
    highp vec4 color0 = texture2D(uSampler0, vTextureCoord0);
    highp vec4 color1 = texture2D(uSampler1, vTextureCoord1);
    highp vec4 color2 = texture2D(uSampler2, vTextureCoord2);
    highp vec4 color3 = texture2D(uSampler3, vTextureCoord3);
    gl_FragColor =  color0 + color1 + color2 + color3;
  }
`;

  const releaseBuf = [];
   //  + color1 + color2 + color3
  // Initialize a shader program; this is where all the lighting
  // for the vertices and so forth is established.
  const shaderProgram = initShaderProgram(gl, vsSource, fsSource, releaseBuf);

  // Collect all the info needed to use the shader program.
  // Look up which attributes our shader program is using
  // for aVertexPosition, aVertexColor and also
  // look up uniform locations.
  const programInfo = {
    program: shaderProgram,
    attribLocations: {
      vertexPosition: gl.getAttribLocation(shaderProgram, "aVertexPosition"),
      textureCoord: [gl.getAttribLocation(shaderProgram, "aTextureCoord0"),gl.getAttribLocation(shaderProgram, "aTextureCoord1"),gl.getAttribLocation(shaderProgram, "aTextureCoord2"), gl.getAttribLocation(shaderProgram, "aTextureCoord3")],
    },
    uniformLocations: {
      projectionMatrix: gl.getUniformLocation(
        shaderProgram,
        "uProjectionMatrix"
      ),
      modelViewMatrix: gl.getUniformLocation(shaderProgram, "uModelViewMatrix"),
      uSampler: [gl.getUniformLocation(shaderProgram, "uSampler0"),gl.getUniformLocation(shaderProgram, "uSampler1"),gl.getUniformLocation(shaderProgram, "uSampler2"),gl.getUniformLocation(shaderProgram, "uSampler3")],
    },
  };

  // Here's where we call the routine that builds all the
  // objects we'll be drawing.
  const buffers = initBuffers(gl);

  // Load textures
  const textures = loadTextures(gl, files, durations);
  // Flip image pixels into the bottom-to-top order that WebGL expects.
  gl.pixelStorei(gl.UNPACK_FLIP_Y_WEBGL, true);

  let then = 0;
  function releaseBuffers() {
     logicTexture.length = 0;
     let n = textures.length;
     for(let i=0;i<n;i++) {
         gl.deleteTexture(textures[i].texture);
     }
     textures.length=0;
     n = releaseBuf.length;
     for(let i=0;i<n;i++) {
         releaseBuf[i]();
     }
     releaseBuf.length=0;
  }
  // Draw the scene repeatedly
  function render(now) {
    now *= 0.002; // convert to seconds
    deltaTime = now - then;
    then = now;
    drawScene(gl, programInfo, buffers, textures, cubeRotation);
    cubeRotation += deltaTime;
    if (drawResolve) {
       releaseBuffers();
       drawResolve();
       drawResolve = null; 
    } else {
       requestAnimationFrame(render);
    }
  }
  requestAnimationFrame(render);
}

//
// Initialize a shader program, so WebGL knows how to draw our data
//
function initShaderProgram(gl, vsSource, fsSource, releaseBuf) {
  const vertexShader = loadShader(gl, gl.VERTEX_SHADER, vsSource);
  const fragmentShader = loadShader(gl, gl.FRAGMENT_SHADER, fsSource);

  // Create the shader program

  const shaderProgram = gl.createProgram();
  gl.attachShader(shaderProgram, vertexShader);
  gl.attachShader(shaderProgram, fragmentShader);
  gl.linkProgram(shaderProgram);

  releaseBuf.push(gl.deleteProgram.bind(gl, shaderProgram));
  releaseBuf.push(gl.deleteShader.bind(gl, fragmentShader));
  releaseBuf.push(gl.deleteShader.bind(gl, vertexShader));
  // If creating the shader program failed, alert

  if (!gl.getProgramParameter(shaderProgram, gl.LINK_STATUS)) {
    alert(
      `Unable to initialize the shader program: ${gl.getProgramInfoLog(
        shaderProgram
      )}`
    );
    return null;
  }

  return shaderProgram;
}

//
// creates a shader of the given type, uploads the source and
// compiles it.
//
function loadShader(gl, type, source) {
  const shader = gl.createShader(type);

  // Send the source to the shader object

  gl.shaderSource(shader, source);

  // Compile the shader program

  gl.compileShader(shader);

  // See if it compiled successfully

  if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
    alert(
      `An error occurred compiling the shaders: ${gl.getShaderInfoLog(shader)}`
    );
    gl.deleteShader(shader);
    return null;
  }

  return shader;
}


function loadTextures(gl, urls, durations) {
  const textures = [];
  const n = urls.length;
  for(let i=0;i<n;i++) {
      textures[i] = loadTexture(gl, urls[i], durations[i] || 0);
  }
  return textures;
}
//
// Initialize a texture and load an image.
// When the image finished loading copy it into the texture.
//
function loadTexture(gl, url, duration) {
  const texture = gl.createTexture();
  gl.bindTexture(gl.TEXTURE_2D, texture);

  // Because images have to be downloaded over the internet
  // they might take a moment until they are ready.
  // Until then put a single pixel in the texture so we can
  // use it immediately. When the image has finished downloading
  // we'll update the texture with the contents of the image.
  const level = 0;
  const internalFormat = gl.RGBA;
  const width = 1;
  const height = 1;
  const border = 0;
  const srcFormat = gl.RGBA;
  const srcType = gl.UNSIGNED_BYTE;
  const pixel = new Uint8Array([0, 0, 0, 0]); // opaque blue
  gl.texImage2D(
    gl.TEXTURE_2D,
    level,
    internalFormat,
    width,
    height,
    border,
    srcFormat,
    srcType,
    pixel
  );
  const isVideo = (url) => {
     const pos = url.lastIndexOf(".");
     if (pos<=0) {
         return false;
     }
     const ext = url.substring(pos+1).toLowerCase();
     return ext==="mp4" || ext==="webm" || ext==="ogg";
  };
  const isPowerOf2 = (value)=>{
     return (value & (value - 1)) === 0;
  };
  if (isVideo(url)) {
  // Turn off mips and set wrapping to clamp to edge so it
  // will work regardless of the dimensions of the video.
  	gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
  	gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
  	gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR);

        return getVideoTexture(url, texture);
  }
  const durationInfo = (nmb) => {
     if (!(nmb>0)) {
         return "00:00";
     }
     const sec = nmb % 60;
     const minutes = Math.round((nmb-sec)/60);
     return (minutes<10 ? "0" : "") + minutes + ":" + (sec<10? "0" : "")+sec;
  };
  const dstImage = new Image();
  const remakeImage = (img, fn) => {
      const cnv = document.createElement("canvas");
      const headerHeight = img.width > 300 && img.height>100 ? 70: 0;
      cnv.width = img.width + 2;
      cnv.height = img.height + 2 + headerHeight;
      const ctx = cnv.getContext("2d");
      ctx.fillStyle = "#000000";
      ctx.fillRect(0, 0, cnv.width, cnv.height);
      ctx.drawImage(img, 1, 1 + headerHeight);
      if (headerHeight>0) {
         ctx.fillStyle = "#00bfff";
         ctx.fillRect(1, 1, cnv.width-2, headerHeight);
         ctx.font = "64px serif";
         ctx.fillStyle = "#ffffff";
         let textX =  cnv.width / 2 - 50;
         ctx.fillText(durationInfo(duration), textX, 60);
      }
      const base64URI = cnv.toDataURL();
      dstImage.onload = fn; 
      dstImage.src = base64URI;
  }; 
  const image = new Image();
  image.onload = () => {
    delete image.onload;
    remakeImage(image,()=>{
      	delete dstImage.onload;
    	gl.bindTexture(gl.TEXTURE_2D, texture);
    	gl.texImage2D(
      		gl.TEXTURE_2D,
      		level,
      		internalFormat,
      		srcFormat,
      		srcType,
      		dstImage
    	);

    // WebGL1 has different requirements for power of 2 images
    // vs non power of 2 images so check if the image is a
    // power of 2 in both dimensions.
    	if (isPowerOf2(dstImage.width) && isPowerOf2(dstImage.height)) {
      // Yes, it's a power of 2. Generate mips.
      		gl.generateMipmap(gl.TEXTURE_2D);
    	} else {
      // No, it's not a power of 2. Turn off mips and set
      // wrapping to clamp to edge
      		gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
      		gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
      		gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR);
    	}
    });
  };
  image.src = url;
  const imageBlock = {
      texture: texture,
      getTexture: (gl) => {
           return imageBlock.texture;   
      },
  };
  return imageBlock;
}

let dvno=0;
function showImage(img) {
        const canva = document.getElementById("workcanvas");
        if (!canva) {
             return;
        } 
        const ctx = canva.getContext("2d");
        ctx.drawImage(img, 0, 0);
}

function preupdateVideo(videoBlock) {
        let canva = videoBlock.canvas;
        if (!canva) {
           canva = document.createElement("canvas");
           // TODO may be use videoBlock.video.videoHeight, videoBlock.video.videoWidth to optimize view
           canva.width = 320;
           canva.height =  240;
           videoBlock.canvas = canva;
           document.body.appendChild(canva);
        }
        const ctx = canva.getContext("2d");
        ctx.fillStyle = "#000000";
        ctx.fillRect(0, 0, canva.width, canva.height);  
        ctx.drawImage(videoBlock.video, 1, 1,318,238);
        const base64URI = canva.toDataURL();
        dstImage = new Image();    
        dstImage.onload = ()=>{
            videoBlock.image = dstImage; 
            dstImage.onload = null; 
        }; 
        dstImage.src = base64URI;
}
function getVideoTexture(url, texture) {
  const updateTexture = (gl, texture, img, nr) => {
	const level = 0;
  	const internalFormat = gl.RGBA;
  	const srcFormat = gl.RGBA;
  	const srcType = gl.UNSIGNED_BYTE;
        gl.activeTexture(gl.TEXTURE4 + (nr % 4));
  	gl.bindTexture(gl.TEXTURE_2D, texture);
  	gl.texImage2D(
    		gl.TEXTURE_2D,
    		level,
    		internalFormat,
    		srcFormat,
    		srcType,
    		img
  	);
        showImage(img);
        console.log("image", nr);
  }
  const video = document.createElement("video");
  const videoBlock = {
     texture: texture,
     video: video,
     copyVideo: false,
     nr: dvno++,
     getTexture: (gl) => {
         return videoBlock.texture;
     },
     prepareTexture: (gl,nr) => {
      	if (videoBlock.copyVideo) {
                preupdateVideo(videoBlock);
                if (videoBlock.image) {
      		    updateTexture(gl, videoBlock.texture, videoBlock.image, videoBlock.nr);
                }
    	}
     }
  } 

  let playing = false;
  let timeupdate = false;

  video.playsInline = true;
  video.muted = true;
  video.loop = true;

  // Waiting for these 2 events ensures
  // there is data in the video

  const checkReady = ()=>{
    if (playing && timeupdate) {
      videoBlock.copyVideo = true;
      console.log("started video", videoBlock);
    }
  }

  const playingListener = () => {
      playing = true;
      checkReady();
      video.removeEventListener("playing", playingListener, true);	
  };

  video.addEventListener("playing", playingListener, true);

  const timeUpdateListener = () => {
      timeupdate = true;
      checkReady();
      video.removeEventListener("timeupdate", timeUpdateListener, true);	
  };
  video.addEventListener("timeupdate", timeUpdateListener, true);

  video.src = url;
  video.play();

  return videoBlock;
}

