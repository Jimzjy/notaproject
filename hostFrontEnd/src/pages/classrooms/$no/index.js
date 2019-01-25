import React, { PureComponent } from "react";
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Page } from 'components'
import styles from './index.less'

@connect(({ classroomDetail }) => ({ classroomDetail }))
class ClassroomDetail extends PureComponent {
  render() {
    const { classroomDetail } = this.props
    const { data } = classroomDetail
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

ClassroomDetail.protoTypes = {
  classroomDetail: PropTypes.object,
}

export default ClassroomDetail