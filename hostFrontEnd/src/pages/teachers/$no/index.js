import React, { PureComponent } from "react";
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Page } from 'components'
import styles from './index.less'
import { Row, Collapse } from 'antd'

@connect(({ teacherDetail }) => ({ teacherDetail }))
class TeacherDetail extends PureComponent {
  render() {
    const { teacherDetail } = this.props
    const { data } = teacherDetail
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
            <Collapse.Panel header="教师信息" key="1">
              <div className={styles.content}>{content}</div>
            </Collapse.Panel>
          </Collapse>
        </Row>
      </Page>
    )
  }
}

TeacherDetail.protoTypes = {
  teacherDetail: PropTypes.object,
}

export default TeacherDetail