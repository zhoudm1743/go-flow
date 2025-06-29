<template>
      
        <el-drawer
        v-model="drawer"
        title="资源授权"
        direction="rtl"
        :before-close="handleClose"
    >
    <div  style="display: flex; flex-direction: column; gap: 10px;  padding-bottom: 15px;">
        <el-row>
            <el-button type="primary" @click="selectAll" size="small">全部选中</el-button>
            <el-button type="info" @click="unSelectAll" size="small">全部取消</el-button>
        </el-row>
        <el-row>
            <el-tree
                ref="treeRef"
                style="max-width: 600px"
                :data="data"
                show-checkbox
                :expand-on-click-node="false"
                default-expand-all
                node-key="id"
                check-strictly
                @check="handleCheck"
                :props="{class: customNodeClass,label: 'name' }"
            >
                
                <template #default="{ node, data }">
                    <span @click.stop="handleTextClick(node, data)"  class="custom-tree-node-content">
                        {{ node.label }}
                    </span>
                </template>
            </el-tree>
        </el-row>
    </div>  
        <template #footer>
            <el-button @click="drawer = false">取消</el-button>
            <el-button type="primary" @click="saveResource">确定</el-button>
        </template>
    </el-drawer>
  
    
</template>

<script setup>
import { onMounted, reactive, ref, watch } from 'vue';
import { allmenus } from '@/api/sys/menu';
import { addApi, listApi } from '@/api/sys/resource';
//  权限配置页面
const treeRef = ref(null);
const data = ref([]);

// 选中的数据，提交表单的数据格式
const selectData = ref([]);

const customNodeClass = (data, node) => {
  if (data.isPenultimate) {
    return 'is-penultimate'
  }
  if(data.isButton){
    return 'is-button'
  }
  return null
}
const drawer = ref(false)

const  emits = defineEmits(['tableReload']);

onMounted(()=>{
    allmenus().then(res => {
        data.value = res.rows;
    })
  
})
const fromInfo = reactive({
    fromType: '',
    fromId: ''
})


const openPage = (params) => {
    drawer.value = true;
    fromInfo.fromId = params.fromId;
    fromInfo.fromType = params.fromType;
    const selectKeys = [];
    listApi({fromId: params.fromId, fromType: params.fromType}).then(res => {
        res.rows.forEach(item => {
            selectKeys.push(item.resourceId)
        })
        treeRef.value.setCheckedKeys(selectKeys)
    })
    

}

const selectAll = () => {
    treeRef.value.setCheckedKeys(getAllKeys(data.value))
   
}
const unSelectAll = () => {
    treeRef.value.setCheckedKeys([])
}

// 全选或反全选监听
// watch(selectData, (newV, oldV) => {
//     const keys = [];
//     newV.forEach(item => {
//         keys.push(item.resourceId)
//     })
//     treeRef.value.setCheckedKeys(keys)
// }, {
//     deep: true,
// })

const saveResource = () => {
    const checkNodeList = treeRef.value.getCheckedNodes();
    const resourceList = [];
    for(let i = 0; i < checkNodeList.length; i++){
        resourceList.push({
            resourceId: checkNodeList[i].id,
            resourceType: checkNodeList[i].resourceType
        })
    }

    addApi({fromId: fromInfo.fromId, fromType: fromInfo.fromType, resourceList: resourceList}).then(res => {
        drawer.value = false;
    })
    
}
const getAllKeys = (nodes) => {
    let keys = [];
    nodes.forEach(node => {
    keys.push(node.id);
    if (node.children) {
        keys = keys.concat(getAllKeys(node.children));
    }
    });
    return keys;
};




const selectChildren =(data, checked) => {
    data.children.forEach(child => {
      treeRef.value.setChecked(child.id, checked);
      if (child.children) {
        selectChildren(child, checked);
      }
    });
  }
  const parentNodesChange = (node, checked) =>{
    if (node.parent) {
        treeRef.value.setChecked(node.parent.id, checked);
      parentNodesChange(node.parent, checked);
    }
  }

const handleTextClick = (node) => {
    const data = node.data;
    // 如果为选中状态， 点击文字就是取消选中
    if(node.checked){
        treeRef.value.setChecked(data.id, false);
        if (data.children) {
            //子节点也要选中
            selectChildren(data, false);
        }
    }else{
        treeRef.value.setChecked(data.id, true);
        parentNodesChange(data, true);
        if (data.children) {
            selectChildren(data, true);
        }
    }
}



function handleCheck(data, { checkedKeys }) {
    const selectedData = [];
    // 包含data.id 表示选中
    if(checkedKeys.includes(data.id)){
        // 父节点也要选中
        parentNodesChange(data, true);
        if (data.children) {
            //子节点也要选中
            selectChildren(data, true);
        }
    }else{
        if (data.children) {
                // 取消选中当前节点时，只取消子节点，不处理父节点
                selectChildren(data, false);
        }
    }
  
  }


defineExpose({
    openPage
})

</script>

<style  scoped>
    :deep .is-button > .el-tree-node__content {
    color: #626aef;
    }
    :deep  .custom-tree-node-content:hover {
        background-color: var(--el-tree-node-hover-bg-color);
    }
    :deep .el-tree-node__content ::selection {
        background-color: transparent;
        color: inherit;
       
    }

    :deep .el-tree-node__content {
        cursor: pointer !important;
    }
    :deep .el-tree .el-tree-node.is-penultimate > .el-tree-node__children {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
    }
    :deep .el-tree-node__content:hover {
        background-color: unset !important;
    }
    :deep .el-tree-node:focus>.el-tree-node__content {
        background-color: unset !important;
    }

 
</style>