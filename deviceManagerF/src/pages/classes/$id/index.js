import React, { PureComponent } from "react";
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Button } from 'antd'
import { Page } from 'components'
import styles from './index.less'
import { apiPrefix } from 'utils/config'
import Konva from 'konva';
import { Stage, Layer, Rect, Text, Circle, Line } from 'react-konva';

@connect(({ classDetail }) => ({ classDetail }))
class ClassDetail extends PureComponent {
  state = {
    showFaceCount: false,
    faceCountData: {
      showBackground: true,
      personCount: 0,
      backgroundImage: "",
      personData: []
    }
  }

  componentDidMount() {
    this.setState()
  }

  handleFaceCountStart = () => {
    const { classDetail } = this.props
    const { data } = classDetail

    this.setState({
      showFaceCount: true,
      faceCountData: {
        showBackground: true,
        personCount: 0,
        backgroundImage: "",
        personData: []
      }
    })

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

      const { person_count, image_url, detected_data } = evt.data
      let _personData = this.state.faceCountData.personData
      _personData.push(detected_data)

      this.setState({
        faceCountData: {
          personCount: person_count,
          backgroundImage: image_url,
          personData: _personData,
        }
      })
    }
    ws.onerror = function(evt) {
      console.log("error: " + evt.data)
    }
    return
  }

  render() {
    const { classDetail } = this.props
    const { data } = classDetail
    const { showFaceCount, faceCountData } = this.state
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
    const faceImage = faceCountData.showBackground ? `${apiPrefix}/images/${faceCountData.backgroundImage}` : ''
    return (
      <Page inner>
        <div className={styles.content}>{content}</div>
        <Button onClick={this.handleFaceCountStart}>点名</Button>
        {showFaceCount && (
          <div style={{ backgroundImage: faceImage }}>
            <Stage>

            </Stage>
          </div>
        )}
      </Page>
    )
  }
}

ClassDetail.protoTypes = {
  studentDetail: PropTypes.object,
}

export default ClassDetail