
// 查找节点的父元素
function familyTree(treeList, id) {
    // 返回数据集合
    var temp = [];
    // 声明digui函数
    var forFn = function (arr, id) {
      // 遍历树
      for (var i = 0; i < arr.length; i++) {
        var item = arr[i];
        if (item.id === id) {
          // 查找到指定节点加入集合
          temp.push(item);
          // 查找其父节点
          forFn(treeList, item.parentId);
          // 不必向下遍历，跳出循环
          break;
        } else {
          if (item.children) {
            // 向下查找到id
            forFn(item.children, id);
          }
        }
      }
    };
    // 调用函数
    forFn(treeList, id);
    // 返回结果
    return temp;
  }



export function getPositionInfo(position, deptTree) {
    
    return familyTree(deptTree, position.deptId);

}