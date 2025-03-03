import { qr_login } from '../api/api';

const QRLOGIN = 'QRLOGIN';
const CLEAR = 'CLEAR';
const initState = {
  tokenInfo: {},
  qrUserInfo: {}
}

export function qrlogin(state = initState, action) {
  switch (action.type) {
    case QRLOGIN:
      const { tokenInfo, qrUserInfo, ...others } = action.payload;
      return { ...state, ...others, tokenInfo, qrUserInfo };
    case CLEAR:
      return initState;
    default:
      return state;
  }
}

export function getQrLoginInfo() {
  return async dispatch => {
    const { tokenInfo, qrUserInfo } = await qr_login();
    dispatch({
      type: QRLOGIN,
      payload: {
        tokenInfo,
        qrUserInfo
      }
    })
  }
}

export function clearlist() {
  return async dispatch => {
    dispatch({
      type: CLEAR
    })
  }
}





