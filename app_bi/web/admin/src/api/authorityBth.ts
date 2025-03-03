
import service from '@/utils/request'

export const getAuthorityBtnApi = (data:any) => {
  return service({
    url: '/authorityBtn/getAuthorityBtn',
    method: 'post',
    data
  })
}

export const setAuthorityBtnApi = (data:any) => {
  return service({
    url: '/authorityBtn/setAuthorityBtn',
    method: 'post',
    data
  })
}

export const canRemoveAuthorityBtnApi = (params:any) => {
  return service({
    url: '/authorityBtn/canRemoveAuthorityBtn',
    method: 'post',
    params
  })
}

