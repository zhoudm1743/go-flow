export default {
  mounted: (el, binding) => {
    el.resizeListener = () => {
      setHeight(el, binding)
    }
    setHeight(el, binding)
    window.addEventListener('resize', el.resizeListener)
  },
  unmounted(el) {
    window.removeEventListener('resize', el.resizeListener)
  },
  updated(el, binding) {
    setHeight(el, binding)
  }   
}

// set el-table height
function setHeight(el, binding) {
  const top = el.offsetTop
  const bottom = binding?.value?.bottom || 64
  const pageHeight = window.innerHeight

  el.style.height = pageHeight - top - bottom + 'px'
  el.style.overflowY = 'auto'  // 新增加
}
