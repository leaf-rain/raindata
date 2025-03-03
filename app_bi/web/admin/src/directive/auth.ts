// 权限按钮展示指令
import { useUserStore } from '@/pinia/modules/user'
import type { App, DirectiveBinding } from 'vue'

export default {
  install: (app: App) => {
    const userStore = useUserStore()
    app.directive('auth', {
      mounted: function(el: HTMLElement, binding: DirectiveBinding) {
        const userInfo = userStore.userInfo
        if (!binding.value) {
          el.parentNode?.removeChild(el)
          return
        }

        const waitUse = Array.isArray(binding.value) ? binding.value : [binding.value]
        let flag = waitUse.some(item => Number(item) === userInfo.authorityId)
        if (binding.modifiers.not) {
          flag = !flag
        }
        if (!flag) {
          el.parentNode?.removeChild(el)
        }
      }
    })
  }
}