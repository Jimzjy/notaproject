import { PureComponent } from "react";
import PropTypes from 'prop-types'
import { connect } from 'dva'
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
        <div className={styles.content}>{content}</div>
      </Page>
    )
  }
}

ClassDetail.protoTypes = {
  studentDetail: PropTypes.object,
}

export default ClassDetail