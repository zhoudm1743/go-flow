import Layout from '~/components/layouts/index.vue';

import FlowConfig from '~/views/flowconfig/index.vue'

export const innerRouter = 
    { 
      component: Layout,
      children:[
        { path: '/flowconfig/:id', component: FlowConfig, name: "myconfig" , meta:{title: "流程配置"} },
  
      ] 
    }
   
  
