import type { ProxyOptions } from 'vite';
import { createServiceConfig } from '../../src/utils/service';

export function createViteProxy(env: Env.ImportMeta, enable: boolean) {
    const isEnableHttpProxy = enable && env.VITE_HTTP_PROXY === 'true';
    if (!isEnableHttpProxy) return undefined;
    const {baseURL, proxyPattern, other} = createServiceConfig(env);
    const proxy: Record<string, ProxyOptions> = createProxyItem({baseURL, proxyPattern});
    other.forEach(element => {
        Object.assign(proxy, createProxyItem(element));
    });
    return proxy;
} 

function createProxyItem(item: App.Service.ServiceConfigItem) {
    const proxy: Record<string, ProxyOptions> = {};
    proxy[item.proxyPattern] = {
      target: item.baseURL,
      changeOrigin: true,
      rewrite: path => path.replace(new RegExp(`^${item.proxyPattern}`), '')
    };
    return proxy;
  }
  