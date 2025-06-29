export function uuid() {
    var s = [];
    var hexDigits = "0123456789abcdef";
    for (var i = 0; i < 36; i++) {
        s[i] = hexDigits.substr(Math.floor(Math.random() * 0x10), 1);
    }
    s[14] = "4";  // bits 12-15 of the time_hi_and_version field to 0010
    s[19] = hexDigits.substr((s[19] & 0x3) | 0x8, 1);  // bits 6-7 of the clock_seq_hi_and_reserved to 01
    s[8] = s[13] = s[18] = s[23] = "-";
 
    var uuid = s.join("");
    return uuid;
}



export function nodeList(nodeData, list){
    for(const i in nodeData){
       const node =  nodeData[i];
       if(node.nodeType == 'start' || node.nodeType == 'between'   // || node.nodeType == 'parallel'  || node.nodeType == 'serial' 
           || node.nodeType == 'parallel-node'){
           list.push({
               nodeId: node.nodeId,
               nodeName: node.nodeName,
               nodeType: node.nodeType,
           });
       }
       
       const childNodes = node.childNodes;
       if(childNodes != null && childNodes.length > 0){
            for(const j in childNodes){
                nodeList(childNodes[j], list);
            }    
       }
        

    }
}


/**
 * 
 * @param {解析表单配置中的字段信息} parseFormJson 
 */
export function parseFormJson(formFields){
    const formColumns = [];
    formFields.forEach(element => {
        // 判断 input = true就行
        if(element.input){
            formColumns.push({
                label: element.label,
                field: element.field,
                required: element.rules?.[0].required
            })
        }else{
            if(element.children){
               return  formColumns.concat(parseFormJson(element.children))  
            }
           
        }
       
    });
    return formColumns;
}

export function changeNodeStatus(nodeData, hisList){
    // 修改节点状态
    
    nodeData.nodeStatus = []
    const hiss = hisList.filter(r => r.nodeCode == nodeData.nodeId);
    if(hiss && hiss.length > 0){
        for(const j in hiss){
            nodeData.nodeStatus.push(hiss[j].flowStatus);
        }
        
    }
    if(nodeData.childNode){
        changeNodeStatus(nodeData.childNode, hisList);
    }
 
   
    if(nodeData.conditionNodes){
        nodeData.conditionNodes.forEach(childNode => {
            changeNodeStatus(childNode, hisList);
        })
    }
   

}