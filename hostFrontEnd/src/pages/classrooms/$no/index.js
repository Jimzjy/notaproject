import React, { PureComponent } from "react";
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Row, Col, Collapse, Card } from 'antd'
import Konva from 'konva';
import { Stage, Layer, Rect } from 'react-konva';
import { Page } from 'components'
import styles from './index.less'
import { timetrans } from 'utils/request'

@connect(({ classroomDetail }) => ({ classroomDetail }))
class ClassroomDetail extends PureComponent {

  personLayoutWidth = 0

  getPersonLayoutWidth = (e) => {
    if (e) {
      this.personLayoutWidth = e.getBoundingClientRect().width
    }
  }

  render() {
    const { classroomDetail } = this.props
    const { data, stats } = classroomDetail
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

    if (stats == null || stats.classroom_stats == null) {
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

    const imageScale = this.personLayoutWidth * 0.7 / stats.classroom_stats.global_width
    const imageHeight = stats.classroom_stats.global_width * imageScale
    const rects = []
    const persons = stats.classroom_stats.persons
    for (let index in persons) {
      rects.push(
        <Rect
          key={index}
          x={persons[index].left * imageScale} 
          y={persons[index].top * imageScale}
          width={persons[index].width * imageScale}
          height={persons[index].height * imageScale}
          stroke={'#42a4f4'}
        />
      )
    }

    return (
      <Page inner>
        <Row>
          <Collapse defaultActiveKey={['1']} >
            <Collapse.Panel header="教室信息" key="1">
              <div className={styles.content}>{content}</div>
            </Collapse.Panel>
          </Collapse>
        </Row>
        <Row>
          <Col lg={24} md={24} style={{ marginTop: 24 }}>
            <Card>
            <Card.Meta
              title="教室状态"
              description={`更新时间: ${timetrans(stats.update_time)} 人数: ${persons.length}`}
            />
              <div ref={this.getPersonLayoutWidth}>
                <Stage width={window.innerWidth} height={imageHeight}>
                  <Layer>
                    { rects }
                  </Layer>
                </Stage>
              </div>
            </Card>
          </Col>
        </Row>
      </Page>
    )
  }
}

ClassroomDetail.protoTypes = {
  classroomDetail: PropTypes.object,
}

export default ClassroomDetail