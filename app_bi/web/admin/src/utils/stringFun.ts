/* eslint-disable */
export const toUpperCase = (str:string) => {
    if (str[0]) {
        return str.replace(str[0], str[0].toUpperCase())
    } else {
        return ''
    }
}

export const toLowerCase = (str:string) => {
    if (str[0]) {
        return str.replace(str[0], str[0].toLowerCase())
    } else {
        return ''
    }
}

// 驼峰转换下划线
export const toSQLLine = (str:string) => {
    if (str === 'ID') return 'ID'
    return str.replace(/([A-Z])/g, "_$1").toLowerCase();
}

// 下划线转换驼峰
export const toHump = (str:string) => {
    return str.replace(/\_(\w)/g, function(all, letter) {
        return letter.toUpperCase();
    });
}

export const convertToString = (input: any): string => {
    if (typeof input === 'string' || typeof input === 'number' || typeof input === 'boolean' || input === null || input === undefined) {
        return String(input);
    }
    if (typeof input === 'object' && input !== null && typeof input.toString === 'function') {
        return input.toString();
    }
    return JSON.stringify(input);
}