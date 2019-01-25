import React, { PureComponent } from "react";
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Page } from 'components'
import styles from './index.less'

@connect(({ cameraDetail }) => ({ cameraDetail }))
class CameraDetail extends PureComponent {
  render() {
    const { cameraDetail } = this.props
    const { data } = cameraDetail
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

CameraDetail.protoTypes = {
  cameraDetail: PropTypes.object,
}

export default CameraDetail