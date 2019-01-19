import React, { PureComponent } from "react";
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Button, Progress, Switch } from 'antd'
import { Stage, Layer, Rect, Text } from 'react-konva';
import Konva from 'konva';
import KImage from '../components/KImage'
import { Page } from 'components'
import styles from './index.less'
import { apiPrefix } from 'utils/config'

@connect(({ classDetail }) => ({ classDetail }))
class ClassDetail extends PureComponent {
  state = {
    showFaceCount: false,
    showFaceCountBackground: true,
    currentPerson: 0,
    faceCountData: {
      personCount: 0,
      backgroundImage: "",
      positionData: [],
      studentNos: [],
      global_height: 0,
      global_width: 0,
    }
  }

  handleFaceCountAdd = (data) => {
    const json = JSON.parse(data)
    const { detected_data, face_token, global_width, global_height, person_count, image_url } = json

    let _positionData = this.state.faceCountData.positionData
    _positionData.push(detected_data)
    let _studentNos = this.state.faceCountData.studentNos
    _studentNos.push(face_token)

    let _currentPerson = this.state.currentPerson + 1

    this.setState({
      showFaceCount: true,
      currentPerson: _currentPerson,
      faceCountData: {
        personCount: person_count,
        backgroundImage: image_url,
        positionData: _positionData,
        studentNos: _studentNos,
        global_height: global_height,
        global_width: global_width,
      }
    })

    console.log(this.state)
  }

  handleFaceCountStart = () => {
    const { classDetail } = this.props
    const { data } = classDetail
    const handleFaceCountAdd = this.handleFaceCountAdd

    this.setState({
      showFaceCount: false,
      currentPerson: 0,
      faceCountData: {
        personCount: 0,
        backgroundImage: "",
        positionData: [],
        studentNos: [],
        global_height: 0,
        global_width: 0,
      }
    })

    let ws = new WebSocket("ws://localhost:8000/face_count?class_id=" + data.class_id)
    ws.onopen = function(evt) {
      console.log("open: " + evt.data)
    }
    ws.onclose = function(evt) {
      console.log("close: " + evt.data)
    }
    ws.onmessage = function(evt) {
      console.log("message: " + evt.data)
      handleFaceCountAdd(evt.data)
    }
    ws.onerror = function(evt) {
      console.log("error: " + evt.data)
    }
    return
  }

  handleFaceCountBackground = (checked) => {
    this.setState({
      showFaceCountBackground: checked,
    })
  }

  render() {
    const { classDetail } = this.props
    const { data } = classDetail
    const { faceCountData, showFaceCount, showFaceCountBackground, currentPerson } = this.state
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

    const faceImage = `${apiPrefix}/images/${faceCountData.backgroundImage}`
    const imageScale = window.innerWidth * 0.7 / faceCountData.global_width
    const imageHeight = faceCountData.global_height * imageScale
    const imageWidth = faceCountData.global_width * imageScale

    const rects = []
    const texts = []
    const { positionData, studentNos } = faceCountData
    for (let index in faceCountData.positionData) {
      rects.push(
        <Rect
          key={index}
          x={positionData[index].x0 * imageScale} 
          y={positionData[index].y0 * imageScale}
          width={(positionData[index].x1 - positionData[index].x0) * imageScale}
          height={(positionData[index].y1 - positionData[index].y0) * imageScale}
          stroke={'#27F5F5'}
        />
      )

      texts.push(
        <Text 
          key={index}
          text={studentNos[index]} 
          fontSize={imageWidth / 50} 
          x={positionData[index].x0 * imageScale} 
          y={positionData[index].y0 * imageScale} 
          fill={'#6CC4C4'}
        />
      )
    }

    return (
      <Page inner>
        <div className={styles.content}>{content}</div>
        <Button onClick={this.handleFaceCountStart}>点名</Button>
        <Switch defaultChecked  checkedChildren="开" unCheckedChildren="关" onChange={this.handleFaceCountBackground} />
        {showFaceCount && (
          <div>
            <Progress size="small" percent={Math.ceil(currentPerson * 100 / faceCountData.personCount)} />
            <Stage width={window.innerWidth} height={imageHeight + 20}>
              <Layer>
                { showFaceCountBackground && <KImage src={faceImage} width={imageWidth} height={imageHeight}/> }
                { rects }
                { texts }
              </Layer>
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