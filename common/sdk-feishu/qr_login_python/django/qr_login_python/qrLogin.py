import os
import json
import requests
from urllib import parse
from django.http import JsonResponse

# get env
APP_ID = "cli_a72d9f7014f89013"
APP_SECRET = "GxB2VAGUYx93eNSHIyZOzfyxP3hlrfBT"
LARK_PASSPORT_HOST = "https://passport.feishu.cn/suite/passport/oauth/"
LARK_BASE_URL = "https://open.feishu.cn/open-apis"


def qr_login(request):
    if request.method == "POST":
        request_body = request.body
        # init config
        qrLogin = QrLogin(APP_ID, APP_SECRET, LARK_PASSPORT_HOST, LARK_BASE_URL)
        print("body:", request_body)
        # get token
        tokenInfo = qrLogin.get_token_info(json_param=json.loads(request_body.decode()))

        # get user
        qrUserInfo = qrLogin.get_user_info()

        if not tokenInfo and qrUserInfo:
            response = {"code": -1, "msg": "Scan qr code to obtain user information."}
            return JsonResponse(response, safe=False)

        response = {"code": 0, "msg": "get userinfo success", "tokenInfo": tokenInfo,
                    "qrUserInfo": qrUserInfo}

        return JsonResponse(response, safe=False)


class QrLogin(object):
    def __init__(self, app_id, app_secret, lark_passport_host, lark_base_url):
        self.lark_passport_host = lark_passport_host
        self.lark_base_url = lark_base_url
        self.app_id = app_id
        self.app_secret = app_secret
        self._token_info = {}
        self._user_info = {}

    def app_access_token(self):
        param = {
                    "app_id": self.app_id,
                    "app_secret": self.app_secret
                }
        headers = {"Content-Type": "application/json; charset=utf-8"}
        token_res = json.loads(requests.post(self._gen_url(uri="/auth/v3/app_access_token/internal"), param, headers).text)
        return token_res.get("tenant_access_token")

    def get_token_info(self, json_param):

        url_param_result = parse.urlparse(json_param.get("url"))
        query_dict = parse.parse_qs(url_param_result.query)
        code_arr = query_dict.get("code", [])
        if not code_arr:
            return {}
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": "Bearer " + self.app_access_token()
        }
        payload = json.dumps({
            "grant_type": "authorization_code",
            "code": code_arr[0]
        })
        token_res = json.loads(requests.post(self._gen_url(uri="/authen/v1/oidc/access_token"), data=payload, headers=headers).text)
        tokenInfo = {}
        if token_res.get("code") == 0:
            token_res = token_res.get("data")
            tokenInfo = {"accessToken": token_res.get("access_token"), "refreshToken": token_res.get("refresh_token"),
                        "tokenType": token_res.get("token_type")}
        print(token_res)
        self._token_info = tokenInfo
        return tokenInfo

    def get_user_info(self):
        header = {}
        response = {}
        qrUserInfo = {}

        header["Content-Type"] = "application/json;charset=UTF-8"
        header["Authorization"] = "%s %s" % (self._token_info.get("tokenType"), self._token_info.get("accessToken"))
        try:
            qr_login_user = requests.get(url=self._gen_url(uri="/authen/v1/user_info"), headers=header).text
        except Exception as e:
            print(e)
            return {}

        userInfoObj = json.loads(qr_login_user)
        if userInfoObj.get("code") == 0:
            userInfoObj = userInfoObj.get("data")
            qrUserInfo["name"] = userInfoObj.get("name")
            qrUserInfo["openId"] = userInfoObj.get("open_id")
            qrUserInfo["userId"] = userInfoObj.get("open_id")
            qrUserInfo["tenantKey"] = userInfoObj.get("tenant_key")
            qrUserInfo["avatarUrl"] = userInfoObj.get("avatar_url")
        self._user_info = qrUserInfo
        return qrUserInfo

    def _gen_url(self, uri):
        return "{}{}".format(self.lark_base_url, uri)

