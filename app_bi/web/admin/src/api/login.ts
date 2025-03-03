import request from "../utils/request";
import type {LoginResponse} from "./api";

export function login(data: any): Promise<LoginResponse> {
    // return request({
    //     url: import.meta.env.VITE_BASE_API+'login/',
    //     method: 'post',
    //     data: data,
    // }) as Promise<LoginResponse>;
    return new Promise((resolve, reject) => {
        resolve({
            code: 0,
            msg: 'success',
            data: {
                token: '1234567890',
                user: {
                    username: 'admin',
                    nickname: '管理员',
                    avatar: 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif',
                }
            }
        })
    })
}
