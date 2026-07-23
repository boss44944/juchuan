let deviceId = localStorage.getItem('juchuan_device_id');

if (!deviceId) {
 deviceId = 'device-' + Math.random().toString(36).substring(2);
 localStorage.setItem('juchuan_device_id', deviceId);
}

function registerDevice(){
 fetch('/api/device/register',{
  method:'POST',
  headers:{'Content-Type':'application/json'},
  body:JSON.stringify({id:deviceId,name:navigator.userAgent})
 });
}

function connectWS(){
 let ws=new WebSocket('ws://' + location.host + '/ws?device=' + deviceId);
 ws.onmessage=function(e){
  let msg=JSON.parse(e.data);
  if(msg.type==='file'){
   if(confirm('收到文件: '+msg.filename)){
    location.href=msg.url;
   }
  }
 };
 ws.onclose=function(){setTimeout(connectWS,3000);};
}

function loadHistory(){
 fetch('/api/history?page=1&size=20')
 .then(r=>r.json())
 .then(list=>{
  document.getElementById('history').innerHTML=list.map(h=>
   h.type==='file'
   ? '<div>'+h.filename+' <a href="/download/'+h.id+'">下载</a></div>'
   : '<div>'+h.content+'</div>'
  ).join('');
 });
}

function loadDevices(){
 fetch('/api/devices')
 .then(r=>r.json())
 .then(list=>{
  document.getElementById('devices').innerHTML=list.map(d=>
   '<div>'+d.Name+'</div>'
  ).join('');
 });
}

function sendText(){
 fetch('/api/text',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({content:document.getElementById('text').value})});
}

function uploadFile(){
 let f=document.getElementById('file').files[0];
 let d=new FormData();
 d.append('file',f);
 fetch('/upload',{method:'POST',body:d});
}

registerDevice();
connectWS();
loadHistory();
loadDevices();
