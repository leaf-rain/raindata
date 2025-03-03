import axios from 'axios';
import { message } from 'antd';


axios.interceptors.request.use(config => {
  config = {
    ...config,
    withCredentials: false,
    baseURL: '',
  }
  return config
})

axios.interceptors.response.use(response => {
  return response
})
axios.defaults.baseURL = "http://127.0.0.1";


export async function qr_login() {
  return new Promise((resolve, reject) => {
    const appId =process.env.REACT_APP_ID;
    const redirect_uri = `http://120.76.192.134:3001/qrLogin/`;
    axios.post('http://127.0.0.1:3000/qrLogin/', {
      loginTime: `${new Date().getTime()}`,
      redirect_uri: redirect_uri,
      url: window.location.href,
    })
      .then(res => {
        var gotoUrl = `https://passport.feishu.cn/suite/passport/oauth/authorize?client_id=${appId}&redirect_uri=${encodeURI(redirect_uri)}&response_type=code&state=success_login`;
        var QRLoginObj = window.QRLogin({
          id: "login_container",
          goto: `${gotoUrl}`,
          style: "width: 300px; height: 300px; margin-left: 3em; border: 0; background-color: #E6F4F3; background-size: cover",
        });

        var handleMessage = function (event) {
          var origin = event.origin;
          if (QRLoginObj.matchOrigin(origin)
            && window.location.href.indexOf('/qrLogin') > -1
          ) {
            var loginTmpCode = event.data;
            window.location.href = `${gotoUrl}&tmp_code=${loginTmpCode}`;
          }
        };
        if (typeof window.addEventListener !== 'undefined') {
          window.addEventListener('message', handleMessage, false);
        }
        else if (typeof window.attachEvent !== 'undefined') {
          window.attachEvent('onmessage', handleMessage);
        }
        const { code, tokenInfo, qrUserInfo } = res.data;
        console.log(JSON.stringify(res))
        code === 0 ? resolve({ tokenInfo, qrUserInfo }) || message.success(res.data.msg) : message.info(res.data.msg);
      })
      .catch(err => {
        resolve(false)
      })
  })
}
