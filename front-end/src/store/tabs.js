import { defineStore } from 'pinia'
 
export const useTabsStore = defineStore('tabs', {
  state() {
    return {
      tabsViewList: [], //tabsView数据存储
    }
  },
  actions: {
    // 添加导航标签方法
    addTabsView(tabItem) {
      let row = this.tabsViewList.find((v) => v.fullPath === tabItem.fullPath)
      if (!row) this.tabsViewList.push(tabItem)
    },
    // 删除导航标签方法
    removeTabsView(
      fullPath,
      isCurrent,
      callback
    ) {
      //如果不符合直接就删除
      if (isCurrent) {
        this.tabsViewList.forEach((item, index) => {
          if (item.fullPath === fullPath) {
            let navIndex =
              this.tabsViewList[index + 1] || this.tabsViewList[index - 1]
            if (navIndex) callback(navIndex.fullPath) //跳转到对应的路由
          }
        })
      }
      this.tabsViewList = this.tabsViewList.filter(
        (v) => v.fullPath !== fullPath,
      ) //删除面包屑每一项
    },
    // 关闭其他
    closeOtherTabsView(fullPath) {
      // noClosable是true和当前路由留下
      this.tabsViewList = this.tabsViewList.filter((item) => {
        return item.meta.noClosable || item.fullPath === fullPath
      })
    },
    // 关闭左侧Or右侧
    closeTabsViewOnSide(fullPath, type) {
      // 找到当前index
      let currentIndex = this.tabsViewList.findIndex(
        (item) => item.fullPath === fullPath,
      )
 
      // 判断一下必须存在才会执行
      if (currentIndex !== -1) {
        let range =
          type === 'left'
            ? [0, currentIndex]
            : [currentIndex + 1, this.tabsViewList.length]
 
        // 是左侧还是右侧 item.meta.noClosable固定留下的
        this.tabsViewList = this.tabsViewList.filter((item, index) => {
          return index < range[0] || index >= range[1] || item.meta.noClosable
        })
      }
    },
    //关闭全部
    closeAllTabsView() {
      this.tabsViewList = this.tabsViewList.filter((item) => {
        return item.meta.noClosable
      })
    },
    //从新设置tabsViewList
    setTabsViewList(list) {
      this.tabsViewList = list
    },
    // 设置标签栏标题
    setTabsViewTitle(fullPath, title) {
      this.tabsViewList.forEach((item) => {
        if (item.fullPath === fullPath) item.meta.title = title
      })
    },
  },
  getters: {},
})