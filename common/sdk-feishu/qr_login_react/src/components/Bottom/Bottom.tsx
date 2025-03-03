import React, { PureComponent } from 'react';
import './bottom.less';
class Bottom extends PureComponent<{}, {}>{
  // eslint-disable-next-line no-useless-constructor
  constructor(props: {}) {
    super(props)
  }
  render() {
    return (
      <div>
        <p className="beian">© LOCAL DEMO </p>
      </div>
    )
  }
}
export default Bottom;