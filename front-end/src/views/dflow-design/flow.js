export function backNodeList(nodeData){
    let backNodeList = [];
    return  getBackNodes(nodeData, backNodeList);
}

function getBackNodes(nodeData, backNodeList){
    const nodeType = nodeData.nodeType;

    if(nodeType == 'start' || nodeType == 'end' || nodeType == 'between'){
        backNodeList.push({
            nodeId: nodeData.nodeId,
            nodeName: nodeData.nodeName
        })
    }
    const conditionNodes = nodeData.conditionNodes;
    if(conditionNodes != null && conditionNodes.length > 0){
        for(const i in conditionNodes){
            getBackNodes(conditionNodes[i], backNodeList);
        }
    }
    const childNode = nodeData.childNode;
    if(childNode != null){
        getBackNodes(childNode, backNodeList);
    }

    return backNodeList;

}