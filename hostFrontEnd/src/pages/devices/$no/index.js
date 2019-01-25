import React, { PureComponent } from "react";
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Page } from 'components'
import styles from './index.less'

@connect(({ deviceDetail }) => ({ deviceDetail }))
class DeviceDetail extends PureComponent {
  render() {
    const { deviceDetail } = this.props
    const { data } = deviceDetail
    const content = []
    for (let key in data) {
      if ({}.hasOwnProperty.call(data, key)) {
        content.push(
          <div key={key} className={styles.item}>
            <div>{key}</div>
            <div>{String(data[key])}</div>
          </div>
        )
      }
    }
    return (
      <Page inner>
        <div className={styles.content}>{content}</div>
      </Page>
    )
  }
}

DeviceDetail.protoTypes = {
  deviceDetail: PropTypes.object,
}

export default DeviceDetail