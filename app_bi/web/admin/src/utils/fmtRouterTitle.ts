export const fmtTitle = (title: string, route: any) => {
    const reg = /\$\{(.+?)\}/;
    const reg_g = /\$\{(.+?)\}/g;
    const result = title.match(reg_g);
    if (result) {
        result.forEach((item) => {
            const matchResult = item.match(reg);
            if (matchResult && matchResult[1]) {
                const key = matchResult[1];
                const value = route.params?.[key] || route.query?.[key] || '';
                title = title.replace(item, value);
            }
        });
    }
    return title;
}