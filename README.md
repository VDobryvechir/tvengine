# tvengine
Backend service for TV info 

1. Computer and Groups 
Computer are united into groups. There is a root group, which includes all groups and computers.
Groups can consist of computers and other groups
Each computer has parameters
  id key parameter
  title *
  url * 
Each group has parameters
  id key parameter
  title *
  computers []
  subgroups []
2. Screens
Each screen has parameters as follows
id key parameter
Title * 
text:
1 message 
  message
  color *
  size *
2 message
  message
  color *
  size *
Image
  id
  title
  url 
  height in %  *
Background
  color *
  image url
Video
  id
  title
  url *
Duration is seconds *
allImage: url or  internal parameter

3. Presentation
  id key parameter
  Title *
  Screens []
  Groups []


1. GROUP API
GET /api/v1/group/{parent id}
  id=0 the ALL group
  retrieves information about a group
  {
      id: 0,
      title: "All"
      groups: [{id, title}],
      computers: [{id,title,url}],
      allGroups:[{id,title}] 
  }
POST /api/v1/group/{parent id}
  creates a group
  {
      title: "Work"
      groupIds: [id],
      computers: [{id,title,url}] 
  }
PUT /api/v1/group/{group id}
  updates a group
  {
      title: "Work"
      groupIds: [id],
      computers: [{id,title,url}] 
  }
DELETE /api/v1/group/{group id}
  deletes a group


2. PICTURE API
GET /picture/name
     load picture
GET /api/v1/picture
    [{id,title,url}] get full list of pictures
POST /api/v1/picture
     {title,file} upload a picture
DELETE /api/v1/picture/{picture id}
     delete a picture


3. VIDEO API
GET /video/name
     load video
GET /api/v1/video
    [{id,title,url}] get full list of videos
POST /api/v1/video
     {title,file} upload a video
DELETE /api/v1/video/{video id}
     delete a video

4. SCREEN API
GET /api/v1/screen
[{
   id: 1,
   title: "",
   text: [{message,color,size}],
   image: {id,title,url,height},
   backgroundColor: #675645,
   backgroundImage: url,
   video: {id, title, url},
   duration: 10,
   allImage: url   
}]
POST /api/v1/screen
creates a new screen
{
   title: "",
   text: [{message,color,size}],
   image: {id,title,url,height},
   backgroundColor: #675645,
   backgroundImage: url,
   video: {id, title, url},
   duration: 10,
   allImage: url    
}                                            
PUT /api/v1/screen
updates an existing screen
{
   id: 6,
   title: "",
   text: [{message,color,size}],
   image: {id,title,url,height},
   backgroundColor: #675645,
   backgroundImage: url,
   video: {id, title, url},
   duration: 10,
   allImage: url   
}
DELETE /api/v1/screen/{id}
   removes screen and all its use
5. PRESENTATION API
GET /api/v1/presentation
  get a list of all presentations
[{
    id:0
    title: "",
    groups: [{id,title}],
    screens: [{id,title,allImage,video,duration}], 
}
]
POST /api/v1/presentation
  create a new presentation
[{
    title: "",
    groupIds: [id],
    screenIds: [id], 
}
]
PUT /api/v1/presentation
  update a presentation
[{
    id:0
    title: "",
    groupIds: [id],
    screenIds: [id], 
}
]
DELETE /api/v1/presentation/id1,id2,id3
   delete a presentation
POST /api/v1/control
   activate presentations
{
  presentations: [1,..]
}  






