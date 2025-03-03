import React, { PureComponent } from 'react';
import { connect } from 'react-redux';

import {
  Divider,
  Form,
  Input,
  Tooltip,
  Icon,
} from 'antd';

import Header from '../../components/Header/Header';
import Bottom from '../../components/Bottom/Bottom';
import { getQrLoginInfo } from '../../reducer/qrlogin.redux';
import './QrLogin.less';

interface Props {
  getQrLoginInfo: () => Record<string, string>,
  qrlogin: {
    tokenInfo: Record<string, string>,
    qrUserInfo: Record<string, string>,
  }
}
interface State {
  qrName: string,
}

// @ts-ignore
@connect(
  state => state,
  { getQrLoginInfo }
)
class QrLogin extends PureComponent<Props, State>{
  constructor(props: Props) {
    super(props)
    this.state = {
      qrName: 'login_container',
    }
  }
  componentDidMount() {
    this.props.getQrLoginInfo();
  }



  render() {

    const { qrlogin } = this.props;
    const { tokenInfo = {}, qrUserInfo = {} } = qrlogin;
  
    let { accessToken = '',refreshToken = '' } = tokenInfo;
    let { name = '', openId = '', userId = '', tenantKey = '', avatarUrl = ''} = qrUserInfo;

    const formItemLayout = {
      labelCol: {
        xs: { span: 24 },
        sm: { span: 8 },
      },
      wrapperCol: {
        xs: { span: 24 },
        sm: { span: 10 },
      },
    };

    return (
      <div>
        <Header active={2} />
        <Divider orientation="left" style={{ color: "#948C76", fontSize: "2em" }}>Step 1：</Divider>
        <p>You can embed the QR code in the webpage:</p>

        <div className="qr_div">
          <div id={this.state.qrName} ></div>
        </div>


        <Divider orientation="left" style={{ color: "#948C76", fontSize: "2em" }}>Step 2：</Divider>
        <p>Scan qr code to obtain user information:</p>

        <Form {...formItemLayout}>

          <Form.Item
            label={
              <span>
                access_token&nbsp;
                <Tooltip title="A unique credential to obtain user information.">
                  <Icon type="question-circle-o" />
                </Tooltip>
              </span>
            }
          >
            <Input value={accessToken} />
          </Form.Item>
          <Form.Item
            label={
              <span>
                refresh_token&nbsp;
                <Tooltip title="To refresh access token.">
                  <Icon type="question-circle-o" />
                </Tooltip>
              </span>
            }
          >
            <Input value={refreshToken} />
          </Form.Item>
          <Form.Item
            label={
              <span>
                name&nbsp;
                <Tooltip title="The name of the login user.">
                  <Icon type="question-circle-o" />
                </Tooltip>
              </span>
            }
          >
            <Input value={name} />
          </Form.Item>
          <Form.Item
            label={
              <span>
                open_id&nbsp;
              </span>
            }
          >
            <Input value={openId} />
          </Form.Item>
          <Form.Item
            label={
              <span>
                user_id&nbsp;
              </span>
            }
          >
            <Input value={userId} />
          </Form.Item>
          <Form.Item
            label={
              <span>
                tenant_key&nbsp;
                <Tooltip title="Unique identifier of a tenant.">
                  <Icon type="question-circle-o" />
                </Tooltip>
              </span>
            }
          >
            <Input value={tenantKey} />
          </Form.Item>

        </Form>


        <Bottom />
      </div>
    )
  }
}
export default QrLogin;