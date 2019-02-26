import React, { PureComponent } from "react";
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Row, Collapse } from 'antd'
import { Page } from 'components'
import styles from './index.less'

@connect(({ classDetail }) => ({ classDetail }))
class ClassDetail extends PureComponent {
  render() {
    const { classDetail } = this.props
    const { data } = classDetail
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
            <Collapse.Panel header="班级信息" key="1">
              <div className={styles.content}>{content}</div>
            </Collapse.Panel>
          </Collapse>
        </Row>
      </Page>
    )
  }
}

ClassDetail.protoTypes = {
  classDetail: PropTypes.object,
}

export default ClassDetail