import React, { PureComponent } from "react";
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Page } from 'components'
import styles from './index.less'
import { Row, Collapse } from 'antd'

@connect(({ studentDetail }) => ({ studentDetail }))
class StudentDetail extends PureComponent {
  render() {
    const { studentDetail } = this.props
    const { data } = studentDetail
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
            <Collapse.Panel header="学生信息" key="1">
              <div className={styles.content}>{content}</div>
            </Collapse.Panel>
          </Collapse>
        </Row>
      </Page>
    )
  }
}

StudentDetail.protoTypes = {
  studentDetail: PropTypes.object,
}

export default StudentDetail