
import NoMatch from '../containers/NoMatch/NoMatch';
import QrLogin from '../containers/QrLogin/QrLogin';

const routers = [{
  path: '/qrLogin',
  exact: false,
  component: QrLogin
},{
  path: '',
  exact: false,
  component: NoMatch
}];
export default routers;
