<template>
     <EBuilder  ref="ebForm"  :formData="variable" :pageSchema="data" :disabled="readonly" />
</template>

<script setup>
//  设计表单容器
    import { ref } from "vue";
    import { EBuilder } from "epic-designer";
    defineProps({
        data: {
            type: Object,
            default: () => {},
            required: true
        },
        variable: {
            type: Object,
            default: () => {},
            required: true
        },
        readonly: { 
            type: Boolean,
            default: false,
        }
    })

    const ebForm = ref(null)


    const submitForm = async () => {
        return new Promise((resolve, reject) => {
            ebForm.value.validate().then(res => {
                resolve(res)
            }).catch(err => {
                reject(err)
            })
        })
 
    }
    defineExpose({submitForm})


</script>

<style  scoped>

</style>