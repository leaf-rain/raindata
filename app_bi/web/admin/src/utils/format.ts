import { formatTimeToStr } from '@/utils/date';
import { ref } from 'vue';

// 定义 DictItem 类型
interface DictItem {
  value: any;
  label: string;
}

export const formatBoolean = (bool: boolean) => {
  if (bool !== null) {
    return bool ? '是' : '否';
  } else {
    return '';
  }
};

export const formatDate = (time: Date) => {
  return formatTimeToStr(time, 'yyyy-MM-dd hh:mm:ss');
};

export const filterDict = (value: any, options: DictItem[] | undefined) => {
  const rowLabel = options?.filter((item: DictItem) => item.value === value);
  return rowLabel && rowLabel[0] && rowLabel[0].label;
};

export const filterDataSource = (dataSource: DictItem[], value: any) => {
  if (Array.isArray(value)) {
    return value.map((item) => {
      const rowLabel = dataSource.find((i: DictItem) => i.value === item);
      return rowLabel?.label;
    });
  }
  const rowLabel = dataSource.find((item: DictItem) => item.value === value);
  return rowLabel?.label;
};

const path = import.meta.env.VITE_BASE_PATH + ':' + import.meta.env.VITE_SERVER_PORT + '/';
export const ReturnArrImg = (arr: any) => {
  const imgArr: string[] = [];
  if (arr instanceof Array) {
    for (const arrKey in arr) {
      if (arr[arrKey].slice(0, 4) !== 'http') {
        imgArr.push(path + arr[arrKey]);
      } else {
        imgArr.push(arr[arrKey]);
      }
    }
  } else {
    if (arr.slice(0, 4) !== 'http') {
      imgArr.push(path + arr);
    } else {
      imgArr.push(arr);
    }
  }
  return imgArr;
};

export const onDownloadFile = (url: string) => {
  window.open(path + url);
};


const colorToHex = (u: string): string[] => {
    let e = u.replace('#', '').match(/../g);
    if (!e) {
      throw new Error('Invalid color format');
    }
    return e.map((hexPart) => parseInt(hexPart, 16).toString(16));
  };

const hexToColor = (u: number, e: number, t: number) => {
  let a = [u.toString(16), e.toString(16), t.toString(16)];
  for (let n = 0; n < 3; n++) {
    a[n].length === 1 && (a[n] = `0${a[n]}`);
  }
  return `#${a.join('')}`;
};

const generateAllColors = (u: string, e: number) => {
    let t = colorToHex(u);
    const target = [10, 10, 30];
    for (let a = 0; a < 3; a++) {
      const newValue = Math.floor(parseInt(t[a], 16) * (1 - e) + target[a] * e);
      t[a] = newValue.toString(16).padStart(2, '0'); // 确保每个部分都是两位数的十六进制字符串
    }
    return hexToColor(parseInt(t[0], 16), parseInt(t[1], 16), parseInt(t[2], 16));
  };


  const generateAllLightColors = (u: string, e: number) => {
    let t = colorToHex(u);
    const target = [240, 248, 255]; // RGB for blue white color
    for (let a = 0; a < 3; a++) {
      const newValue = Math.floor(parseInt(t[a], 16) * (1 - e) + target[a] * e);
      t[a] = newValue.toString(16).padStart(2, '0'); // 确保每个部分都是两位数的十六进制字符串
    }
    return hexToColor(parseInt(t[0], 16), parseInt(t[1], 16), parseInt(t[2], 16));
  };

function addOpacityToColor(u: string, opacity: number) {
  let t = colorToHex(u);
  return `rgba(${t[0]}, ${t[1]}, ${t[2]}, ${opacity})`;
}

export const setBodyPrimaryColor = (primaryColor: string, darkMode: string | null) => {
  let fmtColorFunc = generateAllColors;
  if (darkMode === 'light') {
    fmtColorFunc = generateAllLightColors;
  }

  document.documentElement.style.setProperty('--el-color-primary', primaryColor);
  document.documentElement.style.setProperty('--el-color-primary-bg', addOpacityToColor(primaryColor, 0.4));
  for (let times = 1; times <= 2; times++) {
    document.documentElement.style.setProperty(`--el-color-primary-dark-${times}`, fmtColorFunc(primaryColor, times / 10));
  }
  for (let times = 1; times <= 10; times++) {
    document.documentElement.style.setProperty(`--el-color-primary-light-${times}`, fmtColorFunc(primaryColor, times / 10));
  }
  document.documentElement.style.setProperty('--el-menu-hover-bg-color', addOpacityToColor(primaryColor, 0.2));
};

const baseUrl = ref(import.meta.env.VITE_BASE_API);

export const getBaseUrl = () => {
  return baseUrl.value === '/' ? '' : baseUrl.value;
};