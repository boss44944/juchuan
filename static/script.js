function sendText(){
 fetch('/api/text',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({content:document.getElementById('text').value})})
}
function uploadFile(){
 let f=document.getElementById('file').files[0];
 let d=new FormData();d.append('file',f);
 fetch('/api/upload',{method:'POST',body:d});
}
