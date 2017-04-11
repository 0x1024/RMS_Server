document.getElementById("pid").value="10000"
document.getElementById("dtype").value="新机型"
document.getElementById("client").value="研发"
document.getElementById("tags").value="设计中"
document.getElementById("passwd").value="123"
document.getElementById("created").value=robot.getDate("-")+" "+robot.getTime() ;
document.getElementById("updated").value=robot.getDate("-")+" "+robot.getTime() ;

var mysocket = new WebSocket("ws://10.1.11.151:9003");
	mysocket.onopen=function(){
		console.log("连接成功");
		init_req();

	}
  mysocket.onmessage = function(data) {
 	  		console.log(data);
     var msg = JSON.parse(data.data);
      var data=msg.data;
      // if (msg.cmd == "data_single") {
      //     var contStr = "";
      //     for (var i = 0; i < data.length; i++) {
      //         contStr += "<tr><td>" + data[i].Pid + "</td><td>" + data[i].Dtype + "</td><td>" + data[i].Client + "</td><td>" + data[i].Tags + "</td><td>" + data[i].Passwd + "</td><td>" + data[i].Created + "</td><td>" + data[i].Updated + "</td></tr>";
      //     }
      //     document.getElementById("tcont").innerHTML = contStr;
      // }

      switch(msg.cmd){
          case "data_single":
          var contStr = "";
  
              contStr += "<tr><td>" + data.Pid + "</td><td>" + data.Dtype + "</td><td>" + data.Client + "</td><td>" + data.Tags + "</td><td>" + data.Passwd + "</td><td>" + data.Created + "</td><td>" + data.Updated + "</td></tr>";

          document.getElementById("tcont").innerHTML = contStr;
          break;

           case "HB":
			ack_HB()
          break;

          case "data_all":
          var contStr = "";
          for (var i = 0; i < data.length; i++) {
              contStr += "<tr><td>" + data[i].Pid + "</td><td>" + data[i].Dtype + "</td><td>" + data[i].Client + "</td><td>" + data[i].Tags + "</td><td>" + data[i].Passwd + "</td><td>" + data[i].Created + "</td><td>" + data[i].Updated + "</td></tr>";
          }
          document.getElementById("tcont").innerHTML = contStr;
          break;
          default:
          break;

     
      }

  }

function ack_HB(){
  	  var sentdata={};
  	  sentdata.cmd="HB";
      mysocket.send(JSON.stringify(sentdata));
}

  function req(){
  	  var sentdata={};
  	  sentdata.cmd="req";
	  sentdata.pid=document.getElementById("Req ID").value;
     	 mysocket.send(JSON.stringify(sentdata));
}

  function delete_id(){
  	  var sentdata={};
  	  sentdata.cmd="delete_id";
	  sentdata.pid=document.getElementById("delete_id").value;
     	 mysocket.send(JSON.stringify(sentdata));
}

  function init_req(){
  	  var sentdata={};
  	  sentdata.cmd="all";

     	 mysocket.send(JSON.stringify(sentdata));
}

  function submit(){
  	  var sentdata={};
 	  sentdata.cmd="comitone";		
  	  sentdata.pid=document.getElementById("pid").value;
  	  sentdata.dtype=document.getElementById("dtype").value;
  	  sentdata.client=document.getElementById("client").value;
  	  sentdata.tags=document.getElementById("tags").value;
  	  sentdata.passwd=document.getElementById("passwd").value;
  	  sentdata.created=document.getElementById("created").value;
  	  sentdata.updated=document.getElementById("updated").value;
  	  
      mysocket.send(JSON.stringify(sentdata));
  }