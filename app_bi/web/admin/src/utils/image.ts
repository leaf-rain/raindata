export default class ImageCompress {
    file: File;
    fileSize: number;
    maxWH: number;

    constructor(file: File, fileSize: number, maxWH: number = 1920) {
        this.file = file;
        this.fileSize = fileSize;
        this.maxWH = maxWH; // 最大长宽
    }

    compress() {
        // 压缩逻辑保持不变
        const fileType = this.file.type;
        const fileSize = this.file.size / 1024;
        return new Promise<File>((resolve) => {
            const reader = new FileReader();
            reader.readAsDataURL(this.file);
            reader.onload = () => {
                const canvas = document.createElement('canvas');
                const img = document.createElement('img');
                img.src = reader.result as string;
                img.onload = () => {
                    const ctx = canvas.getContext('2d')!;
                    const _dWH = this.dWH(img.width, img.height, this.maxWH);
                    canvas.width = _dWH.width;
                    canvas.height = _dWH.height;

                    // 清空后, 重写画布
                    ctx.clearRect(0, 0, canvas.width, canvas.height);
                    ctx.drawImage(img, 0, 0, canvas.width, canvas.height);

                    const newImgData = canvas.toDataURL(fileType, 0.90);

                    // 压缩宽高后的图像大小
                    const newImgSize = this.fileSizeKB(newImgData);

                    if (newImgSize > this.fileSize) {
                        console.log('图片尺寸太大!' + fileSize + ' >> ' + newImgSize);
                    }

                    const blob = this.dataURLtoBlob(newImgData, fileType);
                    const nfile = new File([blob], this.file.name, { type: fileType });
                    resolve(nfile);
                };
            };
        });
    }

    /**
     * 长宽等比缩小
     * 图像的一边(长或宽)为最大目标值
     */
    dWH(srcW: number, srcH: number, dMax: number) {
        const defaults = {
            width: srcW,
            height: srcH
        };
        if (Math.max(srcW, srcH) > dMax) {
            if (srcW > srcH) {
                defaults.width = dMax;
                defaults.height = Math.round(srcH * (dMax / srcW));
                return defaults;
            } else {
                defaults.height = dMax;
                defaults.width = Math.round(srcW * (dMax / srcH));
                return defaults;
            }
        } else {
            return defaults;
        }
    }

    fileSizeKB(dataURL: string): number {
        let sizeKB = 0;
        sizeKB = Math.round((dataURL.split(',')[1].length * 3 / 4) / 1024);
        return sizeKB;
    }

    /**
    * 转为Blob
    */
   dataURLtoBlob(dataURL: string, fileType?: string): Blob {
       const byteString = atob(dataURL.split(',')[1]);
       let mimeString = dataURL.split(',')[0].split(':')[1].split(';')[0];
       const ab = new ArrayBuffer(byteString.length);
       const ia = new Uint8Array(ab);
       for (let i = 0; i < byteString.length; i++) {
           ia[i] = byteString.charCodeAt(i);
       }
       if (fileType) {
           mimeString = fileType;
       }
       return new Blob([ab], { type: mimeString });
   }
}

const path = import.meta.env.VITE_FILE_API;
export const getUrl = (url: string): string => {
    if (url && url.slice(0, 4) !== 'http') {
        if (path === '/') {
            return url;
        }
        if (url.slice(0, 1) === '/') {
            return path + url;
        }
        return path + '/' + url;
    } else {
        return url;
    }
};

export const isVideoExt = (url: string): boolean =>
    url.endsWith('.mp4') || url.endsWith('.mov') || url.endsWith('.webm') || url.endsWith('.ogg');

export const isVideoMime = (type: string): boolean =>
    type === 'video/mp4' || type === 'video/webm' || type === 'video/ogg';

export const isImageMime = (type: string): boolean =>
    type === 'image/jpeg' || type === 'image/png' || type === 'image/webp' || type === 'image/svg+xml';