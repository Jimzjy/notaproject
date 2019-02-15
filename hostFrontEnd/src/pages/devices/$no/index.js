import React, { PureComponent } from "react";
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Row, Collapse, Col, Card } from 'antd'
import { Page, SystemStats } from 'components'
import styles from './index.less'

@connect(({ deviceDetail }) => ({ deviceDetail }))
class DeviceDetail extends PureComponent {
  render() {
    const { deviceDetail } = this.props
    const { data, systemStats } = deviceDetail
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
            <Collapse.Panel header="设备信息" key="1">
              <div className={styles.content}>{content}</div>
            </Collapse.Panel>
          </Collapse>
        </Row>
        <Row>
          <Col lg={24} md={24} style={{ marginTop: 24 }}>
            <Card
              bodyStyle={{
                padding: '24px 36px 24px 0',
              }}
            >
              <SystemStats data={systemStats} />
            </Card>
          </Col>
        </Row>
      </Page>
    )
  }
}

DeviceDetail.protoTypes = {
  deviceDetail: PropTypes.object,
}

export default DeviceDetail