import axios from 'axios'
import { ElMessage, ElLoading } from 'element-plus'

// 定义 API 返回的数据类型
export interface RequestBody {
    url: string;
    method: string;
    data: any;
    params: Map<string, any>;
}


const service = axios.create({
    baseURL: import.meta.env.VITE_BASE_API,
    timeout: 15000,
    headers: {
        "Content-Type": "application/json;charset=utf-8"
    }
})

let loading: any
// 正在请求的数量
let requestCount: number = 0
// 显示loding
const showLoading = () => {
    if (requestCount === 0 && !loading) {
        loading = ElLoading.service({
            text: "正在加载中",
            background: "rgba(0, 0, 0, 0.7)",
            spinner: "el-icon-loading",
        })
    }
    requestCount++
}

// 隐藏loading
const hideLoading = () => {
    requestCount--
    if (requestCount == 0) {
        loading.close()
    }
}

// 请求拦截
service.interceptors.request.use(config => {
    showLoading()
    // 是否需要设置 token
    // config.headers['Authorization'] = 'Bearer ' + getToken() // 让每个请求携带自定义token 请根据实际情况自行修改
    // get请求映射params参数
    if (config.method === 'get' && config.params) {
        let url = config.url + '?'
        for (const propName of Object.keys(config.params)) {
            const value = config.params[propName]
            var part = encodeURIComponent(propName) + "="
            if (value !== null && typeof (value) !== "undefined") {
                if (typeof value === 'object') {
                    for (const key of Object.keys(value)) {
                        let params = propName + '[' + key + ']'
                        var subPart = encodeURIComponent(params) + "="
                        url += subPart + encodeURIComponent(value[key]) + "&"
                    }
                } else {
                    url += part + encodeURIComponent(value) + "&"
                }
            }
        }
        url = url.slice(0, -1);
        config.params = {};
        config.url = url;
    }
    return config
}, err => {
    console.log(err)
    Promise.reject(err)
})

// 响应拦截器
service.interceptors.response.use((res: axios.AxiosResponse) => {
    hideLoading()
    // 未设置状态码则默认成功状态
    const code = res.data['code'] || 200;
    // 获取错误信息
    const msg = res.data['msg'] || '未知错误';
    if (code === 200 || code === 0) {
        return Promise.resolve(res.data)
    } else {
        ElMessage.error(msg)
        return Promise.reject(res.data)
    }
}, error => {
    console.log('err' + error)
    hideLoading()
    let { message } = error;
    if (message == "Network Error") {
        message = "后端接口连接异常";
    }
    else if (message.includes("timeout")) {
        message = "系统接口请求超时";
    }
    else if (message.includes("Request failed with status code")) {
        message = "系统接口" + message.substr(message.length - 3) + "异常";
    }
    ElMessage.error({
        message: message,
        duration: 5 * 1000
    })
    return Promise.reject(error)
})

export default service;