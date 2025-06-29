<template>
    <div  style="display: flex; flex-direction: column; gap: 10px;">
        <el-card :body-style="{'padding-bottom': '0px'}">
            <el-form :inline="true" :model="formInline" class="demo-form-inline">
                <el-form-item label="Approved by">
                    <el-input v-model="formInline.user" placeholder="Approved by" clearable />
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="onSubmit">Query</el-button>
                </el-form-item>
            </el-form>
        </el-card>
        <el-card >
            
            <el-container>
                <el-aside width="200px" class="container-aside">   
                    <el-tree
                        style="max-width: 200px; "
                        :data="data"
                        :props="defaultProps"
                        @node-click="handleNodeClick"
                    />
                </el-aside>
                <el-main  class="container-main">
                    <el-table :data="tableData" style="width: 80%">
                        <el-table-column prop="date" label="Date" width="180" />
                        <el-table-column prop="name" label="Name" width="180" />
                        <el-table-column prop="address" label="Address" />
                    </el-table>
                </el-main>
            </el-container>

        </el-card>
    </div>
</template>

<script setup>
    import { reactive } from "vue"

    const formInline = reactive({
        user: '',
        region: '',
        date: '',
    })
    const tableData = [
    {
        date: '2016-05-03',
        name: 'Tom',
        address: 'No. 189, Grove St, Los Angeles',
    },
    {
        date: '2016-05-02',
        name: 'Tom',
        address: 'No. 189, Grove St, Los Angeles',
    },
    {
        date: '2016-05-04',
        name: 'Tom',
        address: 'No. 189, Grove St, Los Angeles',
    },
    {
        date: '2016-05-01',
        name: 'Tom',
        address: 'No. 189, Grove St, Los Angeles',
    },
    ]

    const data = [
        {
            label: 'Level one 1',
           
        },
        {
            label: 'Level one 2',
           
        },
        {
            label: 'Level one 3',
           
        },
        {
            label: 'Level one 4',
           
        },
        ]
</script>

<style  scoped>
    .container-aside{
        border: 1px solid var(--el-menu-border-color);
        margin-right: 5px;
        padding-top: 20px;
    }
    .container-main{
        border: 1px solid var(--el-menu-border-color);
    }
</style>