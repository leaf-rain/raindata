import { fmtTitle } from '@/utils/fmtRouterTitle'
export default function getPageTitle(meta:any, route:any) {
    const titleField = ensureStringField(meta, 'title')
    const appname = import.meta.env.VITE_APP_NAME
    if (titleField !== '') {
        const title = fmtTitle(titleField, route)
        return `${title} - ${appname}`
    }
    return appname
}

function ensureStringField(data: any, fieldName: string): string {
    if (fieldName in data && typeof data[fieldName] === 'string') {
        return data[fieldName];
    }
    return '';
}

