import React, { PureComponent } from "react";
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Button, Progress, Switch, Spin, Row, Col } from 'antd'
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
    showSpinning: false,
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
    const { face, global_width, global_height, person_count, image_url } = json

    let _positionData = this.state.faceCountData.positionData
    _positionData.push(face.face_rectangle)
    let _studentNos = this.state.faceCountData.studentNos
    _studentNos.push(face.face_token)

    let _currentPerson = this.state.currentPerson + 1

    this.setState({
      showFaceCount: true,
      showSpinning: false,
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
  }

  handleFaceCountStart = () => {
    const { classDetail } = this.props
    const { data } = classDetail
    const handleFaceCountAdd = this.handleFaceCountAdd
    const handleFaceCountClose = this.handleFaceCountClose

    this.setState({
      showFaceCount: false,
      showSpinning: true,
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
      console.log("open websocket")
    }
    ws.onclose = function(evt) {
      console.log("close websocket")
      handleFaceCountClose()
    }
    ws.onmessage = function(evt) {
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

  handleFaceCountClose = () => {
    this.setState({
      showSpinning: false,
    })
  }

  render() {
    const { classDetail } = this.props
    const { data } = classDetail
    const { faceCountData, showFaceCount, showFaceCountBackground, currentPerson, showSpinning } = this.state
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
          x={positionData[index].left * imageScale} 
          y={positionData[index].top * imageScale}
          width={positionData[index].width * imageScale}
          height={positionData[index].height * imageScale}
          stroke={'#27F5F5'}
        />
      )

      texts.push(
        <Text 
          key={index}
          text={studentNos[index]} 
          fontSize={imageWidth / 70} 
          x={positionData[index].left * imageScale} 
          y={positionData[index].top * imageScale} 
          fill={'#6CC4C4'}
        />
      )
    }

    return (
      <Page inner>
        <div className={styles.content}>{content}</div>
        <Row>
          <Col lg={6} md={12}>
            <Button onClick={this.handleFaceCountStart}>点名</Button>
          </Col>
          <Col lg={6} md={12}>
            <Switch defaultChecked  checkedChildren="开" unCheckedChildren="关" onChange={this.handleFaceCountBackground} />
          </Col>
          <Col lg={24} md={24} style={{ marginTop: 24 }}>
            <Spin spinning={showSpinning}>
              {showFaceCount && (
                <div>
                  <Progress percent={Math.ceil(currentPerson * 100 / faceCountData.personCount)} />
                  <Stage width={window.innerWidth} height={imageHeight + 20}>
                    <Layer>
                      { showFaceCountBackground && <KImage src={faceImage} width={imageWidth} height={imageHeight}/> }
                      { rects }
                      { texts }
                    </Layer>
                  </Stage>
                </div>
              )}
            </Spin>
          </Col>
        </Row>
      </Page>
    )
  }
}

ClassDetail.protoTypes = {
  studentDetail: PropTypes.object,
}

export default ClassDetail