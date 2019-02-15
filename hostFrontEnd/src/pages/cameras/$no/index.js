import React, { PureComponent } from "react";
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Page } from 'components'
import styles from './index.less'
import { Row, Collapse } from 'antd'

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
        <Row>
          <Collapse defaultActiveKey={['1']} >
            <Collapse.Panel header="摄像头信息" key="1">
              <div className={styles.content}>{content}</div>
            </Collapse.Panel>
          </Collapse>
        </Row>
      </Page>
    )
  }
}

CameraDetail.protoTypes = {
  cameraDetail: PropTypes.object,
}

export default CameraDetail