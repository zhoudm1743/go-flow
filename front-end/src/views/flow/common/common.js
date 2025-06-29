import { defineAsyncComponent, markRaw } from 'vue'

const modules = import.meta.glob('/src/views/customform/**.vue')


export const openCustomForm = (def) => {
    const formPath= getFormpath(def);
    if(formPath == null || formPath == ""){ 
        return
    }
    const formPathStr = `/src/views/customform/${formPath}`;
    return  defineAsyncComponent(modules[formPathStr]) ;
}

const  getFormpath = (def) => {
    // 判断为设计表单
    if(def.formCustom == 'N'){
        return 'designForm.vue'
    }
    return def.formPath;
}