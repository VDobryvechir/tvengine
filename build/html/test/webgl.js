let cubeRotation = 0.0;
let deltaTime = 0;
// will set to true when video can be copied to texture
let copyVideo = false;


//
// start here
//
function presentPictureVideosAndDuration(files, durations, canvas) {
  // Initialize the GL context
  const gl = canvas.getContext("webgl");

  // Only continue if WebGL is available and working
  if (gl === null) {
    alert(
      "Unable to initialize WebGL. Your browser or machine may not support it."
    );
    return;
  }

  // Set clear color to black, fully opaque
  gl.clearColor(0.0, 0.0, 0.0, 0.0);
  // Clear the color buffer with specified clear color
  gl.clear(gl.COLOR_BUFFER_BIT);

  // Vertex shader program

  const vsSource = `
  attribute vec4 aVertexPosition;
  attribute vec2 aTextureCoord;

  uniform mat4 uModelViewMatrix;
  uniform mat4 uProjectionMatrix;

  varying highp vec2 vTextureCoord;

  void main(void) {
    gl_Position = uProjectionMatrix * uModelViewMatrix * aVertexPosition;
    vTextureCoord = aTextureCoord;

  }
`;

  // Fragment shader program

  const fsSource = `
  varying highp vec2 vTextureCoord;

  uniform sampler2D uSampler;

  void main(void) {
    highp vec4 texelColor = texture2D(uSampler, vTextureCoord);

    gl_FragColor = texelColor;
  }
`;

  // Initialize a shader program; this is where all the lighting
  // for the vertices and so forth is established.
  const shaderProgram = initShaderProgram(gl, vsSource, fsSource);

  // Collect all the info needed to use the shader program.
  // Look up which attributes our shader program is using
  // for aVertexPosition, aVertexColor and also
  // look up uniform locations.
  const programInfo = {
    program: shaderProgram,
    attribLocations: {
      vertexPosition: gl.getAttribLocation(shaderProgram, "aVertexPosition"),
      textureCoord: gl.getAttribLocation(shaderProgram, "aTextureCoord"),
    },
    uniformLocations: {
      projectionMatrix: gl.getUniformLocation(
        shaderProgram,
        "uProjectionMatrix"
      ),
      modelViewMatrix: gl.getUniformLocation(shaderProgram, "uModelViewMatrix"),
      uSampler: gl.getUniformLocation(shaderProgram, "uSampler"),
    },
  };

  // Here's where we call the routine that builds all the
  // objects we'll be drawing.
  const buffers = initBuffers(gl);

  //const textures = [];
  //const textur = initTexture(gl);
  //textures[0] = get
  const textures = loadTextures(gl, files, durations);

  // Flip image pixels into the bottom-to-top order that WebGL expects.
  gl.pixelStorei(gl.UNPACK_FLIP_Y_WEBGL, true);

  let then = 0;

  // Draw the scene repeatedly
  function render(now) {
    now *= 0.0004; // convert to seconds
    deltaTime = now - then;
    then = now;
    textures[0].prepareTexture(gl,0);
    let txt = textures[0].getTexture(gl);
    // if (video.copyVideo) {
    //  updateTexture(gl, texture, video.video);
    // }

    drawScene(gl, programInfo, buffers, txt, cubeRotation);
    cubeRotation += deltaTime;

    requestAnimationFrame(render);
  }

  requestAnimationFrame(render);
}

//
// Initialize a shader program, so WebGL knows how to draw our data
//
function initShaderProgram(gl, vsSource, fsSource) {
  const vertexShader = loadShader(gl, gl.VERTEX_SHADER, vsSource);
  const fragmentShader = loadShader(gl, gl.FRAGMENT_SHADER, fsSource);

  // Create the shader program

  const shaderProgram = gl.createProgram();
  gl.attachShader(shaderProgram, vertexShader);
  gl.attachShader(shaderProgram, fragmentShader);
  gl.linkProgram(shaderProgram);

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
function initTexture(gl) {
  const texture = gl.createTexture();
  gl.bindTexture(gl.TEXTURE_2D, texture);

  // Because video has to be download over the internet
  // they might take a moment until it's ready so
  // put a single pixel in the texture so we can
  // use it immediately.
  const level = 0;
  const internalFormat = gl.RGBA;
  const width = 1;
  const height = 1;
  const border = 0;
  const srcFormat = gl.RGBA;
  const srcType = gl.UNSIGNED_BYTE;
  const pixel = new Uint8Array([0, 0, 255, 255]); // opaque blue
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

  // Turn off mips and set wrapping to clamp to edge so it
  // will work regardless of the dimensions of the video.
  gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
  gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
  gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR);

  return texture;
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

