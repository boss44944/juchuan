let deviceId = localStorage.getItem('juchuan_device_id');

if (!deviceId) {
 deviceId = 'device-' + Math.random().toString(36).substring(2);
 localStorage.setItem('juchuan_device_id', deviceId);
}

function registerDevice(){
 fetch('/api/device/register',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({id:deviceId,name:navigator.userAgent})});
}

function connectWS(){
 let ws=new WebSocket('ws://' + location.host + '/ws?device=' + deviceId);
 ws.onmessage=function(e){
  let msg=JSON.parse(e.data);
  if(msg.type==='file' && confirm('收到文件: '+msg.filename)) location.href=msg.url;
 };
 ws.onclose=function(){setTimeout(connectWS,3000);};
}

function loadHistory(){
 fetch('/api/history?page=1&size=20').then(r=>r.json()).then(list=>{
  let box=document.getElementById('history');
  if(!box)return;
  box.innerHTML=list.map(h=>h.type==='file'?'<div>'+h.filename+' <a href="/download/'+h.id+'">下载</a></div>':'<div>'+h.content+'</div>').join('');
 });
}

function loadDevices(){
 fetch('/api/devices').then(r=>r.json()).then(list=>{
  let box=document.getElementById('devices');
  if(!box)return;
  box.innerHTML=list.map(d=>'<div>'+d.Name+'</div>').join('');
 });
}

function sendText(){
 fetch('/api/text',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({content:document.getElementById('text').value})});
}

function uploadBlob(blob,name){
 let d=new FormData();
 d.append('file',blob,name||'clipboard.png');
 fetch('/upload',{method:'POST',body:d});
}

function uploadFile(){
 let f=document.getElementById('file').files[0];
 if(f) uploadBlob(f,f.name);
}

document.addEventListener('dragover',e=>e.preventDefault());
document.addEventListener('drop',e=>{
 e.preventDefault();
 if(e.dataTransfer.files.length) uploadBlob(e.dataTransfer.files[0],e.dataTransfer.files[0].name);
});

document.addEventListener('paste',e=>{
 for(let item of e.clipboardData.items){
  if(item.type.startsWith('image/')) uploadBlob(item.getAsFile(),'clipboard.png');
 }
});

registerDevice();
connectWS();
loadHistory();
loadDevices();
