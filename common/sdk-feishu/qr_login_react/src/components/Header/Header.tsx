import React, { PureComponent } from 'react';
import { Link } from 'react-router-dom';
import myAvator from './avatar.svg';
import './header.less';

interface Props {
  active: number
}
interface State {
}
class Header extends PureComponent<Props, State>{
  constructor(props: Props) {
    super(props)
    this.state = {
    }
  }
  componentDidMount() {

  }
  render() {
    const { active } = this.props;
    const links = [{
      label: 'login',
      path: '/'
    }];
    return (

      <div className="Header">
        <div className="header-box fadein"><img src={myAvator} alt="photo" /> <p className="my-id">Open Platform Demo</p></div>
        <div></div>
        <div className="my-sort">
          {links.map((v, i) => (
            <Link to={v.path} key={i} className={`${active === i ? "active" : ""}`}>{v.label}</Link>
          ))}
        </div>
      </div>
    )
  }
}
export default Header;