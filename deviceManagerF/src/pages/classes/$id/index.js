import React, { PureComponent } from "react";
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Button } from 'antd'
import { pathMatchRegexp } from 'utils'
import { Page } from 'components'
import styles from './index.less'

@connect(({ classDetail }) => ({ classDetail }))
class ClassDetail extends PureComponent {
  handleFaceCountStart = () => {
    const { classDetail } = this.props
    const { data } = classDetail

    let ws = new WebSocket("ws://localhost:8000/face_count?class_id=" + data.class_id)
    ws.onopen = function(evt) {
      console.log("open: " + evt.data)
    }
    ws.onclose = function(evt) {
      console.log("close: " + evt.data)
      ws = null;
    }
    ws.onmessage = function(evt) {
      console.log("message: " + evt.data)
    }
    ws.onerror = function(evt) {
      console.log("error: " + evt.data)
    }
    return
  }

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
        <div><Button onClick={this.handleFaceCountStart}>点名</Button></div>
      </Page>
    )
  }
}

ClassDetail.protoTypes = {
  studentDetail: PropTypes.object,
}

export default ClassDetail