

export function addEventListen(
    target:any,
    event:any,
    handler: Function,
    capture = false
  ) {
    if (
      target.addEventListener &&
      typeof target.addEventListener === 'function'
    ) {
      target.addEventListener(event, handler, capture);
    }
  }
  
  export function removeEventListen(
    target:any,
    event:any,
    handler: Function,
    capture = false
  ) {
    if (
      target.removeEventListener &&
      typeof target.removeEventListener === 'function'
    ) {
      target.removeEventListener(event, handler, capture);
    }
  }
  
  
  