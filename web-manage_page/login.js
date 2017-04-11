



var mysocket = new WebSocket("ws://10.1.11.151:9003");
document.getElementById("ack").value=  JSON.stringify(mysocket.status)	


mysocket.onopen=function(){
		console.log("ws建立连接");
	}	
mysocket.onclose=function(){
		console.log("ws关闭连接");
	}	
	
mysocket.onmessage = function(data) {
 	  		console.log(new Date(),data);
     var msg = JSON.parse(data.data);
      var data=msg.data;

      switch(msg.cmd){
          case "auth_ok":
			stayHome=0;
			window.location.href="index.html";
          break;
		  
          case "auth_failed":

  			document.getElementById("ack").value='登陆异常';
        break;
		  
           case "auth_name_fault":
 			document.getElementById("ack").value='用户名不存在';
         break;

           case "auth_pwd_fault":
			document.getElementById("ack").value='密码错误';
          break;
		  
			default:
			document.getElementById("ack").value='什么鬼！！';

           case "HB":
			ack_HB()
          break;
		  
          break;

     
      }

}	
	
//comment	
 function login(){
	if(mysocket){	 
		var username=document.getElementById("user").value;
		var password=document.getElementById("pswd").value;
		
			
		if (username=="") {
			document.getElementById("ack").value='用户名不能空';
		}else if (password=="") {
		   document.getElementById("ack").value='密码不能空'; 
		}else{
			
		
		auth_req();

		}
	}
	
};



	



  
function ack_HB(){
  	  var sentdata={};
  	  sentdata.cmd="HB";
      mysocket.send(JSON.stringify(sentdata));
}


function auth_req(){
	if(mysocket){
  	  var sentdata={};
  	  sentdata.cmd="auth_req";
	  sentdata.user=document.getElementById("user").value;
	  sentdata.pswd=document.getElementById("pswd").value;
	  
     	 mysocket.send(JSON.stringify(sentdata));
		 }
}


