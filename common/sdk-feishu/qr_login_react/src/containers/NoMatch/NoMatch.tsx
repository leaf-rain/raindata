import React, { PureComponent } from 'react';
import { Row, Col } from 'antd';
import Header from '../../components/Header/Header';
import Bottom from '../../components/Bottom/Bottom';
import {Redirect} from  'react-router-dom'
import './nomatch.less'

interface State {
}
interface Props {
}
class NoMatch extends PureComponent<Props, State>{
  // eslint-disable-next-line no-useless-constructor
  constructor(props: Props) {
    super(props)


  }
  render() {
    return <div>
      <Redirect from={"/"} to={'/qrLogin'}/>
      <Header active={0} />
      <Row gutter={16} className="rowStyle">
        <Col span={14}>
          <img width="100%" src={require('./none.svg')} alt="none" />
        </Col>
        <Col span={10}>
          <div className="statusStyle"></div>
          <div className="smallStyle">redirect...</div>
        </Col>
      </Row>
      <Bottom />
    </div>
  }
}

export default NoMatch