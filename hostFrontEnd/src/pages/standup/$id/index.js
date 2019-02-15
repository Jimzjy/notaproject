import React, { PureComponent } from "react";
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Button, Progress, Switch, Spin, Row, Col, Popover, Collapse, Carousel, notification, Upload, Icon, message } from 'antd'
import { Stage, Layer, Rect, Text } from 'react-konva';
import Konva from 'konva';
import { Document, Page as PDFPage } from 'react-pdf';
import KImage from '../components/KImage'
import QRcode from 'qrcode.react'
import { Page } from 'components'
import styles from './index.less'
import { apiPrefix } from 'utils/config'

@connect(({ standupDetail, app }) => ({ standupDetail, app }))
class StandupDetail extends PureComponent {
  state = {
    standtupStatus: false,
    standupData: {
      WReadMWriteIndex: -1,
    },
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
    },
    numPages: null,
    pageNumber: 1,
    pdfFile: null,
  }

  suws = null
  // carousel = React.createRef()
  faceCountLayoutWidth = 0

  getFaceCountLayoutWidth = (e) => {
    if (e) {
      this.faceCountLayoutWidth = e.getBoundingClientRect().width
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
    if (this.suws == null) {
      notification.open({
        message: '不处于上课状态',
        description: '请点击上课进入上课状态后进行点名',
        duration: 3,
      })
      return
    }

    const { standupDetail } = this.props
    const { data } = standupDetail
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
      console.log("open ws")
    }
    ws.onclose = function(evt) {
      console.log("close ws")
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

  handleFaceCountClose = () => {
    this.setState({
      showSpinning: false
    })

    if (this.suws) {
      this.suws.send(JSON.stringify({ FaceCountClose: true }))
    }
  }

  handleSetState = (data) => {
    this.setState(data)
  }

  handleStandUpStart = () => {
    if (this.state.standtupStatus == true) {
      return
    }

    const { numPages, pdfFile } = this.state
    if (numPages == null || pdfFile == null) {
      message.error('没有演示文稿')
      return
    }

    const { standupDetail, app } = this.props
    const { data } = standupDetail
    const getStandUpStatus = this.getStandUpStatus
    const getStandUpData = this.getStandUpData
    const handleSetState = this.handleSetState
    const handlePageToChange = this.handlePageToChange
    const handlePDFChange = this.handlePDFChange
    
    this.suws = new WebSocket(`ws://localhost:8000/stand_up?class_id=${data.class_id}&teacher_no=${app.user.username}&pdf_url=${pdfFile}&pdf_num_pages=${numPages}`)
    this.suws.onopen = function(evt) {
      console.log("open suws")
    }
    this.suws.onclose = function(evt) {
      console.log("close suws")
      handleSetState({
        standtupStatus: false,
      })   
    }
    this.suws.onmessage = function(evt) {
      const message = JSON.parse(evt.data)

      if (!getStandUpStatus()) {
        notification.open({
          message: '开始上课',
          duration: 3,
        })
        const _standupData = getStandUpData()
        _standupData.WReadMWriteIndex = message.WReadMWriteIndex
        handleSetState({standtupStatus: true, standupData: _standupData})
        handlePDFChange()
      }

      handlePageToChange(message.ChangePDFPage)
    }
    this.suws.onerror = function(evt) {
      console.log("error: " + evt.data)
      notification.open({
        message: 'ERROR',
        description: evt.data,
        duration: 3,
      })
    }
  }

  handlePageToChange = (page) => {
    if (page > 0) {
      this.handlePageChange(1)
    } else if (page < 0) {
      this.handlePageChange(-1)
    }
  }

  handleStandUpStop = () => {
    if (this.suws == null) {
      return
    }

    this.suws.close()
    this.setState({
      standtupStatus: false,
    })
    notification.open({
      message: '下课',
      duration: 3,
    })
    this.suws = null
  }

  handlePageChange = (stepDir) => {
    const numPages = this.state.numPages
    if (numPages <= 0) {
      return
    }

    let currentPage = this.state.pageNumber
    if (stepDir > 0) {
      if (currentPage + 1 <= numPages) {
        currentPage += 1
      } else {
        currentPage = 1
      }
    } else {
      if (currentPage - 1 >= 1) {
        currentPage -= 1
      } else {
        currentPage = numPages
      }
    }
    this.setState({ pageNumber: currentPage })

    if (this.suws != null) {
      this.suws.send(JSON.stringify({ CurrentPDFPage: currentPage }))
    }
  }

  getStandUpStatus = () => {
    return this.state.standtupStatus
  }

  getStandUpData = () => {
    return this.state.standupData
  }

  onDocumentLoadSuccess = ({ numPages }) => {
    this.setState({ numPages })
    // this.handlePDFChange()
  }

  handlePDFChange = () => {
    if (this.suws == null) {
      return
    }

    let packet = { CurrentPDFPage: this.state.pageNumber }
    if (this.state.numPages > 0) {
      packet.NumPDFPages = this.state.numPages
      packet.PDFUrl = this.state.pdfFile
    }
    this.suws.send(JSON.stringify(packet))
  }

  render() {
    const { standupDetail } = this.props
    const { data } = standupDetail
    const { faceCountData, showFaceCount, showFaceCountBackground, currentPerson, showSpinning, standtupStatus, standupData, pageNumber, numPages, pdfFile } = this.state
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
    const imageScale = this.faceCountLayoutWidth / faceCountData.global_width
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

    const setPDFFile = (pdfFile) => {
      this.setState({
        pdfFile: pdfFile,
        pageNumber: 1
      })
    }

    const uploadProps = {
      name: 'file',
      action: `${apiPrefix}/pdf`,
      onChange({ file }) {
        if (file.status === 'done') {
          message.success(`${file.name} 上传成功`)
          setPDFFile(`${apiPrefix}/pdf/${file.response.message}`)
        } else if (file.status === 'error') {
          message.error(`${file.name} 上传失败`)
        }
      },
    };

    return (
      <Page inner>
        <Row>
          <Collapse defaultActiveKey={['1']} >
            <Collapse.Panel header="班级信息" key="1">
              <div className={styles.content} style={{ marginBottom: 16 }}>{content}</div>
              <Button type="primary" onClick={this.handleStandUpStart}>上课</Button>
              <Button type="ghost" onClick={this.handleStandUpStop} style={{ marginLeft: 16 }}>下课</Button>
              {standtupStatus && (
                <span style={{ marginLeft: 16 }}>
                  <Popover placement="right" content={<QRcode size={160} value={`/stand_up_mobile?class_id=${data.class_id}&write_channel_index=${standupData.WReadMWriteIndex}`}/>} >
                    <Button type="primary">二维码</Button>
                  </Popover>
                </span>
              )}
              <div>
                <Upload {...uploadProps}>
                  <Button style={{ marginTop: 16 }}>
                    <Icon type="upload" /> 上传演示文稿
                  </Button>
                </Upload>
              </div>
            </Collapse.Panel>
            <Collapse.Panel header="点名信息" key="2">
              <Button onClick={this.handleFaceCountStart}>点名</Button>
              <Switch defaultChecked style={{ marginLeft: 16 }} checkedChildren="背景开" unCheckedChildren="背景关" onChange={(checked) => { this.handleSetState({showFaceCountBackground: checked}) }} />
              <div ref={this.getFaceCountLayoutWidth} >
                <Col lg={24} md={24} style={{ marginTop: 24 }}>
                  <Spin spinning={showSpinning} style={{ marginBottom: 24 }}>
                    {showFaceCount && (
                      <div>
                        <Progress percent={Math.ceil(currentPerson * 100 / faceCountData.personCount)} style={{ width: imageWidth }}/>
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
              </div>
            </Collapse.Panel>
          </Collapse>
        </Row>
        {/* <Row style={{ marginTop: 24 }} >
          <Carousel className={styles.carousel} afterChange={this.handlePageChange} ref={(ref) => {this.carousel = ref}}>
            <div><h3>1</h3></div>
            <div><h3>2</h3></div>
            <div><h3>3</h3></div>
            <div><h3>4</h3></div>
          </Carousel>
        </Row> */}
        { pdfFile != null && (
          <Row style={{ marginTop: 24 }} type="flex" justify="center">
            <Col>
              <Document
                file={pdfFile}
                onLoadSuccess={this.onDocumentLoadSuccess}
              >
                <PDFPage pageNumber={pageNumber}/>
              </Document>
              <p style={{ marginTop: 16 }}>
                Page {pageNumber} of {numPages}
                <Button shape="circle" icon="left" style={{ marginLeft: 24 }} onClick={ () => { this.handlePageChange(-1) } }/>
                <Button shape="circle" icon="right" style={{ marginLeft: 16 }} onClick={ () => { this.handlePageChange(1) } }/>
              </p>
            </Col>
        </Row>)}
      </Page>
    )
  }
}

StandupDetail.protoTypes = {
  standupDetail: PropTypes.object,
}

export default StandupDetail