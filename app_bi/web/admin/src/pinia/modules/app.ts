
import { defineStore } from 'pinia'
import { ref, watchEffect, reactive } from 'vue'
import originSetting from "@/config.json"
import { setBodyPrimaryColor } from '@/utils/format'

export interface Config {
    weakness: boolean,
    grey: boolean,
    primaryColor: string,
    showTabs: boolean,
    darkMode: string | null,
    layout_side_width: number,
    layout_side_collapsed_width: number,
    layout_side_item_height: number,
    show_watermark: boolean,

    side_mode: string
}

export const useAppStore = defineStore('app', () => {
    const theme = ref(localStorage.getItem('theme') || originSetting.darkMode || 'auto')
    const device = ref("")
    const config = reactive({
        weakness: false,
        grey: false,
        primaryColor: '#79B6E7',
        showTabs: true,
        darkMode: 'light',
        layout_side_width: 256,
        layout_side_collapsed_width: 80,
        layout_side_item_height: 48,
        show_watermark: false,

        side_mode: 'normal'
    } as Config)

    // 初始化配置
    Object.assign(config, originSetting)

    if (config.primaryColor) {
        setBodyPrimaryColor(config.primaryColor, config.darkMode)
    }

    if (localStorage.getItem('darkMode')) {
        config.darkMode = localStorage.getItem('darkMode')
    }


    watchEffect(() => {
        if (theme.value === 'dark') {
            document.documentElement.classList.add('dark');
            document.documentElement.classList.remove('light');
            localStorage.setItem('theme', 'dark');
        } else {
            document.documentElement.classList.add('light');
            document.documentElement.classList.remove('dark');
            localStorage.setItem('theme', 'light');
        }
    })

    const toggleTheme = (dark:boolean) => {
        if (dark) {
            theme.value = 'dark';
        } else {
            theme.value = 'light';
        }
    }

    const toggleWeakness = (e:boolean) => {
        config.weakness = e;
        if (e) {
            document.documentElement.classList.add('html-weakenss');
        } else {
            document.documentElement.classList.remove('html-weakenss');
        }
    }

    const toggleGrey = (e:boolean) => {
        config.grey = e;
        if (e) {
            document.documentElement.classList.add('html-grey');
        } else {
            document.documentElement.classList.remove('html-grey');
        }
    }

    const togglePrimaryColor = (e:string) => {
        config.primaryColor = e;
        setBodyPrimaryColor(e, config.darkMode)
    }

    const toggleTabs = (e:boolean) => {
        config.showTabs = e;
    }

    const toggleDevice = (e:string) => {
        device.value = e;
    }

    const toggleDarkMode = (e:string) => {
        config.darkMode = e
        localStorage.setItem('darkMode', e)
        if (e === 'dark') {
            toggleTheme(true)
        } else {
            toggleTheme(false)
        }
    }

    const toggleDarkModeAuto = () => {
        // 处理浏览器主题
        const darkQuery = window.matchMedia('(prefers-color-scheme: dark)')
        const dark = darkQuery.matches
        toggleTheme(dark)
        darkQuery.addEventListener('change', (e) => {
            toggleTheme(e.matches)
        })
    }

    const toggleConfigSideWidth = (e:number) => {
        config.layout_side_width = e;
    }

    const toggleConfigSideCollapsedWidth = (e:number) => {
        config.layout_side_collapsed_width = e;
    }

    const toggleConfigSideItemHeight = (e:number) => {
        config.layout_side_item_height = e;
    }

    const toggleConfigWatermark = (e:boolean) => {
        config.show_watermark = e;
    }

    const toggleSideModel = (e:string) => {
        config.side_mode = e
    }

    if (config.darkMode === 'auto') {
        toggleDarkModeAuto()
    }

    toggleGrey(config.grey)

    return {
        theme,
        device,
        config,
        toggleTheme,
        toggleDevice,
        toggleWeakness,
        toggleGrey,
        togglePrimaryColor,
        toggleTabs,
        toggleDarkMode,
        toggleConfigSideWidth,
        toggleConfigSideCollapsedWidth,
        toggleConfigSideItemHeight,
        toggleConfigWatermark,
        toggleSideModel
    }
})
